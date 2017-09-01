package main

import (
	"fmt"
	. "github.com/nature19862001/Chat/protocol"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
	"github.com/nature19862001/base/php"
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

func (this *Client) Close() {
	fmt.Println("client:" + this.conn.ConnAddr() + " closed")
	if this.conn != nil {
		this.conn.SetMsgParser(nil)
		this.conn.SetListener(nil)

		this.conn.Close()
		this.conn = nil
	}
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
		this.Close()
	case data := <-this.recvChan:
		this.process(data)
	}
	close(this.recvChan)
}

func (this *Client) process(data []byte) {
	msgid := Uint16(data)

	switch msgid {
	case MsgId_ReqLogin:
		uid := Uint64(data[2:10])
		password := data[10:]
		ret := new(MsgRetLogin)
		ret.MsgId = MsgId_RetLogin
		code := gDataManager.checkLogin(uid, string(password))
		if code == ERR_NONE {
			ret.Result = uint16(ERR_NONE)
			//copy(ret.IP[0:], []byte("127.0.0.1"))
			//ret.Port = 9090
			ok := gDataManager.setUserOnline(uid)
			if !ok {
				ret.Result = uint16(ERR_REDIS)
				this.send(Bytes(ret))
				this.Close()
				return
			} else {
				fmt.Println("addr:" + this.conn.ConnAddr() + " logined success")
			}
		} else {
			ret.Result = uint16(code)
			this.send(Bytes(ret))
			this.Close()
		}
		this.send(Bytes(ret))
	case MsgId_ReqAppLogin:
		appname := string(data[2:34])
		password := string(data[34:])
		//check app login
		//if login ok, then wait for app server verify
		code := gDataManager.checkAppLogin(appname, password)
		ret := new(MsgRetAppLogin)
		ret.MsgId = MsgId_RetAppLogin
		ret.Result = uint16(code)
		if code == ERR_NONE {
			//copy(ret.IP[0:], []byte("127.0.0.1"))
			//ret.Port = 9090
			code := gDataManager.setAppOnline(appname)
			ret.Result = uint16(code)
			if code != ERR_NONE {
				this.send(Bytes(ret))
				this.Close()
				return
			} else {
				token := Authcode(String(time.Now().Unix())+":"+appname, "ENCODE")

				ret.Count = byte(len(token))
				ret.Token = []byte(token)

				newAppClient(appname, this.conn)
				fmt.Println("addr:" + this.conn.ConnAddr() + " app logined success")
			}
		} else {
			//else ret login failed.
			this.send(Bytes(ret))
			this.Close()
		}
		this.send(Bytes(ret))
	}
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

func (this *Client) send(buff []byte) {
	this.conn.Send(append(Bytes(int16(len(buff))), buff...))
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
	this.Close()
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
