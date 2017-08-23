package main

import (
	"fmt"
	. "github.com/nature19862001/Chat/protocol"
	. "github.com/nature19862001/base/common"
	"github.com/nature19862001/base/gtnet"
	"github.com/nature19862001/base/php"
	"sync"
	"time"
)

const (
	state_none          = 0
	state_connected int = 1
	state_logined   int = 2
	state_logouted  int = 3
	state_del       int = 4
	state_verify    int = 5
)

var msgProcesserMap map[uint16]func(*Client, []byte)

func init() {
	msgProcesserMap = make(map[uint16]func(*Client, []byte))
	msgProcesserMap[MsgId_ReqFriendList] = OnReqFriendList
	msgProcesserMap[MsgId_ReqFriendAdd] = OnReqFriendAdd
	msgProcesserMap[MsgId_ReqFriendDel] = OnReqFriendDel
	msgProcesserMap[MsgId_ReqUserToBlack] = OnReqUserToBlack
	msgProcesserMap[MsgId_ReqRemoveUserInBlack] = OnReqRemoveUserInBlack
	msgProcesserMap[MsgId_ReqMoveFriendToGroup] = OnReqMoveFriendToGroup
	msgProcesserMap[MsgId_ReqSetFriendVerifyType] = OnReqSetFriendVerifyType
	msgProcesserMap[MsgId_Message] = OnMessage
}

type Client struct {
	conn         gtnet.IConn
	lock         *sync.Mutex
	isVerifyed   bool
	timer        *time.Timer
	verfiycount  int
	countTimeOut int
	tickChan     chan int
	uid          uint64
	password     string
	state        int
	clientAddr   string

	isAppLogin    bool
	isAppVerifyed bool
	appName       string
}

func (this *Client) Close() {
	fmt.Println("client:" + this.clientAddr + "closed")
	if this.conn != nil {
		this.conn.SetMsgParser(nil)
		this.conn.SetListener(nil)
		if this.isVerifyed {
			gDataManager.setUserOffline(this.uid)
			// n, err := RedisConn.Do("SREM", "user:online", String(this.uid))
			// if err != nil {
			// 	fmt.Println(err.Error())
			// }
			// if n == nil {
			// 	fmt.Println("sadd cmd failed!")
			// }
			// n, err = RedisConn.Do("SREM", "user:online:password", this.password)
			// if err != nil {
			// 	fmt.Println(err.Error())
			// }
			// if n == nil {
			// 	fmt.Println("sadd cmd failed!")
			// }
		}

		this.conn.Close()
		this.conn = nil
		this.isVerifyed = false
		this.state = state_del
		removeUidMap(this.uid)
		removeClient(this.clientAddr)
	}

	this.closeTimer()
}

func (this *Client) closeTimer() {
	if this.timer != nil {
		this.timer.Reset(time.Millisecond * 1)
		this.timer = nil
	}
}

func (this *Client) waitForLogin() {
	this.state = state_connected
	this.timer = time.NewTimer(time.Second * 30)

	select {
	case <-this.timer.C:
		this.lock.Lock()
		if !this.isVerifyed {
			this.Close()
		}
		this.lock.Unlock()
	}
	fmt.Println("waitForLogin end")
}

func (this *Client) waitForAppVerify() {
	this.state = state_connected
	this.timer = time.NewTimer(time.Second * 30)

	select {
	case <-this.timer.C:
		this.lock.Lock()
		if !this.isAppVerifyed {
			this.Close()
		}
		this.lock.Unlock()
	}
	fmt.Println("waitForAppVerify end")
}

func (this *Client) verifyAppLogin() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.state == state_verify {
		this.isAppVerifyed = true
		this.closeTimer()
		go this.startTick()
		return true
	} else {
		this.Close()
		return false
	}
}

func (this *Client) startTick() {
	this.timer = time.NewTimer(time.Second * 60)
	for {
		select {
		case <-this.timer.C:
			fmt.Println("countTimeOut++")
			this.countTimeOut++
			if this.countTimeOut >= 2 {
				if this.timer != nil {
					this.timer.Stop()
				}
				this.Close()
				return
			}
			if this.timer != nil {
				this.timer.Reset(time.Second * 60)
			}
		case <-this.tickChan:
			this.countTimeOut = 0
			if this.timer != nil {
				this.timer.Reset(time.Second * 60)
			}
		}

		if this.state == state_none {
			break
		}
	}
	fmt.Println("tick end")
}

func (this *Client) ParseHeader(data []byte) int {
	size := Int(data)
	//fmt.Println("header size :", size)
	//p.conn.Send(data)
	return size
}

