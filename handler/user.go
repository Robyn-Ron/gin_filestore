package handler

import (
	"errors"
	"file_store_net_http/db"
	"file_store_net_http/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
)

const (
	pwd_salt ="#hao666dePerson!!@somebody"
)

// url: /user/signup
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile(utils.GetFileAbPath("static","view", "signup.html"))
		if err != nil {
			_, fn, line, _ := runtime.Caller(0)
			fmt.Println(fn,"_",line,", error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	} else if r.Method == http.MethodPost{
		r.ParseForm()

		//post, get方式获取参数, 最好还是使用推荐的方法, 不要使用Form这种形式;
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		phone := r.PostFormValue("phone")

		//密码加salt处理, 存到db中;
		encrypt_password := utils.Sha1([]byte(password + pwd_salt))
		err := db.UserSignup(username, encrypt_password, phone)
		if err != nil {
			_, fn, line, _ := runtime.Caller(0)
			fmt.Println(fn,"_",line,", error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("注册失败"))
			return
		}
		w.Write([]byte("注册成功"))
	}

}

// url: /user/signin
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username :=	r.Form.Get("username")
	password := r.Form.Get("password")

	encrypt_password := utils.Sha1([]byte(password+pwd_salt))
	err := db.UserSignin(username, encrypt_password)
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		fmt.Println(fn,"_",line,", error:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("登陆失败"))
		return
	}

	//生成登陆凭证
	token := GenToken(username)
	//每次登陆都是不同的token
	err = db.UpdateToken(username, token)
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		fmt.Println(fn,"_",line,", error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("服务器端生成登陆token失败"))
		return
	}
	resp := utils.RespMes{
		Code: 0,
		Msg: "OK",
		Data: struct {
			Location string
			Username string
			Token string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token: token,
		},
	}

	w.Write(resp.JsonBytes())

}

// url: /user/info
func UserinfoHandler(w http.ResponseWriter, r *http.Request) {
	//因为在获取userInfo信息之前, 先使用了interceptor中间件, 所以此步骤就是查db, 获取user信息;
	r.ParseForm()

	username := r.Form.Get("username")
	user, err := db.GetUserInfo(username)
	if err != nil{
		_, fn, line, _ := runtime.Caller(0)
		fmt.Println(fn,"_",line,", error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("服务器端查询用户信息失败"))
		return
	}
	resp := utils.RespMes{
		Code: 0,
		Msg:  "success",
		Data: user,
	}
	w.Write(resp.JsonBytes())
}

func GenToken(username string) string {
	ts := fmt.Sprintf("%x",time.Now().Unix())
	tokenPrefix := utils.MD5([]byte(username+ts+"_tokensalt"))
	return tokenPrefix +ts[:8]
}

//bool类型的返回值, 基本上不需要返回error
func IsTokenValid(token string) bool {
	//TODO: 判断token是否过期, 过期则需要重新登录
	//TODO: 判断token是否一致, 未被修改过; 如果是jwt, 则应该使用private_key来验证, 这里直接用查db中的token信息来比较一致性;
	if len(token) < 5 {
		_, fn, line, _ := runtime.Caller(0)
		err := errors.New("token验证失败")
		log.Println(fn,"_",line,", error:", err)
		return false
	}
	return true
}
