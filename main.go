package main

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
			info := &mhyc.S2CRoleInfo{}
			_ = mhyc.Receive.Wait(info)
			go func() {
				channel := mhyc.Receive.CreateChannel(&mhyc.S2CRoleInfo{})
				for data := range channel.Wait() {
					info.Message(data)
				}
			}()
			go func() {
				var obj *mhyc.S2CBattlefieldReport
				channel := mhyc.Receive.CreateChannel(&mhyc.S2CBattlefieldReport{})
				for data := range channel.Wait() {
					obj = &mhyc.S2CBattlefieldReport{}
					obj.Message(data)
					_ = cli.EndFight(obj)
				}
			}()
			go func() {
				pong := &mhyc.Pong{}
				channel := mhyc.Receive.CreateChannel(pong)
				for range channel.Wait() {
					pong.Message(nil)
				}
			}()
		}()

		// 用户背包
		func() {
			info := &mhyc.S2CUserBag{}
			mhyc.Receive.Action(cli.UserBag)
			_ = mhyc.Receive.Wait(info)
			go func() {
				channel := mhyc.Receive.CreateChannel(&mhyc.S2CUserBag{})
				for data := range channel.Wait() {
					info.Message(data)
				}
			}()
			go func() {
				channel := mhyc.Receive.CreateChannel(&mhyc.S2CBagChange{})
				for data := range channel.Wait() {
					(&mhyc.S2CBagChange{}).Message(data)
				}
			}()
			go func() {
				channel := mhyc.Receive.CreateChannel(&mhyc.ItemFly{})
				for data := range channel.Wait() {
					(&mhyc.ItemFly{}).Message(data)
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
		go mhyc.Everyday()
		go mhyc.Mail()
		go mhyc.AFK()
		go mhyc.StageFight()
		go mhyc.FamilyJJC()
		go mhyc.EnterAnimalPark()

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
