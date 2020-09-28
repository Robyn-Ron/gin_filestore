package main

import (
	"fmt"
	"net/http"
)
import "CloudWebOfGin/handler"

func main(){
	//文件相关
	http.HandleFunc("/file/upload", handler.UploadHandler) //上传
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler) //暂时没用
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler) //获取单个文件元信息
	http.HandleFunc("/file/multimeta",handler.FileQueryHandler) //获取批量文件元信息
	http.HandleFunc("/file/download", handler.DownloadHandler) //下载
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))
	http.HandleFunc("/file/mpupload/init", handler.HTTPInterceptor(handler.InitialMultipartUploadHandler)) //分块上传接口
	http.HandleFunc("/file/mpupload/uppart", handler.HTTPInterceptor(handler.UploadPartHandler))//分块上传接口
	http.HandleFunc("/file/mpupload/complete", handler.HTTPInterceptor(handler.CompleteUploadHandler))//通知分块上传完成

	//用户相关
	http.HandleFunc("/user/signup", handler.SignupHandler) //注册
	http.HandleFunc("/user/signin", handler.SigninHandler) //登录
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserinfoHandler)) //加拦截器的路由
	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println(err.Error())
	}
}
