package handler

import (
	myredis "file_store_net_http/cache/redis"
	"file_store_net_http/db"
	"file_store_net_http/utils"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type MultipartUploadInfo struct{
	FileHash string
	FileSize int64
	UploadId string
	ChunkSize int64
	ChunkCount int64
}

//redis用来存储什么?
//	redis用来存储正在分chunk上传的文件的信息, 每一个分块上传的文件都独立在一个Hash中;
//	直到文件上传完毕, 或者用户取消上传, redis内部的记录才删除;
// url: /file/mpupload/init
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析用户请求参数
	r.ParseForm()

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil{
		_, fn, line, _ := runtime.Caller(0)
		log.Println(fn, "_", line, ", error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.NewRespMsg(-1, "Server internal error", nil).JsonBytes())
		return
	}
	//2.获得redis连接
	conn := myredis.RedisPool().Get()
	defer conn.Close()

	uploadInfo := MultipartUploadInfo{
		FileHash:  filehash,
		FileSize:   int64(filesize),
		UploadId:   username+fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int64(math.Ceil(float64(filesize / (5 * 1024 * 1024)))),
	}

	conn.Send("MULTI")
	conn.Send("HSET","MP_"+uploadInfo.UploadId, "chunkcount", uploadInfo.ChunkCount)
	conn.Send("HSET","MP_"+uploadInfo.UploadId, "filehash", uploadInfo.FileHash)
	conn.Send("HSET","MP_"+uploadInfo.UploadId, "filesize", uploadInfo.FileSize)
	_, err = redis.Values(conn.Do("EXEC"))
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println(fn, "_", line, ", error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.NewRespMsg(-1, "redis error", nil).JsonBytes())
		return
	}

	w.Write(utils.NewRespMsg(0, "ok", nil).JsonBytes())
}

// url: /file/mpupload/uppart
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析用户请求参数
	r.ParseForm()

	//username := r.Form.Get("username")
	uploadId := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	//2.获得redis连接池连接
	conn := myredis.RedisPool().Get()
	defer conn.Close()

	filePath := utils.GetFileAbPath("chunk_local_file_store", chunkIndex)
	file,err := os.Create(filePath)
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println(fn, "_", line, ", error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.NewRespMsg(-1, "server opens file error", nil).JsonBytes())
		return
	}
	defer file.Close()

	buf := make([]byte, 5 * 1024 * 1024)
	for {
		n, err := r.Body.Read(buf)
		if err != nil {
			break
		}
		file.Write(buf[:n])
	}

	_,err = redis.Values(conn.Do("HSET", "MP_"+uploadId, "chidx_"+chunkIndex, 1))
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println(fn, "_", line, ", error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.NewRespMsg(-1, "redis error", nil).JsonBytes())
		return
	}
	w.Write(utils.NewRespMsg(0, "ok", nil).JsonBytes())
}

// url: /file/mpupload/complete
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析请求参数
	r.ParseForm()

	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	//2.获得redis连接池中的一个连接
	conn := myredis.RedisPool().Get()
	defer conn.Close()

	//3.校验index的完整性: 这里只统计index个数之和是否等于原先存储于redis中的multiPart文件信息;
	//	实际上, 如果想要实现缺失index重传机制, 就要使用一个list(也可以是位图)来统计所有上传失败的文件index片段;
	data, err := redis.Values(conn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println(fn, "_", line, ", error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.NewRespMsg(-1, "redis error", nil).JsonBytes())
		return
	}

	val := data[0].([]interface{})
		totalCount := 0
		curCount := 0
	//对于"HGETALL"命令, 需要一次遍历2个元素: key, value;
	for i := 0; i < len(val); i+=2{
		k := string(val[i].([]byte))
		v := string(val[i+1].([]byte))
		if k == "chunkcount" {
			cnt64,_ := strconv.ParseInt(v,10,64)
			totalCount = int(cnt64)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			curCount++
		}
	}

	if totalCount != curCount {
		w.Write(utils.NewRespMsg(-2, "invalid request", nil).JsonBytes())
		return
	}

	//4. 合并分块
	//	根据本地文件(chunk_local_file_store)中的index顺序把文件读取合并为一个大文件;

	//5. 更新文件表, 用户文件表
	intFileSize, _ := strconv.Atoi(filesize)
	flag,err := db.OnFileUploadFinished(filehash, filename, int64(intFileSize), "")
	if !flag || err != nil{
		_, fn, line, _ := runtime.Caller(0)
		log.Println(fn, "_", line, ", error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.NewRespMsg(-1, "file upload error", nil).JsonBytes())
		return
	}

	err = db.OnUserFileUploadFinished(username, filehash, filename, int64(intFileSize))
	if err != nil{
		_, fn, line, _ := runtime.Caller(0)
		log.Println(fn, "_", line, ", error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.NewRespMsg(-1, "user_file upload error", nil).JsonBytes())
		return
	}

	//6.向客户端响应处理结果
	w.Write(utils.NewRespMsg(0, "OK", nil).JsonBytes())
}

func CancelUploadPartHandler(w http.ResponseWriter, r *http.Request)  {
	//删除已存在的分块文件

	//删除redis缓存状态

	//更新mysql文件status
}

//查看分块上传状态
func MultipartUploadStatusHandler(w http.Response, r *http.Request)  {
	//检查分块上传状态是否有效

	//获取分块初始化信息

	//获取已上传的分块信息
}

