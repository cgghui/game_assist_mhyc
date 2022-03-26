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

// WareHouseReceiveItem 将仓库内的物品转至背包
// 5 寻找
func (c *Connect) WareHouseReceiveItem(id int32) error {
	body, err := proto.Marshal(&C2SWareHouseReceiveItem{WhId: id})
	if err != nil {
		return err
	}
	return c.send(27305, body)
}

// EndFight 结束战斗
func (c *Connect) EndFight(r *S2CBattlefieldReport) error {
	body, err := proto.Marshal(&C2SEndFight{Idx: r.Idx})
	if err != nil {
		return err
	}
	return c.send(102, body)
}

// RoleInfo 角色信息
func (c *Connect) RoleInfo() error {
	body, err := proto.Marshal(&C2SRoleInfo{})
	if err != nil {
		return err
	}
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
		log.Printf("[S][RoleInfo] %v", err)
		return
	}
	for _, a := range x.A {
		if _, ok := AttrType[a.K]; !ok {
			continue
		}
		log.Printf("[S][RoleInfo] %s\t%v", AttrType[a.K], a.V)
		_, _ = RIF.WriteString(fmt.Sprintf("%s\t%v\n", AttrType[a.K], a.V))
		RoleInfo.Set(AttrType[a.K], a.V)
	}
	for _, b := range x.B {
		if _, ok := AttrType[b.K]; !ok {
			continue
		}
		log.Printf("[S][RoleInfo] %s\t%v", AttrType[b.K], b.V)
		_, _ = RIF.WriteString(fmt.Sprintf("%s\t%v\n", AttrType[b.K], b.V))
		RoleInfo.Set(AttrType[b.K], b.V)
	}
	return
}

////////////////////////////////////////////////////////////

func (x *S2CBattlefieldReport) ID() uint16 {
	return 101
}

// Message S2CBattlefieldReport Code:101
func (x *S2CBattlefieldReport) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
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

func (x *Pong) ID() uint16 {
	return 23
}

// Message Pong Code:23
func (x *Pong) Message(_ []byte) {
	log.Printf("[S][Pong]")
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
