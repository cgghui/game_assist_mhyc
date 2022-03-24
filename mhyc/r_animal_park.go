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

func init() {
	PCK[19074] = &S2CEnterAnimalPark{}
	PCK[19076] = &S2CLeaveAnimalPark{}
	PCK[19078] = &S2CSearchPet{}
	PCK[19079] = &S2CSearchRecord{}
	PCK[19088] = &S2CAnimalParkCatch{}
	PCK[19086] = &S2CAnimalParkGO{}
}

var animalParkThread = make(chan interface{})
var animalParkAction = make(chan struct{})

// EnterAnimalPark 寻找宠物
func (c *Connect) EnterAnimalPark() {
	go func() {
		animalParkAction <- struct{}{}
		t := time.NewTimer(10 * time.Minute)
		for range t.C {
			animalParkAction <- struct{}{}
			t.Reset(RandMillisecond(300, 900)) // 5 ~ 15 分钟
		}
	}()
	ArrestWait := make(chan struct{})
	SearchWait := make(chan *S2CSearchPet)
	run := func(val interface{}) {
		switch ret := val.(type) {
		case *S2CEnterAnimalPark:
			{
				if ret.Pet != nil {
					for _, pet := range ret.Pet {
						_ = c.animalParkGO(&C2SAnimalParkGO{PetId: pet.Id, X: pet.PointX, Y: pet.PointY})
						<-ArrestWait
					}
				}
				var items = make(map[int32]*ItemData)
				var n = int64(0)
				if item, ok := UserBag.Get(ItemPet500, time.Second); ok {
					items[ItemPet500] = item
					n += items[ItemPet500].N
				}
				if item, ok := UserBag.Get(ItemPet501, time.Second); ok {
					items[ItemPet501] = item
					n += items[ItemPet501].N
				}
				if item, ok := UserBag.Get(ItemPet502, time.Second); ok {
					items[ItemPet502] = item
					n += items[ItemPet502].N
				}
				if n > 200 {
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
						if item, ok := UserBag.Get(ItemPet503, time.Second); ok {
							v := item
							if v.N > 0 {
								_ = c.searchPet(&C2SSearchPet{ItemId: ItemPet503})
								<-SearchWait
							}
						}
					}
					for id, item := range items {
						for i := 0; i < int(item.N); i++ {
							_ = c.searchPet(&C2SSearchPet{ItemId: id})
							r := <-SearchWait
							_ = c.animalParkGO(&C2SAnimalParkGO{PetId: r.Pet.Id, X: r.Pet.PointX, Y: r.Pet.PointY})
							<-ArrestWait
						}
					}
				}
				_ = c.leaveAnimalPark()
			}
		case *S2CAnimalParkGO:
			{
				ArrestWait <- struct{}{}
			}
		case *S2CSearchPet:
			{
				SearchWait <- ret
			}
		case *S2CLeaveAnimalPark:
			return
		}
	}
	for {
		select {
		case <-animalParkAction:
			_ = c.enterAnimalPark()
		case val := <-animalParkThread:
			go run(val)
		}
	}
}

func (c *Connect) enterAnimalPark() error {
	body, err := proto.Marshal(&C2SEnterAnimalPark{})
	if err != nil {
		return err
	}
	return c.send(19073, body)
}

func (c *Connect) leaveAnimalPark() error {
	body, err := proto.Marshal(&C2SLeaveAnimalPark{})
	if err != nil {
		return err
	}
	return c.send(19075, body)
}

// searchPet 寻找宠物
func (c *Connect) searchPet(act *C2SSearchPet) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(19077, body)
}

// animalParkGO 抓捕宠物
func (c *Connect) animalParkGO(act *C2SAnimalParkGO) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(19085, body)
}

func (x *S2CEnterAnimalPark) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][EnterAnimalPark] tag=%v pet=%v buff=%v", x.Tag, x.Pet, x.Buff)
	animalParkThread <- x
}

func (x *S2CSearchPet) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SearchPet] tag=%v pet=%v buff=%v", x.Tag, x.Pet, x.Buff)
	animalParkThread <- x
}

func (x *S2CAnimalParkGO) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][AnimalParkGO] tag=%v times=%v del_pets=%v", x.Tag, x.Times, x.DelPets)
	animalParkThread <- x
}

func (x *S2CLeaveAnimalPark) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][LeaveAnimalPark] tag=%v", x.Tag)
	animalParkThread <- x
}

func (x *S2CSearchRecord) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SearchRecord] search_record=%v", x.SearchRecord)
}

func (x *S2CAnimalParkCatch) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][AnimalParkCatch] tag=%v pet_id=%v drop=%v", x.Tag, x.PetId, x.Drop)
}
