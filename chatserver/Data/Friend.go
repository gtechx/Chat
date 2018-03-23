package data

import (
	//"errors"
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
	. "github.com/nature19862001/base/common"
)

var defaultGroupName string = "我的好友"
var userOnlineKeyName string = "user:online"

func (this *RedisDataManager) addFriendReq(uid, fuid uint64, group string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	//check if in black list
	// ret, err := conn.Do("SISMEMBER", "black:"+String(fuid), uid)

	// if err != nil {
	// 	fmt.Println("addFriend error:", err.Error())
	// 	return ERR_REDIS
	// }

	if this.isInBlacklist(fuid, uid, conn) {
		return ERR_FRIEND_IN_BLACKLIST
	}

	//check if group exists
	ret, err := conn.Do("SISMEMBER", "fgroup:"+String(uid), group)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	if Bool(ret) != true {
		//return ERR_FRIEND_GROUP_NOT_EXIST
		group = defaultGroupName
	}

	//check if group exists
	_, err = conn.Do("HSET", "freq", String(uid)+":"+String(fuid), group)

	if err != nil {
		fmt.Println("addFriendReq error:", err.Error())
		return -1
	}

	return ERR_NONE
}

func (this *RedisDataManager) addFriend(uid, fuid uint64, group string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	if this.isInBlacklist(fuid, uid, conn) {
		return ERR_FRIEND_IN_BLACKLIST
	}

	//check if friend is exist
	ret, err := conn.Do("SISMEMBER", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	if Bool(ret) {
		return ERR_FRIEND_EXIST
	}

	//check if fuid user is exists
	ret, err = conn.Do("EXISTS", fuid)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	if Bool(ret) {
		return ERR_USER_NOT_EXIST
	}

	//check if friend count is max
	ret, err = conn.Do("HGET", uid, "maxfriends")

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	maxfriends := Int(ret)

	ret, err = conn.Do("SCARD", "friend:"+String(uid))

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	curfriendcount := Int(ret)

	if curfriendcount >= maxfriends {
		return ERR_FRIEND_MAX
	}

	//check if group exists
	ret, err = conn.Do("SISMEMBER", "fgroup:"+String(uid), group)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	if Bool(ret) != true {
		//return ERR_FRIEND_GROUP_NOT_EXIST
		group = defaultGroupName
	}

	//get fuid's verify type, default is VERIFY_TYPE_NEED_AGREE
	vtype := VERIFY_TYPE_NEED_AGREE
	ret, err = conn.Do("HGET", fuid, "vtype")

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	if ret != nil {
		vtype = Int(ret)
	}

	//check if this is a agree message
	ret, err = conn.Do("HGET", "freq", String(fuid)+":"+String(uid))

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	if ret == nil {
		//if add friend actively, then check fuid's vtype
		if vtype == VERIFY_TYPE_NEED_AGREE {
			//need friend request to fuid
			return ERR_FRIEND_ADD_NEED_REQ
		}

		if vtype == VERIFY_TYPE_REFUSE_ALL {
			return ERR_FRIEND_ADD_REFUSE_ALL
		}
	}

	fuidgroup := String(ret)

	//get friend count of group
	ret, err = conn.Do("HGET", "friend:group:"+String(uid), group)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	groupuidcount := 0

	if ret != nil {
		groupuidcount = Int(ret)
	}

	//get friend's friend count of fuidgroup
	ret, err = conn.Do("HGET", "friend:group:"+String(fuid), fuidgroup)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	fuidgroupuidcount := 0

	if ret != nil {
		fuidgroupuidcount = Int(ret)
	}

	conn.Send("MULTI")
	conn.Send("HSET", "friend:"+String(uid), fuid, group+":"+String(groupuidcount))
	conn.Send("HSET", "friend:group:"+String(uid), group+":"+String(groupuidcount), fuid)
	conn.Send("HINCRBY", "friend:group:"+String(uid), group, 1)

	conn.Send("HSET", "friend:"+String(fuid), uid, fuidgroup+":"+String(fuidgroupuidcount))
	conn.Send("HSET", "friend:group:"+String(fuid), fuidgroup+":"+String(fuidgroupuidcount), uid)
	conn.Send("HINCRBY", "friend:group:"+String(fuid), fuidgroup, 1)

	conn.Send("HDEL", "freq", String(fuid)+":"+String(uid))
	_, err = conn.Do("EXEC")

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	return ERR_NONE
}

func (this *RedisDataManager) deleteFriend(uid, fuid uint64) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	//get group
	ret, err := conn.Do("HGET", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("deleteFriend error:", err.Error())
		return ERR_REDIS
	}

	if ret == nil {
		return ERR_FRIEND_NOT_EXIST
	}

	groupstr := String(ret)
	groupstrarr := strings.Split(groupstr, ":")
	group := groupstrarr[0]

	conn.Send("MULTI")
	conn.Send("HDEL", "friend:"+String(uid), fuid)
	conn.Send("HDEL", "friend:group:"+String(uid), groupstr)
	conn.Send("HINCRBY", "friend:group:"+String(uid), group, -1)
	//conn.Send("SREM", "fgroup:"+String(uid)+":"+group, fuid)
	_, err = conn.Do("EXEC")

	if err != nil {
		fmt.Println("deleteFriend error:", err.Error())
		return ERR_REDIS
	}

	return ERR_NONE
}

