package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
)

func (c *Connect) JoinActive(act *C2SJoinActive) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][JoinActive] aid=%v", act.AId)
	return c.send(1507, body)
}

func (c *Connect) LeaveActive(act *C2SLeaveActive) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][LeaveActive] aid=%v", act.AId)
	return c.send(1509, body)
}

////////////////////////////////////////////////////////////

func (x *S2CJoinActive) ID() uint16 {
	return 1508
}

// Message S2CJoinActive 1508
func (x *S2CJoinActive) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JoinActive] tag=%v tag_msg=%s aid=%v", x.Tag, GetTagMsg(x.Tag), x.AId)
}

////////////////////////////////////////////////////////////

func (x *S2CLeaveActive) ID() uint16 {
	return 1510
}

// Message S2CLeaveActive 1510
func (x *S2CLeaveActive) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][LeaveActive] tag=%v tag_msg=%s aid=%v", x.Tag, GetTagMsg(x.Tag), x.AId)
}
