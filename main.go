package main

// 仙宗 仙殿 仙宗悬赏

import (
	"bytes"
	"encoding/binary"
	"github.com/cgghui/game_assist_mhyc/mhyc"
	"log"
	"os"
	"os/signal"
)

func init() {
}

func main() {
	var err error
	var session *mhyc.Client
	session, err = mhyc.NewClient("66c9351a98e6809c52a1e36a13af3917", "19c3e57e2544f42b30b21104d24b2a94")
	if err != nil {
		log.Fatal("session:", err)
	}
	var cli *mhyc.Connect
	if cli, err = session.Connect(); err != nil {
		log.Fatal("connect:", err)
	}

	mhyc.CLI = cli

	go func() {

		// 角色信息
		func() {
			// role info
			info := mhyc.Receive.CreateChannel(&mhyc.S2CRoleInfo{})
			info.Call.Message(<-info.Wait())
			// user bag
			mhyc.Receive.Action(cli.UserBag)
			ubag := mhyc.Receive.CreateChannel(&mhyc.S2CUserBag{})
			ubag.Call.Message(<-ubag.Wait())
			//
			bgcg := mhyc.Receive.CreateChannel(&mhyc.S2CBagChange{})
			itfy := mhyc.Receive.CreateChannel(&mhyc.ItemFly{})
			//
			pong := mhyc.Receive.CreateChannel(&mhyc.Pong{})
			task := mhyc.Receive.CreateChannel(&mhyc.S2CRoleTask{})
			bfrp := mhyc.Receive.CreateChannel(&mhyc.S2CBattlefieldReport{})
			srtm := mhyc.Receive.CreateChannel(&mhyc.S2CServerTime{})
			redc := mhyc.Receive.CreateChannel(&mhyc.S2CRedState{})
			go func() {
				for {
					select {
					case data := <-info.Wait():
						go (&mhyc.S2CRoleInfo{}).Message(data)
					case data := <-ubag.Wait():
						go (&mhyc.S2CUserBag{}).Message(data)
					case data := <-bgcg.Wait():
						go (&mhyc.S2CBagChange{}).Message(data)
					case data := <-itfy.Wait():
						go (&mhyc.ItemFly{}).Message(data)
					case <-pong.Wait():
						go (&mhyc.Pong{}).Message(nil)
					case <-task.Wait():
						go (&mhyc.S2CRoleTask{}).Message(nil)
					case data := <-bfrp.Wait():
						go (&mhyc.S2CBattlefieldReport{}).Message(data)
					case <-srtm.Wait():
						go (&mhyc.S2CServerTime{}).Message(nil)
					case <-redc.Wait():
						go (&mhyc.S2CRedState{}).Message(nil)
					}
				}
			}()
		}()

		//t := time.NewTimer(time.Second)
		//run := func() {
		//	for range t.C {
		//		_ = cli.StageFight()
		//		t.Reset(time.Second)
		//	}
		//}
		//run()
		//
		//go mhyc.BossPersonal()
		//go mhyc.BossVIP()
		go mhyc.BossMulti()
		//go mhyc.XianDianXDSW()
		//go mhyc.XianDianSSSL()
		//go mhyc.XianDianXDXS()
		//go mhyc.Everyday()
		//go mhyc.Mail()
		//go mhyc.AFK()
		//go mhyc.StageFight()
		//go mhyc.FamilyJJC()
		//go mhyc.EnterAnimalPark()

		//wg := &sync.WaitGroup{}
		//wg.Add(2)
		//
		//action := []func(){
		//
		//	func() {
		//
		//	},
		//
		//	func() {
		//		//_ = cli.HuanLingList()
		//		//_ = cli.GetActTimestamp(&mhyc.C2SGetActTimestamp{ActId: 14})
		//		//_ = cli.GetActTimestamp(&mhyc.C2SGetActTimestamp{ActId: 15})
		//		//_ = cli.GetPetAMergeInfo()
		//		//_ = cli.GetAllEquipData()
		//		//_ = cli.PlayerPractice()
		//		//_ = cli.GetEquipData(&mhyc.C2SGetEquipData{FuncId: 15001, ObjId: 0})
		//		//_ = cli.Beasts()
		//		//_ = cli.GetHeroList()
		//		//_ = cli.GetAlienData()
		//		//_ = cli.YJInfo()
		//		//_ = cli.SLGetData()
		//		//_ = cli.NewStory()
		//		//_ = cli.StagePrize()
		//		//_ = cli.RoleInfo()
		//		//_ = cli.LoginEnd()
		//		//_ = cli.GetActTask(mhyc.DefineGetActTask11002)
		//		//_ = cli.AFKGetBuyInfo()
		//		//_ = cli.WeddingInsInvite()
		//		//_ = cli.ClimbingTowerEnter(mhyc.DefineClimbingTowerEnter5)
		//		//_ = cli.GetActXunBaoInfo(mhyc.DefineXunBaoInfo501)
		//	},
		//
		//	//func() {
		//	//	go func() {
		//	//		t := time.NewTimer(time.Second)
		//	//		for range t.C {
		//	//			_ = cli.RealmTask() // 修仙 - 境界 任务
		//	//			t.Reset(time.Second)
		//	//		}
		//	//	}()
		//	//},
		//	//func() {
		//	//	t := time.NewTimer(100 * time.Millisecond)
		//	//	run := func() {
		//	//		for range t.C {
		//	//			_ = cli.GetHistoryTaskPrize()
		//	//			recv := <-mhyc.ChanBox.GetHistoryTaskPrize
		//	//			log.Printf("[主线奖励] tag=%v %v", recv.Tag, recv)
		//	//			if recv.Tag == 5043 {
		//	//				break
		//	//			}
		//	//			t.Reset(100 * time.Millisecond)
		//	//		}
		//	//	}
		//	//	go func() {
		//	//		for {
		//	//			wg.Wait()
		//	//			run()
		//	//		}
		//	//	}()
		//	//},
		//	//
		//
		//	//func() {
		//	//	go func() {
		//	//		for {
		//	//			_ = cli.GetHistoryTaskPrize()
		//	//		}
		//	//	}()
		//	//},
		//}
		//
		//i := 0
		//t := time.NewTimer(time.Second)
		//for range t.C {
		//	action[i]()
		//	i++
		//	if i >= len(action) {
		//		t.Stop()
		//		return
		//	}
		//	t.Reset(time.Second)
		//}
	}()

	go func() {
		for {
			var message []byte
			if _, message, err = cli.Conn.ReadMessage(); err != nil {
				log.Println("read:", err)
				return
			}
			var id uint16
			err = binary.Read(bytes.NewBuffer(message[2:4]), binary.BigEndian, &id)
			if err != nil {
				log.Printf("recv: %v", err)
				continue
			}
			mhyc.Receive.Notify(id, message[4:])
			//if _, ok := mhyc.PCK[id]; !ok {
			//	log.Printf("recv: id[%d] manage func non-existent", id)
			//	continue
			//}
			//go mhyc.PCK[id].Message(message[4:])
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}

// 发送
// n = a.encode(i).finish()

// 接收
// o = u.decode(n)

//S2CUserBag

//t.prototype.loginRoleInfo
//t.prototype.Proc_S2CRoleInfo
//return e.AttrType = {
