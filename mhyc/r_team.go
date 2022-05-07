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

////////////////////////////////////////////////////////////

func (c *Connect) InviteTeam(teamID int64, inviteType int32) error {
	body, err := proto.Marshal(&C2SInviteTeam{TeamId: teamID, UserId: 0, InviteType: inviteType})
	if err != nil {
		return err
	}
	log.Printf("[C][InviteTeam] team_id=%v invite_type=%v", teamID, inviteType)
	return c.send(27109, body)
}

func (x *S2CInviteTeam) ID() uint16 {
	return 27110
}

// Message S2CInviteTeam 27110
func (x *S2CInviteTeam) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][InviteTeam] tag=%v tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (c *Connect) TeamInfo() error {
	body, err := proto.Marshal(&C2STeamInfo{})
	if err != nil {
		return err
	}
	log.Println("[C][组队信息]")
	return c.send(27103, body)
}

func (x *S2CTeamInfo) ID() uint16 {
	return 27104
}

// Message S2CTeamInfo 27104
func (x *S2CTeamInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][组队信息] tag=%v tag_msg=%s team=%v players=%v", x.Tag, GetTagMsg(x.Tag), x.Team, x.Players)
}

////////////////////////////////////////////////////////////

func (c *Connect) LeaveTeam(teamID int64) error {
	body, err := proto.Marshal(&C2SLeaveTeam{TeamId: teamID})
	if err != nil {
		return err
	}
	log.Println("[C][退出组队]")
	return c.send(27111, body)
}

func (x *S2CLeaveTeam) ID() uint16 {
	return 27112
}

// Message S2CLeaveTeam 27112
func (x *S2CLeaveTeam) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][退出组队] tag=%v tag_msg=%s func_id=%v", x.Tag, GetTagMsg(x.Tag), x.FuncId)
}
