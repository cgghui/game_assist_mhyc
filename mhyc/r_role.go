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

// Message S2CBattlefieldReport Code:101
func (x *S2CBattlefieldReport) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BattlefieldReport] win=%v idx=%v", x.Win, x.Idx)
}

var RoleInfo = &roleInfo{
	s: &sync.Map{},
}

type roleInfo struct{ s *sync.Map }

type roleValue struct {
	val interface{}
}

func (r roleValue) Int64() int64 {
	return r.val.(int64)
}

func (r roleValue) String() string {
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
			if err = Receive.WaitWithContext(ctx, 3, &S2CRoleInfo{}); err != nil {
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
