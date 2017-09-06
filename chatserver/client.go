package main

import (
	"fmt"
	. "github.com/nature19862001/Chat/protocol"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
	"github.com/nature19862001/base/php"
	"strings"
	"time"
)

type Client struct {
	conn     gtnet.IConn
	recvChan chan []byte
}

func newClient(conn gtnet.IConn) *Client {
	c := &Client{conn: conn}
	c.serve()
	return c
}

func (this *Client) serve() {
	this.recvChan = make(chan []byte)
	this.conn.SetMsgParser(this)
	this.conn.SetListener(this)

	go this.startProcess()
}

func (this *Client) startProcess() {
	timer := time.NewTimer(time.Second * 30)

	select {
	case <-timer.C:
		this.conn.Close()
	case data := <-this.recvChan:
		if this.process(data) {
			this.conn.Close()
		}
	}

	this.conn = nil
}

func (this *Client) process(data []byte) bool {
	msgid := Uint16(data)
	result := false

	switch msgid {
	case MsgId_ReqLogin:
		uid := Uint64(data[2:10])
		password := string(data[10:])

		code := gDataManager.checkLogin(uid, password)

		ret := new(MsgRetLogin)
		ret.MsgId = MsgId_RetLogin

		if code == ERR_NONE {
			code = gDataManager.setUserOnline(uid)
			if code == ERR_NONE {
				newChatClient(uid, this.conn)
				fmt.Println("addr:" + this.conn.ConnAddr() + " logined success")
			}
		}
		ret.Result = uint16(code)
		result = code != ERR_NONE
		this.send(Bytes(ret))
	case MsgId_ReqAppLogin:
		appname := string(data[2:34])
		password := string(data[34:])
		//check app login
		//if login ok, then wait for app server verify
		code := gDataManager.checkAppLogin(appname, password)

		ret := new(MsgRetAppLogin)
		ret.MsgId = MsgId_RetAppLogin

		if code == ERR_NONE {
			code = gDataManager.setAppOnline(appname)
			if code == ERR_NONE {
				newAppClient(appname, this.conn)
				fmt.Println("addr:" + this.conn.ConnAddr() + " app logined success")
			}
		}
		ret.Result = uint16(code)
		result = code != ERR_NONE
		this.send(Bytes(ret))
	case MsgId_ReqTokenLogin:
		token := data[2:]
		str := Authcode(string(token))
		pos := strings.Index(str, ":")

		ret := new(MsgRetTokenLogin)
		ret.MsgId = MsgId_RetTokenLogin

		code := ERR_NONE
		timestamp := Int64(str[:pos])

		if time.Now().Unix()-timestamp > 3600 {
			code = ERR_TIME_OUT
			result = true
		} else {
			uid := Uint64(str[pos:])
			newChatClient(uid, this.conn)
			fmt.Println("addr:" + this.conn.ConnAddr() + " logined with token success")
		}
		ret.Result = uint16(code)
		this.send(Bytes(ret))
	}

	return result
}

func (this *Client) send(buff []byte) {
	this.conn.Send(append(Bytes(int16(len(buff))), buff...))
}

func (this *Client) ParseHeader(data []byte) int {
	size := Int(data)
	//fmt.Println("header size :", size)
	//p.conn.Send(data)
	return size
}

func (this *Client) ParseMsg(data []byte) {
	newdata := make([]byte, len(data))
	copy(newdata, data)
	this.recvChan <- newdata
}

func (this *Client) OnError(errorcode int, msg string) {
	//fmt.Println("tcpserver error, errorcode:", errorcode, "msg:", msg)
}

func (this *Client) OnPreSend([]byte) {

}

func (this *Client) OnPostSend([]byte, int) {
	// if this.state == state_logouted {
	// 	this.Close()
	// }
}

func (this *Client) OnClose() {
	//fmt.Println("tcpserver closed:", this.clientAddr)
	//this.Close()
}

