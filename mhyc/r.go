package mhyc

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"time"
)

var RoleLoadWait = make(chan struct{})
var UserBagWait = make(chan struct{})

var CLI *Connect

var PCK = map[uint16]Handle{
	2:     &S2CLogin{},
	3:     &S2CRoleInfo{},
	12:    &S2CNotice{},
	14:    &S2CRespect{},
	23:    &Pong{},
	37:    &S2CNewStory{},
	48:    &S2CRoleTask{},
	51:    &S2CChangeMap{},
	59:    &S2CMonsterEnterMap{},
	60:    &S2CMonsterLeaveMap{},
	100:   &S2CPrizeReport{},
	101:   &S2CBattlefieldReport{},
	104:   &S2CStageFight{},
	110:   &S2CStagePrize{},
	119:   &S2CStageDraw{},
	137:   &S2CGetVipDayGift{},
	403:   &S2CNewChatMsg{},
	433:   &S2CShopBuy{},
	501:   &S2CUserBag{},
	520:   &S2CBagChange{},
	524:   &ItemFly{},
	575:   &S2CAutoMeltGain{},
	605:   &S2CBossPersonalSweep{},
	704:   &S2CGetTaskPrize{},
	714:   &S2CGetHistoryTaskPrize{},
	1001:  &S2CServerTime{},
	1111:  &S2CMultiBossInfo{},
	1542:  &S2CGetActTimestamp{},
	11032: &S2CGetActXunBaoInfo{},
	11036: &S2CActXunBaoDraw{},
	12012: &S2CActGiftNewReceive{},
	12152: &S2CGetActTask{},
	15030: &S2CHomeBossInfo{},
	21001: &S2CRedState{},
	22013: &S2CRealmTask{},
	22303: &S2CSign{},
	22406: &S2CLifeCardDayPrize{},
	22572: &S2CClimbingTowerEnter{},
	22576: &S2CClimbingTowerFight{},
	22628: &S2CWeddingInsInvite{},
	22632: &S2CWeddingInsInviteAck{},
	22641: &S2CWeddingInsReport{},
	22731: &S2CGetPetAMergeInfo{},
	23102: &S2CPlayerPractice{},
	25796: &S2CBeasts{},
	26232: &S2CXsdBossInfo{},
	27002: &S2CGetAllEquipData{},
	27004: &S2CGetEquipData{},
	27152: &S2CHuanLingList{},
	27802: &S2CGetHeroList{},
	28602: &S2CGetAlienData{},
	29504: &S2CSLGetData{},
	52227: &S2CYJInfo{},
	55102: &S2CZSStateInfo{},
}

type Handle interface {
	Message([]byte)
}

func (x *S2CAutoMeltGain) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][AutoMeltGain] items=%v", x.Items)
}

func (x *S2CBattlefieldReport) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BattlefieldReport] win=%v report=%v", x.Win, x)
}

func (x *S2CClimbingTowerFight) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("[22576][ClimbingTowerFight] err=%v", err)
		return
	}
	log.Printf("[22576][ClimbingTowerFight] tag=%v TowerType=%v", x.Tag, x.TowerType)
	_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: x.TowerType})
	return
}

func (x *S2CClimbingTowerEnter) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("[22572][ClimbingTowerEnter] err=%v", err)
		return
	}
	log.Printf("[22572][ClimbingTowerEnter] tag=%v TowerType=%v", x.Tag, x.TowerType)
	if x.Tag != 0 {
		return
	}
	_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: x.TowerType})
	return
}

func (x *S2CMonsterEnterMap) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("[59][MonsterEnterMap] %v", err)
		return
	}
	log.Printf("[59][MonsterEnterMap] %v", x)
	return
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

func (x *S2CMultiBossInfo) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [MultiBossInfo] %v", err)
		return
	}
	log.Printf("recv: [MultiBossInfo] %v", x)
	return
}

func (x *ItemFly) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	for _, item := range x.Item {
		var dat *ItemData
		if val, ok := UserBag.Load(item.IId); ok {
			dat = val.(*ItemData)
		} else {
			dat = &ItemData{}
		}
		dat.N = dat.N + item.N
		UserBag.Store(item.IId, dat)
	}
	log.Printf("[S][ItemFly] %v", x)
	return
}

func (x *S2CBagChange) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [BagChange] %v", err)
		return
	}
	log.Printf("recv: [BagChange] %v", x)
	return
}

func (x *S2CZSStateInfo) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [ZSStateInfo] %v", err)
		return
	}
	log.Printf("recv: [ZSStateInfo] %v", x)
	return
}

func (x *S2CXsdBossInfo) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [XsdBossInfo] %v", err)
		return
	}
	log.Printf("recv: [XsdBossInfo] %v", x)
	return
}

func (x *S2CRedState) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [RedState] %v", err)
		return
	}
	log.Printf("recv: [RedState] %v", x)
	return
}

func (x *S2CHomeBossInfo) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [HomeBossInfo] %v", err)
		return
	}
	log.Printf("recv: [HomeBossInfo] %v", x)
	return
}

