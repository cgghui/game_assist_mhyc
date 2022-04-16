package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func actJZLDZTime() time.Duration {
	cur := time.Now()
	ast := time.Date(cur.Year(), cur.Month(), cur.Day(), 12, 00, 3, 0, time.Local)
	if cur.Before(ast) {
		return ast.Sub(cur)
	}
	if cur.Before(ast.Add(9 * time.Hour)) {
		return 0
	}
	return TomorrowDuration(43203 * time.Second)
}

// HuoDongZJLDZ 活动<家族地战>
func HuoDongZJLDZ(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	f1 := func() time.Duration {
		if td := actJZLDZTime(); td != 0 {
			return td
		}
		week := time.Now().Weekday()
		if is := week == time.Wednesday || week == time.Saturday; !is {
			return TomorrowDuration(RandMillisecond(600, 1800))
		}
		Receive.Action(CLI.GetCityWarData)
		data := &S2CCityWarData{}
		_ = Receive.Wait(data, s3)
		Receive.Action(CLI.GetCityWarChooseItem)
		_ = Receive.Wait(&S2CGetCityWarChooseItem{}, s3)
		roleFamilyId := RoleInfo.Get("FamilyId").String()
		isWarChoose := false
		for _, city := range data.Data.CityData {
			if city.CityState == 0 {
				continue
			}
			for _, family := range city.Familys {
				if family.FamilyId == roleFamilyId {
					isWarChoose = true
					break
				}
			}
			if isWarChoose {
				break
			}
		}
		if isWarChoose {
			return TomorrowDuration(RandMillisecond(600, 1800))
		}
		for _, city := range data.Data.CityData {
			if city.CityState == 1 && len(city.Familys) == 0 {
				go func() {
					_ = CLI.CityWarChoose(city.CityId)
				}()
				_ = Receive.Wait(&S2CCityWarChoose{}, s3)
				isWarChoose = true
				break
			}
		}
		if isWarChoose {
			return TomorrowDuration(RandMillisecond(600, 1800))
		}
		return RandMillisecond(10, 15)
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

// GetCityWarData 信息
func (c *Connect) GetCityWarData() error {
	body, err := proto.Marshal(&C2SGetCityWarData{})
	if err != nil {
		return err
	}
	log.Println("[C][GetCityWarData]")
	return c.send(25501, body)
}

func (x *S2CCityWarData) ID() uint16 {
	return 25502
}

// Message S2CCityWarData Code:25502
func (x *S2CCityWarData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][CityWarData] data=%v", x.Data)
}

////////////////////////////////////////////////////////////

// GetCityWarChooseItem 信息
func (c *Connect) GetCityWarChooseItem() error {
	body, err := proto.Marshal(&C2SGetCityWarChooseItem{})
	if err != nil {
		return err
	}
	log.Println("[C][GetCityWarChooseItem]")
	return c.send(25503, body)
}

func (x *S2CGetCityWarChooseItem) ID() uint16 {
	return 25504
}

// Message S2CGetCityWarChooseItem Code:25504
func (x *S2CGetCityWarChooseItem) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetCityWarChooseItem] id=%v", x.Id)
}

////////////////////////////////////////////////////////////

// CityWarChoose 信息
func (c *Connect) CityWarChoose(id int32) error {
	body, err := proto.Marshal(&C2SCityWarChoose{Id: id})
	if err != nil {
		return err
	}
	log.Printf("[C][CityWarChoose] id=%d", id)
	return c.send(25507, body)
}

func (x *S2CCityWarChoose) ID() uint16 {
	return 25504
}

// Message S2CCityWarChoose Code:25508
func (x *S2CCityWarChoose) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][CityWarChoose] tag=%v id=%v", x.Tag, x.Id)
}