func (this *RedisDataManager) getFriendList(uid uint64) (map[string][]uint64, int) {
	conn := this.redisPool.Get()
	defer conn.Close()

	//
	ret, err := conn.Do("HGETALL", "friend:group:"+String(uid))

	if err != nil {
		fmt.Println("getFriendList error:", err.Error())
		return map[string][]uint64{}, ERR_REDIS
	}

	if ret == nil {
		return map[string][]uint64{}, ERR_NONE
	}

	dataarr, err := redis.Strings(ret, err)
	retdata := map[string][]uint64{}

	for i := 0; i < len(dataarr); i += 2 {
		groupstr := dataarr[i]
		groupstrarr := strings.Split(groupstr, ":")
		if len(groupstrarr) <= 1 {
			continue
		}
		groupname := groupstrarr[0]

		_, ok := retdata[groupname]

		fuid := Uint64(groupstrarr[1])
		if !ok {
			retdata[groupname] = make([]uint64, 0)
		}
		retdata[groupname] = append(retdata[groupname], fuid)
	}

	return retdata, ERR_NONE
}

func (this *RedisDataManager) addFriendGroup(uid uint64, groupname string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("SADD", "fgroup:"+String(uid), groupname)

	if err != nil {
		fmt.Println("addFriendGroup error:", err.Error())
		return ERR_REDIS
	}

	if Int(ret) != 1 {
		return ERR_FRIEND_GROUP_EXIST
	}

	return ERR_NONE
}

func (this *RedisDataManager) deleteFriendGroup(uid uint64, groupname string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HGET", "friend:group:"+String(uid), groupname)

	if err != nil {
		fmt.Println("deleteFriendGroup error:", err.Error())
		return ERR_REDIS
	}

	if ret != nil && Int(ret) > 0 {
		return ERR_FRIEND_GROUP_USER_NOT_EMPTY
	}

	ret, err = conn.Do("SREM", "fgroup:"+String(uid), groupname)

	if err != nil {
		fmt.Println("deleteFriendGroup error:", err.Error())
		return ERR_REDIS
	}

	// if Int(ret) != 1 {
	// 	return 0
	// }

	return ERR_NONE
}

func (this *RedisDataManager) getGroupOfFriend(uid, fuid uint64) (string, int) {
	conn := this.redisPool.Get()
	defer conn.Close()

	//check if friend is exist
	ret, err := conn.Do("HGET", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("getGroupOfFriend error:", err.Error())
		return "", ERR_REDIS
	}

	if ret == nil {
		return "", ERR_FRIEND_NOT_EXIST
	}

	return String(ret), ERR_NONE
}

