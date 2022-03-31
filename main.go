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
				(&mhyc.S2CRoleTask{}).Message(data)
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
		}()

		//go mhyc.Everyday(ctx)
		//go mhyc.Mail(ctx)
		//go mhyc.AFK()
		//go mhyc.StageFight()
		//go mhyc.FamilyJJC()
		//go mhyc.EnterAnimalPark()
		//go mhyc.XianDianXDSW()
		go mhyc.XianDianSSSL()
		//go mhyc.XianDianXDXS()
		//go mhyc.BossPersonal()
		//go mhyc.BossVIP()
		//go mhyc.BossMulti()
		//go mhyc.XuanShangBoss()
		//go mhyc.BossGlobal()
		//go mhyc.BossHome()
		//go mhyc.BossXLD()
		go mhyc.BossXSD()
		//go mhyc.BossXMD()
		//go mhyc.FuBen(ctx)
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
				cancel()
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
