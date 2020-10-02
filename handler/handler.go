package handler

import (
	"encoding/json"
	"file_store_net_http/meta"
	"file_store_net_http/utils"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)


// url: /file/upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回界面
		f := utils.GetFileAbPath("static", "view", "home.html")
		data, err := ioutil.ReadFile(f)
		if err != nil{
			_, fn, line, _ := runtime.Caller(0)
			fmt.Println(fn, "_", line, ", error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))
	} else {//POST
		//上传文件
		//解析表单
		//任务:1.保存metafile; 2.保存file为本地文件;
		file, head, err := r.FormFile("file")
		if err != nil {
			_, fn, line, _ := runtime.Caller(0)
			fmt.Println(fn, "_", line, ", error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "get FormFile error!")
			return
		}
		defer file.Close()
		fmeta := meta.FileMeta{
			FileName: head.Filename,
			Location: utils.GetFileAbPath("local_file_store", head.Filename),
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		meta.UpdateFileMeta(fmeta)
		//存储文件到本地
		newFile, err := os.Create(fmeta.Location)
		if err != nil {
			_, fn, line, _:=runtime.Caller(0)
			fmt.Println(fn, "_", line, ", error:", err)
		}
		defer newFile.Close()

		//ioutil.WriteFile() 这里使用create()更加安全
		fmeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			_, fn, line, _:=runtime.Caller(0)
			fmt.Println(fn, "_", line, ", error:", err)
		}

		//TODO: 计算文件签名比较耗时, 之后可抽取做成微服务;
		newFile.Seek(0,0)
		//这里FileSha1的参数, 是os.File类型, 而不是mime.multipart.File类型(这个interface只能io.reader)
		fmeta.FileSha1 = utils.FileSha1(newFile)

		//保存文件在内存中的map[file_hash]FileMeta
		//meta.UpdateFileMeta(fmeta)
		//保存文件在DB中
		meta.UpdateFileMetaDB(fmeta)

		//重定向到上传成功接口
		http.Redirect(w,r, "/file/upload/suc", http.StatusFound)
	}

}
//url: /file/upload/suc
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished")
}

// url: /file/meta
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// r.Form属性, 相比于r.FormValue优点在于, 表单中可能用多个同名的数据, 使用前者可以接受多个, 使用后者只能接受最前一个;
	fileHash := r.Form["filehash"][0]
	//fileMeta := meta.GetFileMeta(fileHash)
	fileMeta, err := meta.GetFileMetaDB(fileHash)
	data, err := json.Marshal(fileMeta)
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		fmt.Println(fn, "_", line, ", error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//默认有: w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// url: /file/multimeta
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	//获取username
	//获取limit
	//从db中查询fileMeta多个对象信息
	//构建RespMsg, 返回前端
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	//获取MetaFile的hash信息
	r.ParseForm()
	fileHash := r.Form["filehash"][0]
	//用hash查询metaFile信息
	metaFile := meta.GetFileMeta(fileHash)
	//构建返回w信息
	respFileContent,err := ioutil.ReadFile(metaFile.Location)
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		fmt.Println(fn,"_",line,", error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Header可以指定全新的map[string]string, 也可以基于原有的Header()来set, delete(), add(), get()
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+ filepath.Base(metaFile.Location)+"\"")
	w.Write(respFileContent)
}

// url: /file/update
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

// url:/file/delete
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileSha1 := r.Form["filehash"][0]

	metaFile := meta.GetFileMeta(fileSha1)
	os.Remove(metaFile.Location)

	meta.RemoveFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
}