func (this *RedisDataManager) moveFriendToGroup(uid, fuid uint64, destgroup string) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	//check if friend is exist
	ret, err := conn.Do("HEXISTS", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("addFriend error:", err.Error())
		return ERR_REDIS
	}

	if !Bool(ret) {
		return ERR_FRIEND_NOT_EXIST
	}

	//check if group exists
	// ret, err = conn.Do("SISMEMBER", "fgroup:"+String(uid), srcgroup)

	// if err != nil {
	// 	fmt.Println("moveFriendToGroup error:", err.Error())
	// 	return ERR_REDIS
	// }

	// if Bool(ret) != true {
	// 	return ERR_FRIEND_GROUP_NOT_EXIST
	// }

	ret, err = conn.Do("SISMEMBER", "fgroup:"+String(uid), destgroup)

	if err != nil {
		fmt.Println("moveFriendToGroup error:", err.Error())
		return ERR_REDIS
	}

	if Bool(ret) != true {
		return ERR_FRIEND_GROUP_NOT_EXIST
	}

	//get group fuid current in
	ret, err = conn.Do("HGET", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("moveFriendToGroup error:", err.Error())
		return ERR_REDIS
	}

	groupstr := String(ret)
	groupstrarr := strings.Split(groupstr, ":")
	curgroup := groupstrarr[0]

	if curgroup == destgroup {
		//if already in destgroup, return
		return ERR_NONE
	}

	//get destgroup count
	ret, err = conn.Do("HGET", "friend:group:"+String(uid), destgroup)

	if err != nil {
		fmt.Println("moveFriendToGroup error:", err.Error())
		return ERR_REDIS
	}

	destgroupcount := 0

	if ret != nil {
		destgroupcount = Int(ret)
	}

	conn.Send("MULTI")
	//remove from current group
	conn.Send("HDEL", "friend:"+String(uid), fuid)
	conn.Send("HDEL", "friend:group:"+String(uid), groupstr)
	conn.Send("HINCRBY", "friend:group:"+String(uid), curgroup, -1)

	//add to dest group
	conn.Send("HSET", "friend:"+String(uid), fuid, destgroup+":"+String(destgroupcount))
	conn.Send("HSET", "friend:group:"+String(uid), destgroup+":"+String(destgroupcount), fuid)
	conn.Send("HINCRBY", "friend:group:"+String(uid), destgroup, 1)

	// conn.Send("SADD", "fgroup:"+String(uid)+":"+destgroup, fuid)
	// conn.Send("SREM", "fgroup:"+String(uid)+":"+srcgroup, fuid)
	_, err = conn.Do("EXEC")

	if err != nil {
		fmt.Println("deleteFriend error:", err.Error())
		return ERR_REDIS
	}

	// ret, err = conn.Do("SADD", "fgroup:"+String(uid)+":"+destgroup, fuid)

	// if err != nil {
	// 	fmt.Println("moveFriendToGroup error:", err.Error())
	// 	return -1
	// }

	// if Int(ret) != 1 {
	// 	return 0
	// }

	// ret, err = conn.Do("SREM", "fgroup:"+String(uid)+":"+srcgroup, fuid)

	// if err != nil {
	// 	fmt.Println("moveFriendToGroup error:", err.Error())
	// 	return -1
	// }

	// if Int(ret) != 1 {
	// 	return 0
	// }

	return ERR_NONE
}

func (this *RedisDataManager) banFriend(uid, fuid uint64) {

}

func (this *RedisDataManager) unBanFriend(uid, fuid uint64) {

}

func (this *RedisDataManager) isFriend(uid, fuid uint64) bool {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HEXISTS", "friend:"+String(uid), fuid)

	if err != nil {
		fmt.Println("getFriendVerifyType error:", err.Error())
		return false
	}

	return Bool(ret)
}

func (this *RedisDataManager) setFriendVerifyType(uid uint64, vtype byte) int {
	conn := this.redisPool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", uid, "vtype", vtype)

	if err != nil {
		fmt.Println("setFriendVerifyType error:", err.Error())
		return ERR_REDIS
	}

	return ERR_NONE
}

func (this *RedisDataManager) getFriendVerifyType(uid uint64) byte {
	conn := this.redisPool.Get()
	defer conn.Close()

	ret, err := conn.Do("HGET", uid, "vtype")

	if err != nil {
		fmt.Println("getFriendVerifyType error:", err.Error())
		return VERIFY_TYPE_ERR
	}

	return Byte(ret)
}