func (this *Client) OnRecvBusy([]byte) {
	//str := "server is busy"
	//p.conn.Send(Bytes(int16(len(str))))
	//this.conn.Send(append(Bytes(int16(len(str))), []byte(str)...))
}

func (this *Client) OnSendBusy([]byte) {
	// str := "server is busy"
	// p.conn.Send(Bytes(int16(len(str))))
	// p.conn.Send([]byte(str))
}

var UC_KEY string = "1111aaaa"

func Authcode(str string, args ...interface{}) string {
	operation := "DECODE"
	key := ""
	var expiry int64 = 0

	argc := len(args)
	if argc >= 3 {
		texpiry, ok := args[2].(int64)

		if ok {
			expiry = texpiry
		}

		ttexpiry, ok := args[2].(int)

		if ok {
			expiry = int64(ttexpiry)
		}
	}

	if argc >= 2 {
		tkey, ok := args[1].(string)

		if ok {
			key = tkey
		}
	}

	if argc >= 1 {
		toperation, ok := args[0].(string)

		if ok {
			operation = toperation
		}
	}

	ckey_length := 4

	if key == "" {
		key = php.Md5(UC_KEY)
	} else {
		key = php.Md5(key)
	}
	//key = php.Md5(key ? key : UC_KEY)
	keya := php.Md5(php.Substr(key, 0, 16))
	keyb := php.Md5(php.Substr(key, 16, 16))

	keyc := ""
	if ckey_length != 0 {
		if operation == "DECODE" {
			keyc = php.Substr(str, 0, ckey_length)
		} else {
			keyc = php.Substr(php.Md5(String(php.Microtime())), -ckey_length)
		}
	}

	cryptkey := keya + php.Md5(keya+keyc)
	key_length := len(cryptkey)

	if operation == "DECODE" {
		str = php.Base64_decode(php.Substr(str, ckey_length))
	} else {
		var rexpiry int64
		if expiry == 0 {
			rexpiry = 0
		} else {
			rexpiry = expiry + php.Time()
		}
		str1 := php.Sprintf("%010d", rexpiry)
		str = str1 + php.Substr(php.Md5(str+keyb), 0, 16) + str
	}
	string_length := len(str)

	result := ""
	box := php.Range(0, 255, 1)

	j := 0
	i := 0
	a := 0
	rndkey := make([]int, 256)
	for i = 0; i <= 255; i++ {
		rndkey[i] = php.Ord(string(cryptkey[i%key_length]))
	}

	j = 0
	i = 0
	for ; i < 256; i++ {
		j = (j + box[i] + rndkey[i]) % 256
		tmp := box[i]
		box[i] = box[j]
		box[j] = tmp
	}

	j = 0
	i = 0
	for ; i < string_length; i++ {
		a = (a + 1) % 256
		j = (j + box[a]) % 256
		tmp := box[a]
		box[a] = box[j]
		box[j] = tmp
		ntmp := box[(box[a]+box[j])%256]
		//nstr := string(rune(str[i])))
		ndata := int(str[i])
		nres := ndata ^ ntmp
		nstr := php.Chr(nres)
		result = result + nstr //php.Chr(nres)
	}

	if operation == "DECODE" {
		num := Int64(php.Substr(result, 0, 10)) //utils.StrToInt64(string(byteresult[:10])) //php.Substr(result, 0, 10))
		if (num == 0 || num-php.Time() > 0) && php.Substr(result, 10, 16) == php.Substr(php.Md5(php.Substr(result, 26)+keyb), 0, 16) {
			return php.Substr(result, 26) //string(byteresult[26:]) //php.Substr(result, 26)
		} else {
			return ""
		}
	} else {
		return keyc + php.Base64_encode(result) //php.Str_replace("=", "", php.Base64_encode(result)) //base64.StdEncoding.EncodeToString(byteresult)) //php.Base64_encode(result))
	}
}
