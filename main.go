package main

// 仙宗 仙殿 仙宗悬赏

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/cgghui/game_assist_mhyc/mhyc"
	"log"
	"os"
	"os/signal"
	"time"
)

func init() {
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	var session *mhyc.Client
	session, err = mhyc.NewClient("549f82f51e9562594e5572e487585160", "19c3e57e2544f42b30b21104d24b2a94")
	if err != nil {
		log.Fatal("session:", err)
	}
	var cli *mhyc.Connect
	if cli, err = session.Connect(ctx); err != nil {
		log.Fatal("connect:", err)
	}

	mhyc.CLI = cli

	go func() {

		go mhyc.ListenMessage(ctx, &mhyc.Pong{})

		// 角色信息
		func() {
			// role info
			info := mhyc.Receive.CreateChannel(&mhyc.S2CRoleInfo{})
			info.Call.Message(<-info.Wait())
			go func() {
				defer info.Close()
				for {
					select {
					case data := <-info.Wait():
						(&mhyc.S2CRoleInfo{}).Message(data)
					case <-ctx.Done():
						return
					}
				}
			}()
			// user bag
			mhyc.Receive.Action(cli.UserBag)
			ubag := mhyc.Receive.CreateChannel(&mhyc.S2CUserBag{})
			ubag.Call.Message(<-ubag.Wait())
			go func() {
				defer ubag.Close()
				for {
					select {
					case data := <-ubag.Wait():
						(&mhyc.S2CUserBag{}).Message(data)
					case <-ctx.Done():
						return
					}
				}
			}()
			//
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CBagChange{}, func(data []byte) {
				(&mhyc.S2CBagChange{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.ItemFly{}, func(data []byte) {
				(&mhyc.ItemFly{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CRoleTask{}, func(data []byte) {
				task := &mhyc.S2CRoleTask{}
				task.Message(data)
				for i := range task.Task {
					go func(t, id int32) {
						mhyc.Fight.Lock()
						defer mhyc.Fight.Unlock()
						go func(t, id int32) {
							_ = cli.GetTaskPrize(&mhyc.C2SGetTaskPrize{TaskType: t, Multi: 1, TaskId: id})
						}(t, id)
						_ = mhyc.Receive.Wait(&mhyc.S2CGetTaskPrize{}, 3*time.Second)
					}(task.Task[i].T, task.Task[i].Id)
				}
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CBattlefieldReport{}, func(data []byte) {
				(&mhyc.S2CBattlefieldReport{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CTeamInstanceGetReport{}, func(data []byte) {
				(&mhyc.S2CTeamInstanceGetReport{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CServerTime{}, func(data []byte) {
				(&mhyc.S2CServerTime{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CRedState{}, func(data []byte) {
				(&mhyc.S2CRedState{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CStartFight{}, func(data []byte) {
				(&mhyc.S2CStartFight{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CNotice{}, func(data []byte) {
				(&mhyc.S2CNotice{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CPlayerEnterMap{}, func(data []byte) {
				(&mhyc.S2CPlayerEnterMap{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CPlayerLeaveMap{}, func(data []byte) {
				(&mhyc.S2CPlayerLeaveMap{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CChangeMap{}, func(data []byte) {
				(&mhyc.S2CChangeMap{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CMonsterEnterMap{}, func(data []byte) {
				(&mhyc.S2CMonsterEnterMap{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CUpdateAmount{}, func(data []byte) {
				(&mhyc.S2CUpdateAmount{}).Message(data)
			})
			go mhyc.ListenMessageCall(ctx, &mhyc.S2CWestExp{}, func(data []byte) {
				(&mhyc.S2CWestExp{}).Message(data)
			})
		}()

		go mhyc.Everyday(ctx)
		go mhyc.Mail(ctx)
		go mhyc.AFK(ctx)
		go mhyc.StageFight(ctx)
		go mhyc.FamilyJJC(ctx)
		go mhyc.EnterAnimalPark(ctx)
		go mhyc.XianDianXDSW(ctx)
		go mhyc.XianDianSSSL(ctx)
		go mhyc.XianDianXDXS(ctx)
		go mhyc.BossPersonal(ctx)
		go mhyc.BossVIP(ctx)
		go mhyc.BossXYCM(ctx)
		go mhyc.BossMulti(ctx)
		go mhyc.XuanShangBoss(ctx)
		go mhyc.BossHome(ctx)
		go mhyc.BossXLD(ctx)
		go mhyc.BossXSD(ctx)
		go mhyc.BossXMD(ctx)
		go mhyc.BossHLTJ(ctx)
		go mhyc.BossBDJJ(ctx)
		//go mhyc.WorldBoss(ctx)
		go mhyc.KuaFu(ctx)
		go mhyc.FuBen(ctx)
		go mhyc.HuoDongSBHS(ctx)
		go mhyc.HuoDongBusiness(ctx)
		go mhyc.JJC(ctx)
		go mhyc.WZZB(ctx)
		go mhyc.HuoDongXS(ctx)
		//go mhyc.ShenYu(ctx)
		//go mhyc.JXSC(ctx)
		//go mhyc.HuoDongZJLDZ(ctx)
	}()

	go func() {
		for {
			var message []byte
			if _, message, err = cli.Conn.ReadMessage(); err != nil {
				cancel()
				mhyc.Receive.Close()
				log.Println("read:", err)
				return
			}
			var id uint16
			err = binary.Read(bytes.NewBuffer(message[2:4]), binary.BigEndian, &id)
			if err != nil {
				cancel()
				log.Printf("recv: %v", err)
				continue
			}
			mhyc.Receive.Notify(id, message[4:])
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
