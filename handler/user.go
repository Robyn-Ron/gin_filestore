package handler

import (
	"CloudWebOfGin/db"
	"CloudWebOfGin/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//加密盐值
const (
	pwd_salt = "#*888!@"
)

//用户注册
func SignupHandler(w http.ResponseWriter, r *http.Request)  {
	//判断客户端请求方式
	if r.Method == http.MethodGet{
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(data)
	}else if r.Method == http.MethodPost{
		r.ParseForm()

		username := r.Form.Get("username")
		passwd := r.Form.Get("password")
		phone := r.Form.Get("phone")

		//参数校验
		if len(username) < 3 || len(passwd) < 5{
			w.Write([]byte("Invalid parameter"))
			return
		}

		if len(phone) != 11{
			w.Write([]byte("Invalid phone format"))
			return
		}

		//密码进行加密处理
		enc_passwd := util.Sha1([]byte(passwd+pwd_salt))

		flag := db.UserSignup(username, enc_passwd, phone)

		if flag{
			w.Write([]byte("注册成功"))
		}else {
			w.Write([]byte("注册失败"))
		}
	}
}

//用户登录
func SigninHandler(w http.ResponseWriter, r *http.Request)  {
	//校验用户名和密码
	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encry_password := util.Sha1([]byte(password + pwd_salt))

	flag := db.UserSignin(username, encry_password)
	if !flag{
		//登录失败
		w.Write([]byte("登录失败"))
		return
	}

	//生成登录凭证：1.token 或者 2.基于session和cookie
	token := GenToken(username)
	flag = db.UpdateToken(username, token)

	if !flag{
		w.Write([]byte("生成token信息失败"))
		return
	}

	//登录成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))

	resp := util.RespMsg{
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

	w.Write(resp.JSONBytes())
}

//生成token凭证信息
func GenToken(username string) string {
	//md5(username + timestamp + token_salt) + timestamp[:8] 组成40位的token
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

//查询用户信息
func UserinfoHandler(w http.ResponseWriter, r *http.Request)  {
	//1.解析请求参数
	r.ParseForm()

	username := r.Form["username"][0]

	//2.查询用户信息
	user, err := db.GetUserInfo(username)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//3.组装并响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg: "success",
		Data: user,
	}

	w.Write(resp.JSONBytes())
}

func IsTokenValid(token string) bool {
	//TODO: 判断token时效性，是否过期(从token的后八位的时间戳来校验时间是否过期)
	//TODO: 从数据库表tbl_user_token查询username对应的token信息
	//TODO: 对比两个token是否一致
	if len(token) < 5{
		fmt.Println("请求被拦截器拦截住了")
		return false
	}

	return true
}