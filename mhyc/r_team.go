package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
)

// CreateTeam 创建组队
func (c *Connect) CreateTeam(t *C2SCreateTeam) error {
	body, err := proto.Marshal(t)
	if err != nil {
		return err
	}
	log.Printf("[C][CreateTeam] func_id=%v", t.FuncId)
	return c.send(27105, body)
}

func (x *S2CCreateTeam) ID() uint16 {
	return 27106
}

// Message S2CCreateTeam Code:27106
func (x *S2CCreateTeam) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][CreateTeam] %v", x)
}

////////////////////////////////////////////////////////////

// Teams 组队
func (c *Connect) Teams(t *C2STeams) error {
	body, err := proto.Marshal(t)
	if err != nil {
		return err
	}
	log.Printf("[C][Teams] func_id=%v", t.FuncId)
	return c.send(27101, body)
}

func (x *S2CTeams) ID() uint16 {
	return 27102
}

// Message S2CTeams Code:27102
func (x *S2CTeams) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][Teams] %v", x)
}
