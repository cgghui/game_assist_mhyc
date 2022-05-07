package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// JJC 竞技场
func JJC(ctx context.Context) {
	t1 := time.NewTimer(ms10)
	f1 := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "竞技-竞技场")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		Receive.Action(CLI.JJCList)
		ls := &S2CJJCList{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, ls, s3); err != nil {
			return RandMillisecond(0, 2)
		}
		targetId := int32(0)
		targetRank := int32(0)
		fv := RoleInfo.Get("FightValue").Int64()
		for _, role := range ls.Roles {
			ok := false
			for _, a := range role.A {
				if a.K == 9999 {
					if fv > a.V {
						ok = true
						break
					}
				}
			}
			if !ok {
				continue
			}
			targetId = role.UserId
			for _, a := range role.A {
				if a.K == 140 {
					targetRank = int32(a.V)
					break
				}
			}
			break
		}
		if targetId == 0 && targetRank == 0 {
			Receive.Action(CLI.JJCSweep)
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CJJCSweep{}, s3)
			return RandMillisecond(60, 300)
		}
		go func() {
			_ = CLI.JJCFight(targetId, targetRank)
		}()
		r := &S2CJJCFight{}
		if _ = Receive.WaitWithContextOrTimeout(am.Ctx, r, s3); r.Tag == 11002 {
			return RandMillisecond(60, 300)
		}
		return ms500
	}
	defer t1.Stop()
	for {
		select {
		case <-t1.C:
			t1.Reset(f1())
		case <-ctx.Done():
			return
		}
	}
}

func WZZB(ctx context.Context) {
	t1 := time.NewTimer(ms10)
	f1 := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "竞技-王者争霸")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		Receive.Action(CLI.KingMatch)
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CKingMatch{}, s3); err != nil {
			return RandMillisecond(0, 2)
		}
		Receive.Action(CLI.KingFight)
		r := &S2CKingFight{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, r, s3); err != nil {
			return RandMillisecond(0, 2)
		}
		if r.Tag == 0 {
			return ms500
		}
		return RandMillisecond(60, 300)
	}
	defer t1.Stop()
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

// JJCList 列表
func (c *Connect) JJCList() error {
	body, err := proto.Marshal(&C2SJJCList{})
	if err != nil {
		return err
	}
	log.Println("[C][JJCList]")
	return c.send(1101, body)
}

func (x *S2CJJCList) ID() uint16 {
	return 1102
}

// Message S2CJJCList Code:1102
func (x *S2CJJCList) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JJCList] tag=%v role=%v", x.Tag, x.Roles)
}

////////////////////////////////////////////////////////////

// JJCFight 战
func (c *Connect) JJCFight(targetId, targetRank int32) error {
	body, err := proto.Marshal(&C2SJJCFight{TargetId: targetId, TargetRank: targetRank})
	if err != nil {
		return err
	}
	log.Printf("[C][JJCFight] target_id=%v target_rank=%v", targetId, targetRank)
	return c.send(1103, body)
}

func (x *S2CJJCFight) ID() uint16 {
	return 1104
}

// Message S2CJJCFight Code:1104
func (x *S2CJJCFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JJCFight] tag=%v %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

// JJCSweep 扫
func (c *Connect) JJCSweep() error {
	body, err := proto.Marshal(&C2SJJCSweep{})
	if err != nil {
		return err
	}
	log.Println("[C][JJCSweep]")
	return c.send(1109, body)
}

func (x *S2CJJCSweep) ID() uint16 {
	return 1110
}

// Message S2CJJCSweep Code:1110
func (x *S2CJJCSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JJCSweep] tag=%v %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

// KingMatch 匹配
func (c *Connect) KingMatch() error {
	body, err := proto.Marshal(&C2SKingMatch{})
	if err != nil {
		return err
	}
	log.Println("[C][KingMatch]")
	return c.send(14002, body)
}

func (x *S2CKingMatch) ID() uint16 {
	return 14003
}

// Message S2CKingMatch Code:14003
func (x *S2CKingMatch) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][KingMatch] tag=%v %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

// KingFight 战
func (c *Connect) KingFight() error {
	body, err := proto.Marshal(&C2SKingFight{})
	if err != nil {
		return err
	}
	log.Println("[C][KingFight]")
	return c.send(14004, body)
}

func (x *S2CKingFight) ID() uint16 {
	return 14005
}

// Message S2CKingFight Code:14005
func (x *S2CKingFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][KingFight] tag=%v", x.Tag)
}
