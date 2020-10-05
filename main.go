package main

import (
	"file_store_net_http/handler"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

//服务器实现方式1: 1.注册Mux多路选择器; 2.给mux注册handler, handleFunc; 3.创建Server对象; 4.ServeAndListen;
//服务器实现方式2: 1.http.handle(), http.handleFunc()注册处理func; 2.http.ListenAndServe(port, mux);
//	这里使用的都是http中的DefaultMux, DefaultServer;
//	注册的handler与handleFunc怎么被调用? 本质是还是调用其serveHTTP(w, r)方法;

//中间件的实现:
//	调用时:
//		1. 对于handleFunc(), 先调用http.HandlerFunc()转换为Handler类型; 对于http.Handler对象, 则无需转换;
//		2. 将上述返回值, 作为中间件函数参数调用;
//		3. 将上述返回值, 再做为另一个中间件函数调用;
//		4. 中间件函数在http.Handle()中作为参数被调用;
//	写中间件函数
//		1. 参数, 返回值为http.Handler类型;
//		2. 在中间件函数中:
//			<1> 返回值为http.Handler类型;
//			<2> 在返回值函数中, 逻辑为: pre处理 + 函数参数调用 + post处理;
func main() {

	//这些都是要加interceptor, 先验证用户的登录状态, 才提供文件上传/下载服务;
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/file/fastupload", handler.TryFastUploadHandler)

	//这个handler对应的处理请求, 来自于client端;
	http.HandleFunc("/file/mpupload/init", handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart", handler.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", handler.HTTPInterceptor(handler.CompleteUploadHandler))

	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SigninHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserinfoHandler))

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		_, filePath, line, _ := runtime.Caller(0)
		fmt.Printf("file: %v; line: %v; error: %v\n",filePath, line,err)
		time.Sleep(time.Second * 3)
	}
}
