package mhyc

import (
	"context"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// XianDianXDSW 仙宗 - 仙殿 - 仙宗声望 // 晋升
func XianDianXDSW(ctx context.Context) {
	t := time.NewTimer(ms100)
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "仙宗-仙殿-仙宗声望")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		if s := RoleInfo.SectPrestige(); s != nil && int(RoleInfo.Get("SectPrestigeVal").Int64()) >= s.Prestige {
			Receive.Action(CLI.SectPrestigeLevelUp)
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CSectPrestigeLevelUp{}, s3)
		}
		return RandMillisecond(60, 180) // 1 ~ 3 分钟
	}
	for {
		select {
		case <-t.C:
			t.Reset(f())
		case <-ctx.Done():
			return
		}
	}
}

// SectPrestigeRecv 仙宗 - 仙殿 - 仙宗声望 - 每日奉碌
func (c *Connect) SectPrestigeRecv() error {
	body, err := proto.Marshal(&C2SSectPrestigeRecv{})
	if err != nil {
		return err
	}
	log.Println("[C][SectPrestigeRecv]")
	return c.send(19055, body)
}

// SectPrestigeLevelUp 仙宗 - 仙殿 - 仙宗声望 - 晋升
func (c *Connect) SectPrestigeLevelUp() error {
	body, err := proto.Marshal(&C2SSectPrestigeLevelUp{})
	if err != nil {
		return err
	}
	log.Println("[C][SectPrestigeLevelUp]")
	return c.send(19053, body)
}

////////////////////////////////////////////////////////////

func (x *S2CSectPrestigeRecv) ID() uint16 {
	return 19056
}

// Message S2CSectPrestigeRecv 19056
func (x *S2CSectPrestigeRecv) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SectPrestigeRecv] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CSectPrestigeLevelUp) ID() uint16 {
	return 19054
}

// Message S2CSectPrestigeLevelUp 19054
func (x *S2CSectPrestigeLevelUp) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SectPrestigeLevelUp] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

type SectPrestige struct {
	AttrId   int    `json:"AttrId"`
	Items    string `json:"Items"`
	Level    int    `json:"Level"`
	Prestige int    `json:"Prestige"`
	Title    int    `json:"Title"`
}

var sectPrestigeLS []SectPrestige

func init() {
	if err := json.Unmarshal(cfg2SectPrestige, &sectPrestigeLS); err != nil {
		panic(err)
	}
}

func (r *roleInfo) SectPrestige() *SectPrestige {
	lv := int(r.Get("SectPrestigeLevel").Int64())
	for i, s := range sectPrestigeLS {
		if s.Level == lv {
			return &sectPrestigeLS[i]
		}
	}
	return nil
}
