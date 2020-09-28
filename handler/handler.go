package handler

import (
	"CloudWebOfGin/db"
	"CloudWebOfGin/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"CloudWebOfGin/meta"
)

func UploadHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET"{
		//返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil{
			io.WriteString(w, "internel server error")
			return
		}

		io.WriteString(w, string(data))
	}else if r.Method == "POST"{
		//解析请求参数
		r.ParseForm()

		//接收文件流以存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil{
			fmt.Println(err.Error())
			return
		}

		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil{
			fmt.Println(err.Error())
			return
		}

		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil{
			fmt.Printf("Failed to save data into file, err: %s\n",err.Error())
			return
		}

		//TODO:计算文件签名算法比较耗时，之后可以抽取出去做成微服务
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)


		_ = meta.UpdateFileMetaDB(fileMeta) //更新文件表

		//TODO：更新用户文件表信息
		username := r.Form.Get("username")
		flag := db.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)

		if !flag{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}else {
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		}
	}
}

//上传已完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request)  {
	io.WriteString(w, "Upload finished!")
}

//批量获取文件元信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	limit, err := strconv.Atoi(r.Form["limit"][0])
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username := r.Form["limit"][0]

	userFileArray, err := db.QueryUserFileMetas(username, limit)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := util.RespMsg{
		Code: 0,
		Msg: "success",
		Data: userFileArray,
	}

	w.Write(resp.JSONBytes())
}

//获取单个文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request)  {
	//解析客户端的参数信息
	r.ParseForm()

	//获取参数
	fileHash := r.Form["filehash"][0]

	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fileMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	fsha1 := r.Form.Get("filehash")
	fmeta := meta.GetFileMeta(fsha1)

	f, err := os.Open(fmeta.Location)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+ fmeta.FileName +"\"")
	w.Write(data)
}

//更新文件元信息
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0"{
		//403 错误（身份验证）
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST"{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName

	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//删除文件
func FileDeleteHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	fsha1 := r.Form.Get("filehash")

	//删除具体文件
	fMeta := meta.GetFileMeta(fsha1)
	os.Remove(fMeta.Location)

	//删除文件索引
	meta.RemoveFileMeta(fsha1)

	w.WriteHeader(http.StatusOK)
}

//尝试秒传文件接口
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request)  {
	//TODO:1.解析请求参数
	r.ParseForm()

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO:2.从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO:3.查不到记录则返回秒传失败
	if fileMeta == nil{
		resp := util.RespMsg{
			Code: -1,
			Msg: "秒传失败，请访问普通上传文件接口",
		}

		w.Write(resp.JSONBytes())
	}

	//TODO:4.如果上传过则将文件信息写入用户文件表，返回成功
	flag := db.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))

	if !flag{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := util.RespMsg{
		Code: 0,
		Msg: "秒传成功",
	}

	w.Write(resp.JSONBytes())
}