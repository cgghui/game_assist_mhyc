package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// XianDianSSSL 仙宗 - 仙殿 - 神兽试炼
func XianDianSSSL() {
	t := time.NewTimer(ms100)
	f := func() time.Duration {
		if RoleInfo.Get("SectGodAnimalChallenge").Int64() >= 35 {
			curr := time.Now()
			next := SelfWeekMonday(curr).Add(169 * time.Hour) // 168小时 = 7天
			return next.Sub(curr)
		}
		k := time.NewTimer(ms100)
		defer k.Stop()
		for range k.C {
			info := &S2CChallengeGodAnimal{}
			Receive.Action(CLI.ChallengeGodAnimal)
			_ = Receive.Wait(info, s3)
			if info.Tag == 48025 {
				break
			}
			k.Reset(ms500)
		}
		return RandMillisecond(86400, 115200)
	}
	for range t.C {
		t.Reset(f())
	}
}

// SectGodAnimalData 仙宗 - 仙殿 - 神兽试炼 - 信息
func (c *Connect) SectGodAnimalData() error {
	body, err := proto.Marshal(&C2SChallengeGodAnimal{})
	if err != nil {
		return err
	}
	return c.send(19043, body)
}

func (c *Connect) ChallengeGodAnimal() error {
	body, err := proto.Marshal(&C2SChallengeGodAnimal{})
	if err != nil {
		return err
	}
	return c.send(19045, body)
}

////////////////////////////////////////////////////////////

func (x *S2CSectGodAnimalData) ID() uint16 {
	return 19044
}

// Message S2CSectGodAnimalData 19044
func (x *S2CSectGodAnimalData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SectGodAnimalData] %v", x)
}

////////////////////////////////////////////////////////////

func (x *S2CChallengeGodAnimal) ID() uint16 {
	return 19046
}

// Message S2CChallengeGodAnimal 19046
func (x *S2CChallengeGodAnimal) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ChallengeGodAnimal] tag=%v is_win=%v", x.Tag, x.IsWin)
}
