package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// FuBen 副本
func FuBen(ctx context.Context) {
	//
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	f1 := func() time.Duration {
		isEnd := 0
		//
		ims := &S2CInstanceMaterialSweep{}
		Receive.Action(CLI.InstanceMaterialSweep)
		if _ = Receive.Wait(ims, s3); ims.Tag != 0 {
			isEnd += 1
		}
		//
		isl1 := &S2CInstanceSLSweep{}
		Receive.Action(CLI.InstanceSLSweep1)
		if _ = Receive.Wait(isl1, s3); isl1.Tag != 0 {
			isEnd += 1
		}
		//
		isl2 := &S2CInstanceSLSweep{}
		Receive.Action(CLI.InstanceSLSweep2)
		if _ = Receive.Wait(isl2, s3); isl2.Tag != 0 {
			isEnd += 1
		}
		//
		if isEnd >= 3 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		return time.Second
	}
	//
	t2 := time.NewTimer(ms100)
	defer t2.Stop()
	f2 := func() time.Duration {
		enter := &S2CClimbingTowerEnter{}
		go func() {
			_ = CLI.ClimbingTowerEnter(&C2SClimbingTowerEnter{TowerType: 1})
		}()
		if err := Receive.Wait(enter, s3); err != nil {
			return ms500
		}
		for {
			go func() {
				_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: 1, Id: 0})
			}()
			r := &S2CClimbingTowerFight{}
			if err := Receive.Wait(r, s3); err != nil {
				return ms100
			}
			if r.Tag != 0 {
				break
			}
		}
		return ms100
	}
	//
	for {
		select {
		case <-t1.C:
			t1.Reset(f1())
		case <-t2.C:
			t1.Reset(f2())
		case <-ctx.Done():
			return
		}
	}
}

func PaTa() {

}

// InstanceMaterialSweep 副本 材料 一键扫荡
func (c *Connect) InstanceMaterialSweep() error {
	body, err := proto.Marshal(&C2SInstanceMaterialSweep{Id: 0})
	if err != nil {
		return err
	}
	log.Println("[C][InstanceMaterialSweep] ID: 0")
	return c.send(609, body)
}

// InstanceSLSweep1 副本 试炼 一键扫荡
func (c *Connect) InstanceSLSweep1() error {
	body, err := proto.Marshal(&C2SInstanceSLSweep{Id: 1})
	if err != nil {
		return err
	}
	log.Println("[C][InstanceSLSweep] ID: 1")
	return c.send(23015, body)
}

// InstanceSLSweep2 副本 试炼 一键扫荡
func (c *Connect) InstanceSLSweep2() error {
	body, err := proto.Marshal(&C2SInstanceSLSweep{Id: 2})
	if err != nil {
		return err
	}
	log.Println("[C][InstanceSLSweep] ID: 2")
	return c.send(23015, body)
}

// ClimbingTowerEnter 进入爬塔场景
func (c *Connect) ClimbingTowerEnter(enter *C2SClimbingTowerEnter) error {
	body, err := proto.Marshal(enter)
	if err != nil {
		return err
	}
	log.Println("[C][ClimbingTowerEnter] TowerType: 1")
	return c.send(22571, body)
}

// ClimbingTowerFight 副本 - 爬塔 - 战斗
func (c *Connect) ClimbingTowerFight(act *C2SClimbingTowerFight) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(22575, body)
}

////////////////////////////////////////////////////////////

func (x *S2CInstanceMaterialSweep) ID() uint16 {
	return 610
}

// Message S2CInstanceMaterialSweep 610
func (x *S2CInstanceMaterialSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][InstanceMaterialSweep] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CInstanceSLSweep) ID() uint16 {
	return 23016
}

// Message S2CInstanceSLSweep 23016
func (x *S2CInstanceSLSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][InstanceSLSweep] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CClimbingTowerEnter) ID() uint16 {
	return 22572
}

// Message S2CClimbingTowerEnter 22572
func (x *S2CClimbingTowerEnter) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ClimbingTowerEnter] tag=%v tower_type=%v", x.Tag, x.TowerType)
}

////////////////////////////////////////////////////////////

func (x *S2CClimbingTowerFight) ID() uint16 {
	return 22576
}

// Message S2CClimbingTowerFight 22576
func (x *S2CClimbingTowerFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ClimbingTowerFight] tag=%v id=%v tower_type=%v", x.Tag, x.Id, x.TowerType)
}
