package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"time"
)

const (
	ItemPet500 int32 = 500 // 缚妖索
	ItemPet501 int32 = 501 // 高级缚妖索
	ItemPet502 int32 = 502 // 神兽号角
	ItemPet503 int32 = 503 // 金铲子
)

func EnterAnimalPark(ctx context.Context) {
	t := time.NewTimer(ms10)
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "抓捕宠物")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		// 进入
		Receive.Action(CLI.EnterAnimalPark)
		ret := &S2CEnterAnimalPark{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, ret, s10); err != nil {
			return RandMillisecond(6, 12)
		}
		defer func() {
			Receive.Action(CLI.LeaveAnimalPark)
			_ = Receive.Wait(&S2CLeaveAnimalPark{}, s3)
		}()
		if ret.Pet != nil {
			count := len(ret.Pet)
			i := 0
			r := am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
				pet := ret.Pet[i]
				go func(r *PasturePet) {
					_ = CLI.AnimalParkGO(&C2SAnimalParkGO{
						PetId: r.Id,
						X:     r.PointX,
						Y:     r.PointY,
					})
				}(pet)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CAnimalParkGO{}, s3); err != nil {
					return 0, 0
				}
				i++
				if i >= count {
					loop = 0
					next = time.Second
					return
				}
				return ms100, 0
			})
			if r == 0 {
				return RandMillisecond(6, 12)
			}
		}
		var items = make(map[int32]*ItemData)
		var n = int64(0)
		if item := UserBag.Get(ItemPet500); item != nil {
			items[ItemPet500] = item
			n += items[ItemPet500].N
		}
		if item := UserBag.Get(ItemPet501); item != nil {
			items[ItemPet501] = item
			n += items[ItemPet501].N
		}
		if item := UserBag.Get(ItemPet502); item != nil {
			items[ItemPet502] = item
			n += items[ItemPet502].N
		}
		if n < 1000 {
			date := time.Now().Format("2006-01-02")
			if b, err := ioutil.ReadFile("./AnimalParkSearch10.txt"); err == nil && date == string(b) {
				return RandMillisecond(1800, 3600)
			}
			if item := UserBag.Get(ItemPet500); item != nil && item.N >= 10 {
				c := 0
				am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
					if c >= 10 {
						_ = ioutil.WriteFile("./AnimalParkSearch10.txt", []byte(date), 0666)
						return 0, 0
					}
					// s
					go func() {
						_ = CLI.SearchPet(&C2SSearchPet{ItemId: ItemPet500})
					}()
					r := &S2CSearchPet{}
					if err := Receive.WaitWithContextOrTimeout(am.Ctx, r, s3); err == nil && r.Pet != nil {
						c++
						go func(r *S2CSearchPet) {
							_ = CLI.AnimalParkGO(&C2SAnimalParkGO{
								PetId: r.Pet.Id,
								X:     r.Pet.PointX,
								Y:     r.Pet.PointY,
							})
						}(r)
						_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CAnimalParkGO{}, s3)
					}
					return ms100, 0
				})
			}
			return RandMillisecond(1800, 3600)
		}
		// 检测是否需要使用buff
		isBuff := false
		for _, buff := range ret.Buff {
			if buff.BuffId == 2 {
				isBuff = true
				break
			}
		}
		// 使用【金铲子】buff
		if !isBuff {
			if item := UserBag.Wait(ItemPet503, s3, am.Ctx); item != nil {
				v := item
				if v.N > 0 {
					Receive.Action(func() error {
						return CLI.SearchPet(&C2SSearchPet{ItemId: ItemPet503})
					})
					if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CSearchPet{}, s3); err != nil {
						return RandMillisecond(6, 12)
					}
				}
			}
		}
		props := make([]map[string]interface{}, 0)
		for id, item := range items {
			props = append(props, map[string]interface{}{"id": id, "item": item})
		}
		count := len(props)
		i := 0
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			id := props[i]["id"].(int32)
			item := props[i]["item"].(*ItemData)
			am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
				// s
				go func(id int32) {
					_ = CLI.SearchPet(&C2SSearchPet{ItemId: id})
				}(id)
				r := &S2CSearchPet{}
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, r, s3); err == nil && r.Pet != nil {
					go func(r *S2CSearchPet) {
						_ = CLI.AnimalParkGO(&C2SAnimalParkGO{
							PetId: r.Pet.Id,
							X:     r.Pet.PointX,
							Y:     r.Pet.PointY,
						})
					}(r)
					_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CAnimalParkGO{}, s3)
				}
				item.N--
				if item.N <= 0 {
					return 0, 0
				}
				return ms100, 0
			})
			i++
			if i >= count {
				loop = 0
				next = RandMillisecond(600, 1800)
			} else {
				loop = ms100
				next = 0
			}
			return
		})
	}
	for {
		select {
		case <-t.C:
			f()
			t.Reset(RandMillisecond(1800, 3600)) // 30 ~ 60 分钟
		case <-ctx.Done():
			return
		}
	}
}

