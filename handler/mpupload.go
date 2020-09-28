package handler

import (
	"CloudWebOfGin/cache/redis"
	"CloudWebOfGin/db"
	"CloudWebOfGin/util"
	"fmt"
	o_redis "github.com/garyburd/redigo/redis"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

//初始化信息
type MultipartUploadInfo struct {
	FileHash string
	FileSize int
	UploadId string //针对于每一次上传动作的id
	ChunkSize int //分块大小
	ChunkCount int //分块数量
}

//初始化分块上传
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request)  {
	//1.解析用户请求参数
	r.ParseForm()

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil{
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
	}
	//2.获得redis连接
	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	//3.生成分块上传的初始化信息
	uploadInfo := MultipartUploadInfo{
		FileHash: filehash,
		FileSize: filesize,
		UploadId: username + fmt.Sprintf("%x", time.Now().UnixNano()), //定义规则：用户名+时间戳
		ChunkSize: 5 * 1024 * 1024, //5M
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	//4.将初始化信息写入到redis缓存
	rConn.Do("HSET", "MP_"+uploadInfo.UploadId, "chuncount", uploadInfo.ChunkCount) //命令  key  value
	rConn.Do("HSET", "MP_"+uploadInfo.UploadId, "filehash", uploadInfo.FileHash)
	rConn.Do("HSET", "MP_"+uploadInfo.UploadId, "filesize", uploadInfo.FileSize)

	//5.将响应初始化数据返回客户端
	w.Write(util.NewRespMsg(0, "ok", nil).JSONBytes())
}

//上传文件分块
func UploadPartHandler(w http.ResponseWriter, r *http.Request)  {
	//1.解析用户请求参数
	r.ParseForm()

	username := r.Form.Get("username")
	uploadId := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	//2.获得redis连接池连接
	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	//3.获得文件句柄，用于存储分块内容
	//3.1 创建目录
	fpath := "/data/" + uploadId + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil{
		w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()

	var buf = make([]byte, 1024 * 1024) //1M大小
	for{
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil{
			break
		}
	}

	//4.更新redis存储状态
	rConn.Do("HSET", "MP_" + uploadId, "chkidx_" + chunkIndex, 1)

	//5.返回处理结果到客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

//通知上传合并接口
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request)  {
	//1.解析请求参数
	r.ParseForm()

	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	//2.获得redis连接池中的一个连接
	rConn := redis.RedisPool().Get()
	defer rConn.Close()

	//3.通过uploadId查询redis并判断是否所有分块上传完成
	data, err := o_redis.Values(rConn.Do("HGETALL", "MP_"+upid))
	if err != nil{
		w.Write(util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
		return
	}

	totalCount := 0
	chunkCount := 0

	for i := 0; i < len(data); i+=2{
		k := string(data[i].([]byte))
		v := string(data[i + 1].([]byte))

		if k == "chunkcount"{
			totalCount, _ = strconv.Atoi(v)
		}else if strings.HasPrefix(k, "chkidx_") && v == "1"{
			chunkCount++
		}
	}

	if totalCount != chunkCount{
		w.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		return
	}

	//4.合并分块


	//5.更新唯一文件表和用户文件表
	intFileSize, _ := strconv.Atoi(filesize)
	flag := db.OnFileUploadFinished(filehash, filename, int64(intFileSize), "")
	if !flag{
		w.Write(util.NewRespMsg(-1, "", nil).JSONBytes())
		return
	}

	flag = db.OnUserFileUploadFinished(username, filehash, filename, int64(intFileSize))
	if !flag{
		w.Write(util.NewRespMsg(-1, "", nil).JSONBytes())
		return
	}

	//6.向客户端响应处理结果
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())

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