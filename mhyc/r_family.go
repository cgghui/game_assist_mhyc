package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func init() {
	PCK[27358] = &S2CFamilyJJCJoin{}
	PCK[27363] = &S2CFamilyJJCFight{}
}

var familyJJCThread = make(chan interface{})
var familyJJCAction = make(chan struct{})

// FamilyJJC 竞技
// [RoleInfo] FamilyJJC_TimesLeft 剩余次数
// [RoleInfo] FamilyJJC_Times	  使用次数
// [RoleInfo] FamilyJJC_Score	  积分
func (c *Connect) FamilyJJC() {
	go func() {
		familyJJCAction <- struct{}{}
		t := time.NewTimer(time.Minute)
		for range t.C {
			familyJJCAction <- struct{}{}
			t.Reset(10 * time.Minute)
		}
	}()
	run := func(val interface{}) {
		switch ret := val.(type) {
		case *S2CFamilyJJCJoin:
			if ret.Tag == 0 {
				_ = c.familyJJCFight(ret)
				return
			}
			if ret.Tag == 17003 {
				time.AfterFunc(time.Second, func() {
					familyJJCAction <- struct{}{}
				})
				return
			}
			// end
			if ret.Tag == 57606 {

			}
		case *S2CFamilyJJCFight:
			familyJJCAction <- struct{}{}
		}
	}
	for {
		select {
		case <-familyJJCAction:
			if val, ok := RoleInfo.Load("FamilyJJC_Times"); ok {
				if n, y := val.(int64); y && n >= 10 {
					break
				}
			}
			_ = c.familyJJCJoin()
		case val := <-familyJJCThread:
			go run(val)
		}
	}
}

func (c *Connect) familyJJCJoin() error {
	body, err := proto.Marshal(&C2SFamilyJJCJoin{})
	if err != nil {
		return err
	}
	return c.send(27357, body)
}

func (c *Connect) familyJJCFight(act *S2CFamilyJJCJoin) error {
	dat := C2SFamilyJJCFight{
		UserId: make([]int64, 0, 0),
	}
	for _, self := range act.Self {
		dat.UserId = append(dat.UserId, self.UserId)
	}
	body, err := proto.Marshal(&dat)
	if err != nil {
		return err
	}
	return c.send(27359, body)
}

func (x *S2CFamilyInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][FamilyInfo] tag=%v", x.Tag)
	familyJJCThread <- x
}

func (x *S2CFamilyJJCJoin) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][FamilyJJCJoin] tag=%v", x.Tag)
	familyJJCThread <- x
}

func (x *S2CFamilyJJCFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][FamilyJJCFight] tag=%v", x.Tag)
	familyJJCThread <- x
}
