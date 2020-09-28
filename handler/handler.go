package handler

import (
	"03_gin_filestore_own/file_store_net_http/meta"
	"03_gin_filestore_own/file_store_net_http/utils"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	filepath "path/filepath"
	"runtime"
	"time"
)

var(
	dirPath = filepath.Dir()
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//返回界面
		f := utils.GetFileAbPath("static", "view", "home.html")
		data, err := ioutil.ReadFile(f)
		if err != nil{
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
			fmt.Println(err.Error())
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
		//fmeta.FileSha1 = utils.FileSsha1(newFile)
	}
}
