package centity

import "fmt"

const (
	SMALL_MSG_ID_LOGIN = iota
)

func init() {
	msgProcesser[BIG_MSG_ID_LOGIN] = make([]func(IEntity, []byte))
	msgProcesser[BIG_MSG_ID_LOGIN][SMALL_MSG_ID_LOGIN] = onLogin
}

func onLogin(entity IEntity, data []byte) {
	uid := Uint64(data)
	password := string(data[8:])

	code := gDataManager.checkLogin(uid, password)

	ret := new(MsgRetLogin)
	ret.MsgId = MsgId_RetLogin

	if code == ERR_NONE {
		code = gDataManager.setUserOnline(uid)
		if code == ERR_NONE {
			Manager().CreateNullEntity(entity)
			fmt.Println("addr:" + this.conn.ConnAddr() + " logined success")
		}
	}
	ret.Result = uint16(code)
	result = code != ERR_NONE
	this.send(Bytes(ret))
}
