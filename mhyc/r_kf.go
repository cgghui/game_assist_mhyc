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
	// 11 2 1
	f1 := func(IllusionType int32) time.Duration {
		Fight.Lock()
		defer func() {
			Receive.Action(CLI.IllusionDelTeam)
			_ = Receive.Wait(&S2CIllusionDelTeam{}, s3)
			Fight.Unlock()
		}()
		// 章
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
			return RandMillisecond(1800, 3600)
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
		cx, cancel := context.WithTimeout(ctx, s20) // 最大搜寻队友时间
		defer cancel()
		// 组成人员
		ListenMessageCall(cx, &S2CIllusionMyTeam{}, func(data []byte) {
			info := &S2CIllusionMyTeam{}
			info.Message(data)
			if len(info.Team.Users) >= 3 {
				return
			}
		})
		// 邀人
		go func() {
			tc := time.NewTimer(ms100)
			defer tc.Stop()
			for {
				select {
				case <-tc.C:
					MyFV := RoleInfo.Get("FightValue").Int64()
					for _, ir := range []int32{1, 2, 3} {
						go func() {
							_ = CLI.GetIllusionCanInviteList(ir, &teamParam)
						}()
						userList := &S2CGetIllusionCanInviteList{}
						if err := Receive.Wait(userList, s3); err != nil {
							tc.Reset(ms500)
							break
						}
						for _, user := range userList.Users {
							if user.Fv > MyFV {
								go func(ir int32, uid int64) {
									_ = CLI.IllusionTeamInviteUser(ir, uid)
								}(ir, user.UserId)
							}
						}
						tc.Reset(time.Second * 62)
					}
				case <-cx.Done():
					ir := int32(5)
					go func() {
						_ = CLI.GetIllusionCanInviteList(ir, &teamParam)
					}()
					userList := &S2CGetIllusionCanInviteList{}
					if err := Receive.Wait(userList, s3); err != nil {
						tc.Reset(ms500)
						break
					}
					for _, user := range userList.Users {
						go func(ir int32, uid int64) {
							_ = CLI.IllusionTeamInviteUser(ir, uid)
						}(ir, user.UserId)

					}
					return
				}
			}
		}()
		// 监听加入
		ListenMessageCall(cx, &S2CIllusionTeamInvite{}, func(data []byte) {
			ti := &S2CIllusionTeamInvite{}
			ti.Message(data)
			if ti.Tag == 56119 { // 拒绝
				return
			}
		})
		// 组队开战
		func() {
			n := 0
			t := time.NewTimer(time.Second)
			defer t.Stop()
			for range t.C {
				report := &S2CBattlefieldReport{}
				rc := ListenMessageNotify(report, s3)
				Receive.Action(CLI.IllusionFight)
				fr := &S2CIllusionFight{}
				_ = Receive.Wait(fr, s3)
				<-rc
				if fr.Tag == 0 && report.Win == 1 {
					if fr.ChapterId >= 9 && fr.CheckPointId >= 9 {
						return
					}
					n = 0
					next := C2SIllusionTeamNextCheckPoint{IllusionType: IllusionType}
					if fr.CheckPointId == 9 {
						next.ChapterId = fr.ChapterId + 1
						next.CheckPointId = 1
					} else {
						next.ChapterId = fr.ChapterId
						next.CheckPointId = fr.CheckPointId + 1
					}
					go func() {
						_ = CLI.IllusionTeamNextCheckPoint(&next)
					}()
					_ = Receive.Wait(&S2CIllusionTeamNextCheckPoint{}, s3)
				}
				if fr.Tag == 0 && report.Win == 0 && fr.RetStr == "" {
					n++
				}
				if n >= 3 {
					return
				}
				t.Reset(time.Second)
			}
		}()
		return RandMillisecond(300, 600)
	}

	for {
		select {
		case <-t1.C:
			t1.Reset(f1(2))
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

// IllusionDelTeam 离开组队
func (c *Connect) IllusionDelTeam() error {
	body, err := proto.Marshal(&C2SIllusionDelTeam{})
	if err != nil {
		return err
	}
	log.Printf("[C][IllusionDelTeam]")
	return c.send(25719, body)
}

func (x *S2CIllusionDelTeam) ID() uint16 {
	return 25720
}

// Message S2CIllusionDelTeam Code:25720
func (x *S2CIllusionDelTeam) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][IllusionDelTeam] %v", x)
}

////////////////////////////////////////////////////////////

// IllusionFight 组队开战
func (c *Connect) IllusionFight() error {
	body, err := proto.Marshal(&C2SIllusionFight{})
	if err != nil {
		return err
	}
	log.Printf("[C][IllusionFight]")
	return c.send(25721, body)
}

func (x *S2CIllusionFight) ID() uint16 {
	return 25722
}

// Message S2CIllusionFight Code:25722
func (x *S2CIllusionFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][IllusionFight] %v", x)
}

////////////////////////////////////////////////////////////

// IllusionTeamNextCheckPoint 组队下一章节
func (c *Connect) IllusionTeamNextCheckPoint(next *C2SIllusionTeamNextCheckPoint) error {
	body, err := proto.Marshal(next)
	if err != nil {
		return err
	}
	log.Printf("[C][IllusionTeamNextCheckPoint]")
	return c.send(25733, body)
}

func (x *S2CIllusionTeamNextCheckPoint) ID() uint16 {
	return 25734
}

// Message S2CIllusionTeamNextCheckPoint Code:25734
func (x *S2CIllusionTeamNextCheckPoint) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][IllusionTeamNextCheckPoint] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

// IllusionMyTeam 我的组队
func (c *Connect) IllusionMyTeam() error {
	body, err := proto.Marshal(&C2SIllusionMyTeam{})
	if err != nil {
		return err
	}
	log.Printf("[C][IllusionMyTeam]")
	return c.send(25707, body)
}

func (x *S2CIllusionMyTeam) ID() uint16 {
	return 25708
}

// Message S2CIllusionMyTeam Code:25708
func (x *S2CIllusionMyTeam) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	if x.Team == nil {
		return
	}
	log.Printf("[S][S2CIllusionMyTeam] team_id=%v user_len=%v", x.Team.TeamId, len(x.Team.Users))
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