func (x *S2CNotice) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [Notice] %v", err)
		return
	}
	log.Printf("recv: [Notice] %v", x)
	return
}

func (x *S2CGetActTask) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetActTask] %v", err)
		return
	}
	for _, t := range x.Task {
		_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 6, Multi: 1, TaskId: t.Id})
	}
	log.Printf("recv: [GetActTask] %v", x)
	return
}

func (x *S2CGetActXunBaoInfo) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetActXunBaoInfo] %v", err)
		return
	}
	log.Printf("recv: [GetActXunBaoInfo] %v", x)
	return
}

func (x *S2CStagePrize) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [StagePrize] %v", err)
		return
	}
	log.Printf("recv: [StagePrize] %v", x)
	return
}

func (x *S2CUserBag) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("[501][CUserBag] %v", err)
		return
	}
	for _, item := range x.Bag.Items {
		UserBag.Store(item.IId, item)
	}
	log.Printf("[501][CUserBag] %v", x)
	UserBagWait <- struct{}{}
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

func (x *S2CBossPersonalSweep) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [BossPersonalSweep] %v", err)
		return
	}
	log.Printf("recv: [BossPersonalSweep] %v", x)
	return
}

func (x *S2CGetTaskPrize) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetTaskPrize] %v", err)
		return
	}
	log.Printf("recv: [GetTaskPrize] %v", x)
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

func (x *S2CShopBuy) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [ShopBuy] %v", err)
		return
	}
	log.Printf("recv: [ShopBuy] %v", x)
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

func (x *S2CStageDraw) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [StageDraw] %v", err)
		return
	}
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

var RIF, _ = os.OpenFile("./role_info_"+time.Now().Format("2006.01.02")+".txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

// Message S2CRoleInfo 角色信息
func (x *S2CRoleInfo) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [RoleInfo] %v", err)
		return
	}
	for _, a := range x.A {
		if _, ok := AttrType[a.K]; !ok {
			continue
		}
		log.Printf("[RoleInfo] %s\t%v", AttrType[a.K], a.V)
		_, _ = RIF.WriteString(fmt.Sprintf("%s\t%v\n", AttrType[a.K], a.V))
		RoleInfo.Store(AttrType[a.K], a.V)
	}
	for _, b := range x.B {
		if _, ok := AttrType[b.K]; !ok {
			continue
		}
		log.Printf("[RoleInfo] %s\t%v", AttrType[b.K], b.V)
		_, _ = RIF.WriteString(fmt.Sprintf("%s\t%v\n", AttrType[b.K], b.V))
		RoleInfo.Store(AttrType[b.K], b.V)
	}
	RoleLoadWait <- struct{}{}
	return
}

func (x *Pong) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [Pong] %v", err)
		return
	}
	log.Println("Pong")
	return
}

func (x *S2CChangeMap) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [ChangeMap] %v", err)
		return
	}
	log.Printf("ChangeMap: id=%d x=%d y=%d", x.MapId, x.X, x.Y)
	return
}

func (x *S2CActGiftNewReceive) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [ActGiftNewReceive] %v", err)
		return
	}
	log.Printf("ActGiftNewReceive: tag=%d gid=%d aid=%d", x.Tag, x.Gid, x.Aid)
	return
}

func (x *S2CRespect) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [Respect] %v", err)
		return
	}
	log.Printf("Respect: tag=%d type=%d prize=%v", x.Tag, x.Type, x.Prize)
	return
}

func (x *S2CServerTime) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [ServerTime] %v", err)
		return
	}

	log.Printf("ServerTime: %s", time.Unix(x.T, 0).Local().Format("2006-01-02 15:04:05"))
	return
}

func (x *S2CNewChatMsg) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [NewChatMsg] %v", err)
		return
	}
	log.Printf("NewChatMsg: [%s] %s", x.Chatmessage.SenderNick, x.Chatmessage.Content)
	return
}

func (x *S2CRoleTask) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [RoleTask] %v", err)
		return
	}
	log.Printf("[RoleTask] %v", x)
	return
}

func (x *S2CGetVipDayGift) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetVipDayGift] %v", err)
		return
	}
	log.Printf("GetVipDayGift: tag=%d", x.Tag)
	return
}

func (x *S2CActXunBaoDraw) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [ActXunBaoDraw] %v", err)
		return
	}
	log.Printf("ActXunBaoDraw: tag=%d %v", x.Tag, x)
	return
}

func (x *S2CLifeCardDayPrize) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [LifeCardDayPrize] %v", err)
		return
	}
	log.Printf("LifeCardDayPrize: tag=%d", x.Tag)
	return
}

func (x *S2CSign) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [Sign] %v", err)
		return
	}
	log.Printf("Sign: tag=%d", x.Tag)
	return
}

func (x *S2CStageFight) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [StageFight] %v", err)
		return
	}
	return
}

func (x *S2CPrizeReport) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [PrizeReport] %v", err)
		return
	}
	log.Printf("PrizeReport: %v", x)
	return
}

func (x *S2CGetHistoryTaskPrize) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetHistoryTaskPrize] %v", err)
		return
	}
	return
}
