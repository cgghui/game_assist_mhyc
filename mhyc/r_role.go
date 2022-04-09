package mhyc

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"sync"
	"time"
)

////////////////////////////////////////////////////////////

// WareHouseReceiveItem 将仓库内的物品转至背包
// 5 寻宝
func (c *Connect) WareHouseReceiveItem(id int32) error {
	body, err := proto.Marshal(&C2SWareHouseReceiveItem{WhId: id})
	if err != nil {
		return err
	}
	log.Printf("[C][WareHouseReceiveItem] wh_id=%v", id)
	return c.send(27305, body)
}

// EndFight 结束战斗
func (c *Connect) EndFight(r *S2CBattlefieldReport) error {
	body, err := proto.Marshal(&C2SEndFight{Idx: r.Idx})
	if err != nil {
		return err
	}
	log.Printf("[C][EndFight] idx=%v win=%v", r.Idx, r.Win)
	return c.send(102, body)
}

// StartFight 开始战斗
func (c *Connect) StartFight(f *C2SStartFight) error {
	body, err := proto.Marshal(f)
	if err != nil {
		return err
	}
	log.Printf("[C][StartFight] id=%v type=%v", f.Id, f.Type)
	return c.send(61, body)
}

// RoleInfo 角色信息
func (c *Connect) RoleInfo() error {
	body, err := proto.Marshal(&C2SRoleInfo{})
	if err != nil {
		return err
	}
	log.Println("[C][角色信息] >>>")
	return c.send(49, body)
}

var RIF, _ = os.OpenFile("./role_info_"+time.Now().Format("2006.01.02")+".txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

////////////////////////////////////////////////////////////

func (x *S2CRoleInfo) ID() uint16 {
	return 3
}

// Message S2CRoleInfo 角色信息
func (x *S2CRoleInfo) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("[S][角色信息] %v", err)
		return
	}
	for _, a := range x.A {
		if _, ok := AttrType[a.K]; !ok {
			continue
		}
		log.Printf("[S][角色信息] %s\t%v", AttrType[a.K], a.V)
		_, _ = RIF.WriteString(fmt.Sprintf("%s\t%v\n", AttrType[a.K], a.V))
		RoleInfo.Set(AttrType[a.K], a.V)
	}
	for _, b := range x.B {
		if _, ok := AttrType[b.K]; !ok {
			continue
		}
		log.Printf("[S][角色信息] %s\t%v", AttrType[b.K], b.V)
		_, _ = RIF.WriteString(fmt.Sprintf("%s\t%v\n", AttrType[b.K], b.V))
		RoleInfo.Set(AttrType[b.K], b.V)
	}
	return
}

////////////////////////////////////////////////////////////

func (x *S2CChangeMap) ID() uint16 {
	return 51
}

// Message S2CChangeMap Code:51
func (x *S2CChangeMap) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][切换地图] tag=%v id=%d x=%d y=%d", x.Tag, x.MapId, x.X, x.Y)
}

////////////////////////////////////////////////////////////

func (x *S2CPlayerEnterMap) ID() uint16 {
	return 57
}

// Message S2CPlayerEnterMap Code:57
func (x *S2CPlayerEnterMap) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][玩家进入地图] user_id=%v x=%v y=%v", x.UserId, x.X, x.Y)
}

////////////////////////////////////////////////////////////

func (x *S2CPlayerLeaveMap) ID() uint16 {
	return 58
}

// Message S2CPlayerEnterMap Code:58
func (x *S2CPlayerLeaveMap) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][玩家离开地图] user_id=%v", x.UserId)
}

////////////////////////////////////////////////////////////

func (x *S2CBattlefieldReport) ID() uint16 {
	return 101
}

// Message S2CBattlefieldReport Code:101
func (x *S2CBattlefieldReport) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	_ = CLI.EndFight(x)
	log.Printf("[S][BattlefieldReport] win=%v idx=%v", x.Win, x.Idx)
}

////////////////////////////////////////////////////////////

func (x *S2CWareHouseReceiveItem) ID() uint16 {
	return 27306
}

// Message S2CWareHouseReceiveItem Code:27306
func (x *S2CWareHouseReceiveItem) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][WareHouseReceiveItem] tag=%v wh_id=%v", x.Tag, x.WhId)
}

////////////////////////////////////////////////////////////

func (x *S2CStartFight) ID() uint16 {
	return 62
}

// Message S2CStartFight Code:62
func (x *S2CStartFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][StartFight] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CMonsterEnterMap) ID() uint16 {
	return 59
}

// Message S2CMonsterEnterMap Code:59
func (x *S2CMonsterEnterMap) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][MonsterEnterMap] id=%v", x.Id)
}

////////////////////////////////////////////////////////////

func (x *Pong) ID() uint16 {
	return 23
}

// Message Pong Code:23
func (x *Pong) Message(_ []byte) {
	log.Printf("[S][Pong]")
}

////////////////////////////////////////////////////////////

func (x *S2CRoleTask) ID() uint16 {
	return 48
}

// Message S2CRoleTask Code:48
func (x *S2CRoleTask) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][RoleTask] %v", x)
}

////////////////////////////////////////////////////////////

func (x *S2CServerTime) ID() uint16 {
	return 1001
}

func (x *S2CServerTime) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ServerTime] %s", time.Unix(x.T, 0).Local().Format("2006-01-02 15:04:05"))
	return
}

////////////////////////////////////////////////////////////

func (x *S2CRedState) ID() uint16 {
	return 21001
}

func (x *S2CRedState) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][RedState] list=%v", x.List)
	return
}

////////////////////////////////////////////////////////////

func (x *S2CNotice) ID() uint16 {
	return 12
}

func (x *S2CNotice) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][Notice] %v", x)
	return
}

////////////////////////////////////////////////////////////

func (x *S2CNewChatMsg) ID() uint16 {
	return 403
}

func (x *S2CNewChatMsg) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][NewChatMsg] [%s] %s", x.Chatmessage.SenderNick, x.Chatmessage.Content)
	return
}

////////////////////////////////////////////////////////////

var RoleInfo = &roleInfo{
	s: &sync.Map{},
}

type roleInfo struct{ s *sync.Map }

type roleValue struct {
	val interface{}
}

func (r *roleValue) Int64() int64 {
	if r == nil {
		return 0
	}
	return r.val.(int64)
}

func (r *roleValue) String() string {
	if r == nil {
		return ""
	}
	return r.val.(string)
}

func (r *roleInfo) Set(name string, val interface{}) {
	r.s.Store(name, val)
}

func (r *roleInfo) Has(name string) bool {
	_, ok := r.s.Load(name)
	return ok
}

func (r *roleInfo) Get(name string) *roleValue {
	ret, ok := r.s.Load(name)
	if !ok {
		return nil
	}
	return &roleValue{val: ret}
}

func (r *roleInfo) Wait(id int32, timeout time.Duration) *roleValue {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	tm := time.NewTimer(ms100)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tm.C:
			_ = CLI.RoleInfo()
			if err = Receive.WaitWithContext(ctx, &S2CRoleInfo{}); err != nil {
				return nil
			}
			ret, ok := r.s.Load(id)
			if !ok {
				tm.Reset(ms100)
				continue
			}
			return &roleValue{val: ret}
		}
	}
}
