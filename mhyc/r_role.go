package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"time"
)

////////////////////////////////////////////////////////////

func (c *Connect) StartMove(m *C2SStartMove) error {
	body, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	log.Printf("[C][角色移动] P0:%d P1:%d", m.P[0], m.P[1])
	return c.send(52, body)
}

func (x *S2CStartMove) ID() uint16 {
	return 53
}

// Message S2CStartMove 53
func (x *S2CStartMove) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][角色移动] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

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
		RoleInfo.Set(AttrType[a.K], a.V)
	}
	for _, b := range x.B {
		if _, ok := AttrType[b.K]; !ok {
			continue
		}
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
	log.Printf("[S][切换地图] tag=%v tag_msg=%s id=%d x=%d y=%d", x.Tag, GetTagMsg(x.Tag), x.MapId, x.X, x.Y)
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
	log.Printf("[S][WareHouseReceiveItem] tag=%v tag_msg=%s wh_id=%v", x.Tag, GetTagMsg(x.Tag), x.WhId)
}

////////////////////////////////////////////////////////////

func (x *S2CStartFight) ID() uint16 {
	return 62
}

// Message S2CStartFight Code:62
func (x *S2CStartFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][StartFight] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
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

func (x *S2CUpdateAmount) ID() uint16 {
	return 11086
}

func (x *S2CUpdateAmount) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][UpdateAmount] tag=%v tag_msg=%s act_id=%v amount=%v", x.Tag, GetTagMsg(x.Tag), x.ActId, x.Amount)
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

func (r *roleInfo) GetAll() (ret map[string]interface{}) {
	ret = make(map[string]interface{})
	r.s.Range(func(key, value interface{}) bool {
		ret[key.(string)] = value
		return true
	})
	return
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

type ActionRunHistory struct {
	Name        string        `json:"name"`
	RunningTime string        `json:"running_time"`
	TakeUpTime  time.Duration `json:"take_up_time"`
}

type ActionManage struct {
	Name   string
	Ctx    context.Context
	Cancel context.CancelFunc
	Sr     time.Time
}

func (a *ActionManage) End() {
	if len(ActionRunningHistoryList) > 200 {
		ActionRunningHistoryList = ActionRunningHistoryList[1:]
	}
	ActionRunningHistoryList = append(ActionRunningHistoryList, ActionRunHistory{
		Name:        a.Name,
		RunningTime: time.Now().Format("2006-01-02 15:04:05"),
		TakeUpTime:  time.Since(a.Sr),
	})
	a.Cancel()
	a.Name = ""
	a.Ctx = nil
}

func (a *ActionManage) RunAction(ctx context.Context, run func() (loop time.Duration, next time.Duration)) time.Duration {
	tm := time.NewTimer(ms10)
	defer tm.Stop()
	for {
		select {
		case <-tm.C:
			loop, next := run()
			if loop == 0 {
				return next
			}
			tm.Reset(loop)
		case <-ctx.Done():
			ActionRunningHistoryList = append(ActionRunningHistoryList, ActionRunHistory{
				Name:        "主线程结束，等待重新启动",
				RunningTime: time.Now().Format("2006-01-02 15:04:05"),
				TakeUpTime:  time.Since(a.Sr),
			})
			return RandMillisecond(60, 120)
		case <-a.Ctx.Done():
			ActionRunningHistoryList = append(ActionRunningHistoryList, ActionRunHistory{
				Name:        a.Name + " 任务中止，等待下次执行",
				RunningTime: time.Now().Format("2006-01-02 15:04:05"),
				TakeUpTime:  time.Since(a.Sr),
			})
			return RandMillisecond(60, 120)
		}
	}
}

func SetAction(ctx context.Context, name string) *ActionManage {
	am := &ActionManage{Name: name, Sr: time.Now()}
	am.Ctx, am.Cancel = context.WithCancel(ctx)
	for i := range ActionManageList {
		if ActionManageList[i].Name == "" && ActionManageList[i].Ctx == nil {
			ActionManageList[i] = am
			return am
		}
	}
	ActionManageList = append(ActionManageList, am)
	return am
}

func StopAction() {
	for i := range ActionManageList {
		ActionManageList[i].End()
	}
}

var ActionManageList = make([]*ActionManage, 0)
var ActionRunningHistoryList = make([]ActionRunHistory, 0)