// EnterAnimalPark 寻找宠物
func (c *Connect) EnterAnimalPark() error {
	body, err := proto.Marshal(&C2SEnterAnimalPark{})
	if err != nil {
		return err
	}
	log.Println("[C][EnterAnimalPark]")
	return c.send(19073, body)
}

// LeaveAnimalPark 寻找宠物
func (c *Connect) LeaveAnimalPark() error {
	body, err := proto.Marshal(&C2SLeaveAnimalPark{})
	if err != nil {
		return err
	}
	log.Println("[C][LeaveAnimalPark]")
	return c.send(19075, body)
}

// SearchPet 寻找宠物
func (c *Connect) SearchPet(act *C2SSearchPet) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][SearchPet] item_id=%v", act.ItemId)
	return c.send(19077, body)
}

// AnimalParkGO 抓捕宠物
func (c *Connect) AnimalParkGO(act *C2SAnimalParkGO) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][AnimalParkGO] pet_id=%v x=%v y=%v", act.PetId, act.X, act.Y)
	return c.send(19085, body)
}

////////////////////////////////////////////////////////////

func (x *S2CEnterAnimalPark) ID() uint16 {
	return 19074
}

func (x *S2CEnterAnimalPark) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][EnterAnimalPark] tag=%v tag_msg=%s pet=%v buff=%v", x.Tag, GetTagMsg(x.Tag), x.Pet, x.Buff)
}

////////////////////////////////////////////////////////////

func (x *S2CSearchPet) ID() uint16 {
	return 19078
}

func (x *S2CSearchPet) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SearchPet] tag=%v tag_msg=%s pet=%v buff=%v", x.Tag, GetTagMsg(x.Tag), x.Pet, x.Buff)
}

////////////////////////////////////////////////////////////

func (x *S2CAnimalParkGO) ID() uint16 {
	return 19086
}

func (x *S2CAnimalParkGO) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][AnimalParkGO] tag=%v tag_msg=%s times=%v del_pets=%v", x.Tag, GetTagMsg(x.Tag), x.Times, x.DelPets)
}

////////////////////////////////////////////////////////////

func (x *S2CLeaveAnimalPark) ID() uint16 {
	return 19076
}

func (x *S2CLeaveAnimalPark) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][LeaveAnimalPark] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CSearchRecord) ID() uint16 {
	return 19079
}

func (x *S2CSearchRecord) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SearchRecord] search_record=%v", x.SearchRecord)
}

////////////////////////////////////////////////////////////

func (x *S2CAnimalParkCatch) ID() uint16 {
	return 19088
}

func (x *S2CAnimalParkCatch) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][AnimalParkCatch] tag=%v tag_msg=%s pet_id=%v drop=%v", x.Tag, GetTagMsg(x.Tag), x.PetId, x.Drop)
}
