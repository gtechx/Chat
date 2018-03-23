package entity

import "github.com/nature19862001/base/gtnet"

var BigMsgIDCounter uint8 = 0
var SmallMsgIDCounter uint8 = 0

const (
	BIG_MSG_ID_LOGIN uint8 = iota + 10000
)

const BIG_MSG_ID_START = BIG_MSG_ID_LOGIN

var msgProcesser [][]func(IEntity, []byte) = [][]func(IEntity, []byte){}
var bigMsgProcesser []func(IEntity, []byte) = []func(IEntity, []byte){}
var smallMsgProcesser []func(IEntity, []byte) = []func(IEntity, []byte){}

func init() {
}

type EntityManager struct {
	nullEntityMap         map[uint64]IEntity
	userIDEntityMap       map[uint64]IEntity
	userAPPIDUIDEntityMap map[uint64]map[uint64]IEntity
	curID                 uint64

	delChan chan IEntity
	addChan chan IEntity
}

var instanceEntityManager *EntityManager

func Manager() *EntityManager {
	if instanceEntityManager == nil {
		instanceEntityManager = &EntityManager{nullEntityMap: make(map[uint64]IEntity), userIDEntityMap: make(map[uint64]IEntity), userAPPIDUIDEntityMap: make(map[uint64]map[uint64]IEntity)}
	}
	return instanceEntityManager
}

func (this *EntityManager) Initialize() {
	this.delChan = make(chan IEntity, 1024)
	this.addChan = make(chan IEntity, 1024)

	go this.userEntityProcess()
}

func (this *EntityManager) CreateNullEntity(conn gtnet.IConn) IEntity {
	this.curID++
	entity := newNullEntity(this.curID, conn)
	entity.start()
	return entity
}

// func (this *EntityManager) RemoveNullEntity(id uint64) {
// 	delete(this.nullEntityMap, id)
// }

func (this *EntityManager) CreateEntity(etype int, entity IEntity) IEntity {
	switch etype {
	case TYPE_USER:
		newentity := newUserEntity(entity)
		this.addChan <- newentity
		return newentity
	}

	return nil
}

func (this *EntityManager) RemoveEntity(entity IEntity) {
	this.delChan <- entity
}

func (this *EntityManager) doAddEntity(entity IEntity) {
	oldentitymap, ok := this.userAPPIDUIDEntityMap[entity.APPID()]
	eid := entity.ID()
	uid := entity.UID()
	if !ok {
		this.userAPPIDUIDEntityMap[entity.APPID()] = make(map[uint64]IEntity)
	} else {
		oldentity, ok := oldentitymap[entity.UID()]

		if ok {
			oldentityuid := entity.UID()
			oldentityid := entity.ID()

			if ok {
				oldentity.ForceOffline()
				delete(oldentitymap, oldentityuid)
				delete(this.userIDEntityMap, eid)
			}
		}
	}

	userIDEntityMap[entity.ID()] = entity
	userAPPIDUIDEntityMap[entity.APPID()][entity.UID()] = entity
	entity.start()
}

func (this *EntityManager) userEntityProcess() {
	select {
	case entity := <-this.addChan:
		if entity.Conn() != nil {
			this.doAddEntity(entity)
		}
	case entity := <-this.delChan:
		if entity.Conn() != nil {
			this.doRemoveEntity(entity)
		}
	}
}

func (this *EntityManager) doRemoveEntity(entity IEntity) {
	entity, ok := this.userIDEntityMap[id]

	if ok {
		delete(this.userIDEntityMap, id)
		delete(this.userAPPIDUIDEntityMap[entity.APPID()], entity.UID())
	}
}

// func (this *EntityManager) RemoveUserEntityByID(id uint64) {
// 	entity, ok := this.userIDEntityMap[id]

// 	if ok {
// 		delete(this.userIDEntityMap, id)
// 		delete(this.userAPPIDUIDEntityMap[entity.APPID()], entity.UID())
// 	}
// }

// func (this *EntityManager) RemoveUserEntityByAPPIDUID(appid, uid uint64) {
// 	entitymap, ok := this.userAPPIDUIDEntityMap[appid]

// 	if ok {
// 		entity, ok := entitymap[uid]

// 		if ok {
// 			delete(this.userIDEntityMap, entity.ID())
// 			delete(entitymap, uid)
// 		}
// 	}
// }
