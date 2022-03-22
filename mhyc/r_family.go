package mhyc

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func init() {
	PCK[20002] = &S2CFamilyInfo{}
	PCK[27358] = &S2CFamilyJJCJoin{}
	PCK[27363] = &S2CFamilyJJCFight{}
}

var jjcThread = make(chan interface{})
var jjcAction = make(chan struct{})

// JJC 竞技
func (c *Connect) JJC() {
	go func() {
		jjcAction <- struct{}{}
		t := time.NewTimer(time.Minute)
		for range t.C {
			jjcAction <- struct{}{}
			t.Reset(10 * time.Minute)
		}
	}()
	for range jjcAction {
		_ = c.familyJJCJoin()
		join := (<-jjcThread).(*S2CFamilyJJCJoin)
		fmt.Println(join)
		_ = c.familyJJCFight(join)
		fight := (<-jjcThread).(*S2CFamilyJJCFight)
		fmt.Println(fight)
	}
}

func (c *Connect) familyInfo() error {
	body, err := proto.Marshal(&C2SFamilyInfo{FuncType: 0})
	if err != nil {
		return err
	}
	return c.send(20001, body)
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
	jjcThread <- x
}

func (x *S2CFamilyJJCJoin) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][FamilyJJCJoin] tag=%v", x.Tag)
	jjcThread <- x
}

func (x *S2CFamilyJJCFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][FamilyJJCFight] tag=%v", x.Tag)
	jjcThread <- x
}
