package main

import (
	"fmt"
	//. "github.com/nature19862001/Chat/common"
	//"github.com/nature19862001/base/gtnet"
	. "github.com/nature19862001/base/common"
	"io"
	"net/http"
)

func loadBanlanceInit() {
	go startHTTPServer()
	go starUserRegister()
}

func startHTTPServer() {
	http.HandleFunc("/serverlist", getServerList)
	http.ListenAndServe(":9001", nil)
}

func getServerList(rw http.ResponseWriter, req *http.Request) {
	serverlist := gDataManager.getServerList()

	ret := "{\r\n\tserverlist:\r\n\t[\r\n"
	for i := 0; i < len(serverlist); i++ {
		ret += "\t\t{ addr:\"" + serverlist[i] + "\" },\r\n"
	}
	ret += "\t]\r\n"
	ret += "}\r\n"

	io.WriteString(rw, ret)
}

func starUserRegister() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/create", create)
	http.ListenAndServe(":8080", nil)
}

func register(rw http.ResponseWriter, req *http.Request) {
	//req.ParseForm()
	fmt.Println(req.Method)
	fmt.Println(req.RemoteAddr)
	// fmt.Println(req.PostForm)
	// fmt.Println(req.Form["username"])
	// fmt.Println(req.PostForm["username"])
	// fmt.Println(req.PostFormValue("username"))
	// nickname := req.PostFormValue("nickname")
	// password := req.PostFormValue("password")
	// regip := req.RemoteAddr
	// method := req.Method

	ret := "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">"
	ret += "<html xmlns=\"http://www.w3.org/1999/xhtml\">"
	ret += "<head>"
	ret += "<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" />"
	ret += "<title>register</title>"
	ret += "<meta content=\"GTech Inc.\" name=\"Copyright\" />"
	ret += "<script src=\"http://cdn.bootcss.com/blueimp-md5/1.1.0/js/md5.min.js\"></script>"
	ret += "</head>"
	ret += "<body>"

	ret += "<form method=\"post\" action=\"/create\" onsubmit=\"return true;\">"
	ret += "昵称：<input type=\"text\" name=\"nickname\" />"
	ret += "<br/>"
	ret += "密码：<input type=\"password\" name=\"password1\" oninput=\"document.getElementById('password').value = md5(this.value);\" onpropertychange=\"document.getElementById('password').value = md5(this.value);\" />"
	ret += "<input type=\"hidden\" name=\"password\" id=\"password\" />"
	ret += "<br/>"
	ret += "<input type=\"submit\" name=\"login_button\" value=\"提交\">"
	ret += "</form>"

	ret += "</body>"
	ret += "</html>"
	io.WriteString(rw, ret)
}

func create(rw http.ResponseWriter, req *http.Request) {
	var ok bool
	var uid uint64
	nickname := req.PostFormValue("nickname")
	password := req.PostFormValue("password")
	regip := req.RemoteAddr
	method := req.Method

	ret := "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">"
	ret += "<html xmlns=\"http://www.w3.org/1999/xhtml\">"
	ret += "<head>"
	ret += "<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" />"
	ret += "<title>register</title>"
	ret += "<meta content=\"GTech Inc.\" name=\"Copyright\" />"
	ret += "</head>"
	ret += "<body>"

	if method != "POST" {
		// ret := "{\r\n\terrorcode:1,\r\n"
		// ret = "\r\n\terror:\"need post\",\r\n"
		// ret += "}"
		ret += "<span>请使用post方法!</span><br/>"
		ret += "<form method=\"post\" action=\"/create\">"
		ret += "昵称：<input type=\"text\" name=\"nickname\" />"
		ret += "<br/>"
		ret += "密码：<input type=\"password\" name=\"password\" />"
		ret += "<br/>"
		ret += "<input type=\"submit\" name=\"login_button\" value=\"提交\">"
		ret += "</form>"
		//io.WriteString(rw, ret)
		goto end
	}

	if nickname == "" {
		// ret := "{\r\n\terrorcode:1,\r\n"
		// ret = "\r\n\terror:\"need nickname\",\r\n"
		// ret += "}"
		ret += "<span>请输入昵称!</span><br/>"
		ret += "<form method=\"post\" action=\"/create\">"
		ret += "昵称：<input type=\"text\" name=\"nickname\" />"
		ret += "<br/>"
		ret += "密码：<input type=\"password\" name=\"password\" />"
		ret += "<br/>"
		ret += "<input type=\"submit\" name=\"login_button\" value=\"提交\">"
		ret += "</form>"
		//io.WriteString(rw, ret)
		goto end
	}

	if password == "" {
		// ret := "{\r\n\terrorcode:2,\r\n"
		// ret = "\r\n\terror:\"need password\",\r\n"
		// ret += "}"
		ret += "<span>请输入密码!</span><br/>"
		ret += "<form method=\"post\" action=\"/create\">"
		ret += "昵称：<input type=\"text\" name=\"nickname\" />"
		ret += "<br/>"
		ret += "密码：<input type=\"password\" name=\"password\" />"
		ret += "<br/>"
		ret += "<input type=\"submit\" name=\"login_button\" value=\"提交\">"
		ret += "</form>"
		//io.WriteString(rw, ret)
		goto end
	}

	ok, uid = gDataManager.createUser(nickname, password, regip)

	if !ok {
		// ret := "{\r\n\terrorcode:3,\r\n"
		// ret = "\r\n\terror:\"server error\",\r\n"
		// ret += "}"
		ret += "<span>注册失败，服务器内部错误!</span><br/>"
		ret += "<form method=\"post\" action=\"/create\">"
		ret += "昵称：<input type=\"text\" name=\"nickname\" />"
		ret += "<br/>"
		ret += "密码：<input type=\"password\" name=\"password\" />"
		ret += "<br/>"
		ret += "<input type=\"submit\" name=\"login_button\" value=\"提交\">"
		ret += "</form>"
		//io.WriteString(rw, ret)
		goto end
	}

	// ret := "{\r\n\terrorcode:0,\r\n"
	// ret = "\r\n\terror:\"\",\r\n"
	// ret = "\r\n\tuid:" + String(uid) + ",\r\n"
	// ret += "}"
	ret += "<span>注册成功，登录账号：" + String(uid) + "</span><br/>"
end:
	ret += "</body>"
	ret += "</html>"
	io.WriteString(rw, ret)
}
