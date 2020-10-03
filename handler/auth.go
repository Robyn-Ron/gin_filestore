package handler

import "net/http"

//拦截器中间件
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		username := r.Form.Get("username")
		token := r.Form.Get("token")

		if len(username) < 3 || !IsTokenValid(token){
			w.WriteHeader(http.StatusForbidden)
			return
		}
		//最后调用功能函数handlerFunc
		h(w,r)
	}
}
