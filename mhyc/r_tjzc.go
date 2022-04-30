package mhyc

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// TJZC 天降珍宠
func TJZC(ctx context.Context) {
	t1 := time.NewTimer(ms10)
	f1 := func() time.Duration {
		Fight.Lock()
		cx, cancel := context.WithCancel(ctx)
		go ListenMessageCall(cx, &S2CTreatData{}, func(data []byte) {
			r := &S2CTreatData{}
			r.Message(data)
			fmt.Println(r)
		})
		go ListenMessageCall(cx, &S2CSectCrossBossBox{}, func(data []byte) {
			r := &S2CSectCrossBossBox{}
			r.Message(data)
			fmt.Println(r)
		})
		// 列表宠物
		go ListenMessageCall(cx, &S2CCrossPastureTrapList{}, func(data []byte) {
			r := &S2CCrossPastureTrapList{}
			r.Message(data)
			if len(r.TrapList) == 0 {
				return
			}
			for i := range r.TrapList {
				go func(pet *CrossPastureTrap) {
					ts := ms10
					cur := time.Now()
					tmv := time.Unix(pet.Wait, 0).Local().Add(ms500)
					if cur.Before(tmv) {
						ts = tmv.Sub(cur)
					}
					select {
					case <-ctx.Done():
						return
					case <-time.After(ts):
						go func() {
							_ = CLI.StartMove(&C2SStartMove{P: []int32{pet.X, pet.Y}})
						}()
						_ = Receive.Wait(&S2CStartMove{}, s3)
						go func() {
							_ = CLI.CatchCrossPet(pet.PetId, pet.X, pet.Y)
						}()
						_ = Receive.Wait(&S2CCatchCrossPet{}, s3)
						return
					}
				}(r.TrapList[i])
			}
		})
		go ListenMessageCall(cx, &S2CFightBossTimes{}, func(data []byte) {
			r := &S2CFightBossTimes{}
			r.Message(data)
			fmt.Println(r)
		})
		go ListenMessageCall(cx, &S2CSectCrossSeizePets{}, func(data []byte) {
			r := &S2CSectCrossSeizePets{}
			r.Message(data)
			fmt.Println(r)
		})
		// 进入场景
		Receive.Action(CLI.EnterCrossPasture)
		err := Receive.Wait(&S2CEnterCrossPasture{}, s3)
		defer func() {
			if err == nil {
				Receive.Action(CLI.LeaveCrossPasture)
				_ = Receive.Wait(&S2CLeaveCrossPasture{}, s3)
			}
			cancel()
			Fight.Unlock()
		}()
		// 使用抓捕道具
		tm := time.NewTimer(ms10)
		defer tm.Stop()
		items := []int32{504, 505, 506} // 捕获道具
		i := 0
		for {
			select {
			case <-tm.C:
				if i >= 3 {
					goto Next
				}
				go func(i int) {
					_ = CLI.UseTrap(items[i])
				}(i)
				r := &S2CUseTrap{}
				if _ = Receive.Wait(r); r.Tag != 0 {
					i++
				}
				tm.Reset(ms100)
			case <-ctx.Done():
				return s3
			}
		}
	Next:
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

////////////////////////////////////////////////////////////

// EnterCrossPasture 进入
func (c *Connect) EnterCrossPasture() error {
	body, err := proto.Marshal(&C2SEnterCrossPasture{})
	if err != nil {
		return err
	}
	log.Println("[C][EnterCrossPasture]")
	return c.send(19201, body)
}

func (x *S2CEnterCrossPasture) ID() uint16 {
	return 19202
}

// Message S2CEnterCrossPasture Code:19202
func (x *S2CEnterCrossPasture) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][EnterCrossPasture] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

// LeaveCrossPasture 退出
func (c *Connect) LeaveCrossPasture() error {
	body, err := proto.Marshal(&C2SLeaveCrossPasture{})
	if err != nil {
		return err
	}
	log.Println("[C][LeaveCrossPasture]")
	return c.send(19203, body)
}

func (x *S2CLeaveCrossPasture) ID() uint16 {
	return 19204
}

// Message S2CLeaveCrossPasture Code:19204
func (x *S2CLeaveCrossPasture) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][EnterCrossPasture] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CCrossPastureTrapList) ID() uint16 {
	return 19218
}

// Message S2CCrossPastureTrapList Code:19218
func (x *S2CCrossPastureTrapList) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][CrossPastureTrapList] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

// SectCrossBossBox 退出
func (c *Connect) SectCrossBossBox() error {
	body, err := proto.Marshal(&C2SSectCrossBossBox{})
	if err != nil {
		return err
	}
	log.Println("[C][SectCrossBossBox]")
	return c.send(19213, body)
}

func (x *S2CSectCrossBossBox) ID() uint16 {
	return 19214
}

// Message S2CSectCrossBossBox Code:19214
func (x *S2CSectCrossBossBox) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SectCrossBossBox] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CTreatData) ID() uint16 {
	return 1520
}

// Message S2CTreatData Code:1520
func (x *S2CTreatData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][TreatData] items=%v", x.Items)
}

////////////////////////////////////////////////////////////

// FightBossTimes 退出
func (c *Connect) FightBossTimes(id int32) error {
	body, err := proto.Marshal(&C2SFightBossTimes{ActId: id})
	if err != nil {
		return err
	}
	log.Println("[C][FightBossTimes]")
	return c.send(19227, body)
}

func (x *S2CFightBossTimes) ID() uint16 {
	return 19228
}

// Message S2CFightBossTimes Code:19228
func (x *S2CFightBossTimes) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][FightBossTimes] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CSectCrossSeizePets) ID() uint16 {
	return 19230
}

// Message S2CSectCrossSeizePets Code:19230
func (x *S2CSectCrossSeizePets) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SectCrossSeizePets] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

// UseTrap 使用道具
// 505 稀有补兽夹
func (c *Connect) UseTrap(id int32) error {
	body, err := proto.Marshal(&C2SUseTrap{ItemId: id})
	if err != nil {
		return err
	}
	log.Printf("[C][C2SUseTrap] item_id=%d", id)
	return c.send(19205, body)
}

func (x *S2CUseTrap) ID() uint16 {
	return 19206
}

// Message S2CUseTrap Code:19206
func (x *S2CUseTrap) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][UseTrap] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

// CatchCrossPet 抓捕
func (c *Connect) CatchCrossPet(id, x, y int32) error {
	body, err := proto.Marshal(&C2SCatchCrossPet{PetId: id, X: x, Y: y})
	if err != nil {
		return err
	}
	log.Printf("[C][CatchCrossPet] item_id=%d", id)
	return c.send(19207, body)
}

func (x *S2CCatchCrossPet) ID() uint16 {
	return 19208
}

// Message S2CCatchCrossPet Code:19208
func (x *S2CCatchCrossPet) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][CatchCrossPet] tag=%v pet_id=%d", x.Tag, x.PetId)
}
