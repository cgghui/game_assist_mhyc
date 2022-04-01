package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
)

var CLI *Connect
var Fight = &sync.Mutex{}

//var PCK = map[uint16]HandleMessage{
//	2:     &S2CLogin{},
//	12:    &S2CNotice{},
//	23:    &Pong{},
//	37:    &S2CNewStory{},
//	48:    &S2CRoleTask{},
//	51:    &S2CChangeMap{},
//	54:    &S2CPlayerMove{},
//	59:    &S2CMonsterEnterMap{},
//	60:    &S2CMonsterLeaveMap{},
//	66:    &S2CCheckFight{},
//	155:   &S2CRoutePath{},
//	403:   &S2CNewChatMsg{},
//	433:   &S2CShopBuy{},
//	524:   &ItemFly{},
//	575:   &S2CAutoMeltGain{},
//	605:   &S2CBossPersonalSweep{},
//	704:   &S2CGetTaskPrize{},
//	1001:  &S2CServerTime{},
//	1111:  &S2CMultiBossInfo{},
//	1542:  &S2CGetActTimestamp{},
//	11032: &S2CGetActXunBaoInfo{},
//	12012: &S2CActGiftNewReceive{},
//	12152: &S2CGetActTask{},
//	15030: &S2CHomeBossInfo{},
//	19060: &S2CSectIMSeizeReward{},
//	21001: &S2CRedState{},
//	22013: &S2CRealmTask{},
//	22303: &S2CSign{},
//	22406: &S2CLifeCardDayPrize{},
//	22572: &S2CClimbingTowerEnter{},
//	22576: &S2CClimbingTowerFight{},
//	22628: &S2CWeddingInsInvite{},
//	22632: &S2CWeddingInsInviteAck{},
//	22641: &S2CWeddingInsReport{},
//	22731: &S2CGetPetAMergeInfo{},
//	23102: &S2CPlayerPractice{},
//	25796: &S2CBeasts{},
//	26232: &S2CXsdBossInfo{},
//	27002: &S2CGetAllEquipData{},
//	27004: &S2CGetEquipData{},
//	27152: &S2CHuanLingList{},
//	27802: &S2CGetHeroList{},
//	28602: &S2CGetAlienData{},
//	29504: &S2CSLGetData{},
//	52227: &S2CYJInfo{},
//	55102: &S2CZSStateInfo{},
//}

type HandleMessage interface {
	ID() uint16
	Message([]byte)
}

func (x *S2CSectIMSeizeReward) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SectIMSeizeReward] tag=%v %v", x.Tag, x)
}

func (x *S2CCheckFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][CheckFight] tag=%v next_time=%v", x.Tag, x.NextTime)
}

func (x *S2CRoutePath) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][RoutePath] tag=%v map_id=%v point=%v", x.Tag, x.MapId, x.Points)
}

func (x *S2CPlayerMove) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][PlayerMove] userid=%v p=%v", x.UserId, x.P)
}

func (x *S2CAutoMeltGain) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][AutoMeltGain] items=%v", x.Items)
}

func (x *S2CMonsterLeaveMap) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("[60][MonsterLeaveMap] %v", err)
		return
	}
	log.Printf("[60][MonsterLeaveMap] %v", x)
	return
}

func (x *S2CWeddingInsReport) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [WeddingInsReport] %v", err)
		return
	}
	_ = CLI.WeddingInsInvite()
	log.Printf("recv: [WeddingInsReport] %v", x)
	return
}

func (x *S2CWeddingInsInviteAck) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [WeddingInsInviteAck] %v", err)
		return
	}
	log.Printf("recv: [WeddingInsInviteAck] %v", x)
	return
}

func (x *S2CWeddingInsInvite) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [WeddingInsInvite] %v", err)
		return
	}
	log.Printf("recv: [WeddingInsInvite] %v", x)
	if x.Timestamp == 0 && x.P == 0 {
		return
	}
	_ = CLI.WeddingInsFight()
}

func (x *S2CZSStateInfo) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [ZSStateInfo] %v", err)
		return
	}
	log.Printf("recv: [ZSStateInfo] %v", x)
	return
}

//func (x *S2CGetActTask) Message(data []byte) {
//	if err := proto.Unmarshal(data, x); err != nil {
//		log.Printf("recv: [GetActTask] %v", err)
//		return
//	}
//	for _, t := range x.Task {
//		_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 6, Multi: 1, TaskId: t.Id})
//	}
//	log.Printf("recv: [GetActTask] %v", x)
//	return
//}

func (x *S2CGetActXunBaoInfo) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetActXunBaoInfo] %v", err)
		return
	}
	log.Printf("recv: [GetActXunBaoInfo] %v", x)
	return
}

func (x *S2CNewStory) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [NewStory] %v", err)
		return
	}
	log.Printf("recv: [NewStory] %v", x)
	return
}

func (x *S2CSLGetData) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [SLGetData] %v", err)
		return
	}
	log.Printf("recv: [SLGetData] %v", x)
	return
}

func (x *S2CYJInfo) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [YJInfo] %v", err)
		return
	}
	log.Printf("recv: [YJInfo] %v", x)
	return
}

func (x *S2CGetAlienData) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetAlienData] %v", err)
		return
	}
	log.Printf("recv: [GetAlienData] %v", x)
	return
}

func (x *S2CGetHeroList) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetHeroList] %v", err)
		return
	}
	log.Printf("recv: [GetHeroList] %v", x)
	return
}

func (x *S2CBeasts) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [Beasts] %v", err)
		return
	}
	log.Printf("recv: [Beasts] %v", x)
	return
}

func (x *S2CGetEquipData) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetEquipData] %v", err)
		return
	}
	log.Printf("recv: [GetEquipData] %v", x)
	return
}

func (x *S2CPlayerPractice) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [日常-每日/每周] %v", err)
		return
	}
	log.Printf("recv: [日常-每日/每周] %v", x)
	return
}

func (x *S2CGetAllEquipData) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetAllEquipData] %v", err)
		return
	}
	log.Printf("recv: [GetAllEquipData] %v", x)
	return
}

func (x *S2CGetPetAMergeInfo) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetPetAMergeInfo] %v", err)
		return
	}
	log.Printf("recv: [GetPetAMergeInfo] %v", x)
	return
}

func (x *S2CGetActTimestamp) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetActTimestamp] %v", err)
		return
	}
	log.Printf("recv: [GetActTimestamp] %v", x)
	return
}

func (x *S2CRealmTask) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [RealmTask] %v", err)
		return
	}
	for _, task := range x.Tasks {
		if task.S == 1 {
			_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 24, Multi: 1, TaskId: task.Id})
		}
	}
	//log.Printf("recv: [RealmTask] %v", x)
	return
}

func (x *S2CHuanLingList) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [HuanLingList] %v", err)
		return
	}
	log.Printf("recv: [HuanLingList] %v", x)
	return
}

// Message S2CLogin 登录
func (x *S2CLogin) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [Login] %v", err)
		return
	}
	if x.UserId == 0 {
		log.Printf("recv: [Login] 无法登录，请更换【token】")
		return
	}
	log.Printf("recv: [Login] 登录成功 用户ID: %d", x.UserId)
	return
}

func (x *S2CChangeMap) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ChangeMap] id=%d x=%d y=%d", x.MapId, x.X, x.Y)
}