func (this *Client) ParseMsg(data []byte) {
	//fmt.Println("client:", this.conn.ConnAddr(), "say:", String(data))
	msgid := Uint16(data)
	//fmt.Println("msgid:", msgid)
	if this.isVerifyed && !this.isAppLogin {
		this.tickChan <- 1
	} else if this.isVerifyed && this.isAppLogin && this.isAppVerifyed {
		this.tickChan <- 1
	} else if msgid != MsgId_ReqLogin && msgid != MsgId_ReqAppLogin {
		//if not logined, do not response to any msg
		return
	}

	switch msgid {
	case MsgId_ReqLogin:
		if this.isVerifyed {
			//if had logined, do nothing
			return
		}
		uid := Uint64(data[2:10])
		password := data[10:]
		ret := new(MsgRetLogin)
		code := gDataManager.checkLogin(uid, string(password))
		if code == ERR_NONE {
			this.state = state_logined
			this.uid = uid
			this.password = string(password)
			ret.Result = uint16(ERR_NONE)
			addUidMap(uid, this)
			//copy(ret.IP[0:], []byte("127.0.0.1"))
			//ret.Port = 9090
			ok := gDataManager.setUserOnline(uid)
			if !ok {
				ret.Result = uint16(ERR_REDIS)
			} else {
				this.lock.Lock()
				this.isVerifyed = true
				this.tickChan = make(chan int, 2)
				this.closeTimer()
				go this.startTick()
				this.lock.Unlock()
				fmt.Println("addr:" + this.conn.ConnAddr() + " logined success")
			}
		} else {
			ret.Result = uint16(code)
			this.verfiycount++

			if this.verfiycount < 5 {
				this.timer.Reset(time.Second * 30)
			} else {
				this.Close()
			}
		}
		ret.MsgId = MsgId_RetLogin
		this.send(Bytes(ret))
	case MsgId_ReqAppLogin:
		if this.isVerifyed {
			//if had logined, do nothing
			return
		}

		uid := Uint64(data[2:10])
		password := string(data[10:42])
		appname := string(data[42:])
		//check app login
		//if login ok, then wait for app server verify
		code := gDataManager.checkAppLogin(uid, password, appname)
		ret := new(MsgRetAppLogin)
		ret.MsgId = MsgId_RetAppLogin
		if code == ERR_NONE {
			this.state = state_verify
			this.uid = uid
			this.password = string(password)
			this.isAppLogin = true
			this.appName = appname
			ret.Result = uint16(ERR_NONE)
			addUidMap(uid, this)

			//copy(ret.IP[0:], []byte("127.0.0.1"))
			//ret.Port = 9090
			ok := gDataManager.setUserOnline(uid)
			if !ok {
				ret.Result = uint16(ERR_REDIS)
			} else {
				uuid := Authcode(String(uid)+":"+appname, "ENCODE")
				code = gDataManager.setAppVerifyData(uuid, uid)
				if code != ERR_NONE {
					ret.Result = uint16(ERR_REDIS)
				} else {
					ret.Count = byte(len(uuid))
					ret.Uuid = []byte(uuid)
					this.lock.Lock()
					this.isVerifyed = true
					this.tickChan = make(chan int, 2)
					this.closeTimer()
					go this.waitForAppVerify()
					this.lock.Unlock()
					fmt.Println("addr:" + this.conn.ConnAddr() + " logined success")
				}
			}
		} else {
			//else ret login failed.
			ret.Result = uint16(code)
			this.verfiycount++

			if this.verfiycount < 5 {
				this.timer.Reset(time.Second * 30)
			} else {
				this.Close()
			}
		}
		this.send(Bytes(ret))
	case MsgId_Tick:
		ret := new(MsgTick)
		ret.MsgId = MsgId_Tick
		this.send(Bytes(ret))
	case MsgId_ReqLoginOut:
		ret := new(MsgRetLoginOut)
		ret.Result = 1
		ret.MsgId = MsgId_ReqRetLoginOut
		this.send(Bytes(ret))
		this.state = state_logouted
	case MsgId_Echo:
		// ret := new(Echo)
		// ret.MsgId = MsgId_Echo
		// ret.Data = data[2:]
		this.send(data)
	default:
		fn, ok := msgProcesserMap[msgid]
		if ok {
			fn(this, data)
		} else {
			fmt.Println("unknown msgid:", msgid)
		}
	}
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
	if this.state == state_logouted {
		this.Close()
	}
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
	key := "abc"
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
