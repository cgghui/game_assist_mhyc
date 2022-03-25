package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

const (
	ItemPet500 int32 = 500 // 缚妖索
	ItemPet501 int32 = 501 // 高级缚妖索
	ItemPet502 int32 = 502 // 神兽号角
	ItemPet503 int32 = 503 // 金铲子
)

func EnterAnimalPark() {
	t := time.NewTimer(ms10)
	for range t.C {
		ret := &S2CEnterAnimalPark{}
		Receive.Action(CLI.EnterAnimalPark)
		if err := Receive.Wait(19074, ret, s3); err != nil {
			continue
		}
		if ret.Pet != nil {
			for _, pet := range ret.Pet {
				go func(r *PasturePet) {
					_ = CLI.AnimalParkGO(&C2SAnimalParkGO{PetId: r.Id, X: r.PointX, Y: r.PointY})
				}(pet)
				_ = Receive.Wait(19086, &S2CAnimalParkGO{}, s3)
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
		if n > 0 {
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
				if item := UserBag.Wait(ItemPet503, s3); item != nil {
					v := item
					if v.N > 0 {
						Receive.Action(func() error {
							return CLI.SearchPet(&C2SSearchPet{ItemId: ItemPet503})
						})
						_ = Receive.Wait(19078, &S2CSearchPet{}, s3)
					}
				}
			}
			for id, item := range items {
				for i := 0; i < int(item.N); i++ {
					// s
					r := &S2CSearchPet{}
					go func(id int32) {
						_ = CLI.SearchPet(&C2SSearchPet{ItemId: id})
					}(id)
					_ = Receive.Wait(19078, r, s3)
					// a
					if r.Pet == nil {
						continue
					}
					go func(r *S2CSearchPet) {
						_ = CLI.AnimalParkGO(&C2SAnimalParkGO{PetId: r.Pet.Id, X: r.Pet.PointX, Y: r.Pet.PointY})
					}(r)
					_ = Receive.Wait(19086, &S2CAnimalParkGO{}, s3)
				}
			}
		}
		Receive.Action(CLI.LeaveAnimalPark)
		_ = Receive.Wait(19076, &S2CLeaveAnimalPark{}, s3)
		t.Reset(RandMillisecond(1800, 3600)) // 30 ~ 60 分钟
	}
}

// EnterAnimalPark 寻找宠物
func (c *Connect) EnterAnimalPark() error {
	body, err := proto.Marshal(&C2SEnterAnimalPark{})
	if err != nil {
		return err
	}
	return c.send(19073, body)
}

// LeaveAnimalPark 寻找宠物
func (c *Connect) LeaveAnimalPark() error {
	body, err := proto.Marshal(&C2SLeaveAnimalPark{})
	if err != nil {
		return err
	}
	return c.send(19075, body)
}

// SearchPet 寻找宠物
func (c *Connect) SearchPet(act *C2SSearchPet) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(19077, body)
}

// AnimalParkGO 抓捕宠物
func (c *Connect) AnimalParkGO(act *C2SAnimalParkGO) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(19085, body)
}

func (x *S2CEnterAnimalPark) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][EnterAnimalPark] tag=%v pet=%v buff=%v", x.Tag, x.Pet, x.Buff)
}

func (x *S2CSearchPet) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SearchPet] tag=%v pet=%v buff=%v", x.Tag, x.Pet, x.Buff)
}

func (x *S2CAnimalParkGO) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][AnimalParkGO] tag=%v times=%v del_pets=%v", x.Tag, x.Times, x.DelPets)
}

func (x *S2CLeaveAnimalPark) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][LeaveAnimalPark] tag=%v", x.Tag)
}

func (x *S2CSearchRecord) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SearchRecord] search_record=%v", x.SearchRecord)
}

func (x *S2CAnimalParkCatch) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][AnimalParkCatch] tag=%v pet_id=%v drop=%v", x.Tag, x.PetId, x.Drop)
}
