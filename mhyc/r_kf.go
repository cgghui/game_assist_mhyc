package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// KuaFu 跨服
func KuaFu(ctx context.Context) {
	// 幻境
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	f1 := func() time.Duration {
		Fight.Lock()
		defer Fight.Unlock()
		// 章
		IllusionType := int32(2) // 11 2 1
		go func() {
			_ = CLI.IllusionData(IllusionType)
		}()
		data := &S2CIllusionData{}
		if err := Receive.Wait(data, s3); err != nil {
			return ms100
		}
		ChapterData := C2SIllusionChapterData{IllusionType: IllusionType}
		for _, c := range data.Chapters {
			if c.ChapterState == 2 {
				ChapterData.ChapterId = c.ChapterId
				break
			}
		}
		if ChapterData.ChapterId == 0 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		// 节
		go func() {
			_ = CLI.IllusionChapterData(&ChapterData)
		}()
		chapter := &S2CIllusionChapterData{}
		if err := Receive.Wait(chapter, s3); err != nil {
			return ms100
		}
		var teamParam C2SIllusionCreateTeam
		for _, c := range chapter.CheckPoints {
			if c.CheckPointState == 2 {
				teamParam = C2SIllusionCreateTeam{IllusionType: c.IllusionType, ChapterId: c.ChapterId, CheckPointId: c.CheckPointId}
				break
			}
		}
		// 组队
		go func() {
			_ = CLI.IllusionCreateTeam(&teamParam)
		}()
		team := &S2CIllusionCreateTeam{}
		if err := Receive.Wait(team, s3); err != nil {
			return ms100
		}
		//
		SelfFV := RoleInfo.Get("FightValue").Int64()
		for _, ir := range []int32{5, 1, 2, 3} {
			go func() {
				_ = CLI.GetIllusionCanInviteList(ir, &teamParam)
			}()
			userList := &S2CGetIllusionCanInviteList{}
			if err := Receive.Wait(userList, s3); err != nil {
				return ms100
			}
			for _, user := range userList.Users {
				if user.Fv > SelfFV {
					go func(ir int32, uid int64) {
						_ = CLI.IllusionTeamInviteUser(ir, uid)
					}(ir, user.UserId)
				}
			}
		}

		return time.Second
	}

	for {
		select {
		case <-t1.C:
			t1.Reset(f1())
		case <-ctx.Done():
			return
		}
	}
}

////////////////////////////////////////////////////////////

// IllusionChapterData 章节信息
func (c *Connect) IllusionChapterData(t *C2SIllusionChapterData) error {
	body, err := proto.Marshal(t)
	if err != nil {
		return err
	}
	log.Printf("[C][IllusionChapterData] illusion_type=%v chapter_id=%v", t.IllusionType, t.ChapterId)
	return c.send(25703, body)
}

func (x *S2CIllusionChapterData) ID() uint16 {
	return 25704
}

// Message S2CIllusionChapterData Code:25704
func (x *S2CIllusionChapterData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][IllusionChapterData] %v", x)
}

////////////////////////////////////////////////////////////

// IllusionCreateTeam 创建组队
func (c *Connect) IllusionCreateTeam(t *C2SIllusionCreateTeam) error {
	body, err := proto.Marshal(t)
	if err != nil {
		return err
	}
	log.Printf("[C][IllusionCreateTeam] illusion_type=%v chapter_id=%v check_point_id=%v", t.IllusionType, t.ChapterId, t.CheckPointId)
	return c.send(25709, body)
}

func (x *S2CIllusionCreateTeam) ID() uint16 {
	return 25710
}

// Message S2CIllusionCreateTeam Code:25710
func (x *S2CIllusionCreateTeam) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][IllusionCreateTeam] %v", x)
}

////////////////////////////////////////////////////////////

// IllusionTeamInviteAll 邀人所有人加入
// ir 5 战队
// ir 1 仙宗
// ir 2 仙缘
// ir 3 家族
func (c *Connect) IllusionTeamInviteAll(ir int32) error {
	body, err := proto.Marshal(&C2SIllusionTeamInvite{InviteUserId: 0, InviteRange: ir})
	if err != nil {
		return err
	}
	log.Println("[C][IllusionTeamInvite] all")
	return c.send(25713, body)
}

// IllusionTeamInviteUser 邀人指定UserID加入
// ir 5 战队
// ir 1 仙宗
// ir 2 仙缘
// ir 3 家族
func (c *Connect) IllusionTeamInviteUser(ir int32, userID int64) error {
	body, err := proto.Marshal(&C2SIllusionTeamInvite{InviteUserId: userID, InviteRange: ir})
	if err != nil {
		return err
	}
	log.Println("[C][IllusionTeamInvite] all")
	return c.send(25713, body)
}

func (x *S2CIllusionTeamInvite) ID() uint16 {
	return 25714
}

// Message S2CIllusionTeamInvite Code:25714
func (x *S2CIllusionTeamInvite) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][IllusionTeamInvite] %v", x)
}

////////////////////////////////////////////////////////////

// GetIllusionCanInviteList 人员列表
// ir 5 战队
// ir 1 仙宗
// ir 2 仙缘
// ir 3 家族
func (c *Connect) GetIllusionCanInviteList(ir int32, ct *C2SIllusionCreateTeam) error {
	body, err := proto.Marshal(&C2SGetIllusionCanInviteList{IllusionType: ct.IllusionType, ChapterId: ct.ChapterId, CheckPointId: ct.CheckPointId, InviteRange: ir})
	if err != nil {
		return err
	}
	log.Println("[C][GetIllusionCanInviteList] all")
	return c.send(25725, body)
}

func (x *S2CGetIllusionCanInviteList) ID() uint16 {
	return 25726
}

// Message S2CGetIllusionCanInviteList Code:25714
func (x *S2CGetIllusionCanInviteList) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetIllusionCanInviteList] %v", x)
}

////////////////////////////////////////////////////////////

// IllusionData 幻境数据
func (c *Connect) IllusionData(it int32) error {
	body, err := proto.Marshal(&C2SIllusionData{IllusionType: it})
	if err != nil {
		return err
	}
	log.Printf("[C][IllusionData] illusion_type=%v", it)
	return c.send(25701, body)
}

func (x *S2CIllusionData) ID() uint16 {
	return 25702
}

// Message S2CIllusionData Code:25702
func (x *S2CIllusionData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][IllusionData] %v", x)
}
