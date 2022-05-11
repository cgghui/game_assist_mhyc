package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/cgghui/game_assist_mhyc/mhyc"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

func main() {

	go func() {

		web := gin.New()

		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)

		web.GET("/get_task_status", func(ctt *gin.Context) {
			text := make([]string, 0)
			for _, am := range mhyc.ActionManageList {
				if am.Name == "" {
					continue
				}
				text = append(text, "current: "+am.Name+" running time: "+time.Since(am.Sr).String())
			}
			text = append(text, "-----------------------")
			i := len(mhyc.ActionRunningHistoryList) - 1
			for ; i >= 0; i-- {
				am := mhyc.ActionRunningHistoryList[i]
				text = append(text, am.RunningTime+" "+am.Name+" "+am.TakeUpTime.String())
			}
			ctt.String(http.StatusOK, strings.Join(text, "\n"))
		})

		web.GET("/get_role_data", func(ctt *gin.Context) {
			info := mhyc.RoleInfo.GetAll()
			ctt.AsciiJSON(http.StatusOK, info)
		})

		web.GET("/get_user_bag_data", func(ctt *gin.Context) {
			info := mhyc.UserBag.GetAll()
			ctt.AsciiJSON(http.StatusOK, info)
		})

		s := &http.Server{
			Addr:         ":9292",
			Handler:      web,
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
		}
		_ = s.ListenAndServe()
	}()

	tm := time.NewTimer(time.Millisecond)

	reStart := func(r int) {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for range t.C {
			//fmt.Printf("\r系统将在%d秒后，重新运行", r)
			r--
			if r == 0 {
				tm.Reset(time.Millisecond)
				return
			}
		}
	}

	listenAction := []mhyc.HandleMessage{
		&mhyc.S2CBagChange{},
		&mhyc.ItemFly{},
		&mhyc.S2CBattlefieldReport{},
		&mhyc.S2CTeamInstanceGetReport{},
		&mhyc.S2CServerTime{},
		&mhyc.S2CRedState{},
		&mhyc.S2CStartFight{},
		&mhyc.S2CNotice{},
		&mhyc.S2CPlayerEnterMap{},
		&mhyc.S2CPlayerLeaveMap{},
		&mhyc.S2CChangeMap{},
		&mhyc.S2CMonsterEnterMap{},
		&mhyc.S2CUpdateAmount{},
		&mhyc.S2CWestExp{},
		&mhyc.S2CHomeBossInfo{},
		&mhyc.S2CPlayerMove{},
	}

	for range tm.C {

		mhyc.Init()

		mhyc.CTX, mhyc.CancelFunc = context.WithCancel(context.Background())

		wg := &sync.WaitGroup{}

		thread := func(f ...func()) {
			for i := range f {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					f[i]()
				}(i)
			}
		}

		threadCtx := func(f ...func(ctx context.Context)) {
			for i := range f {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					f[i](mhyc.CTX)
				}(i)
			}
		}

		var err error
		var session *mhyc.Client
		session, err = mhyc.NewClient("549f82f51e9562594e5572e487585160", "19c3e57e2544f42b30b21104d24b2a94")
		if err != nil {
			log.Printf("session: %v", err)
			go reStart(60)
			continue
		}
		var cli *mhyc.Connect
		if cli, err = session.Connect(mhyc.CTX); err != nil {
			log.Printf("connect: %v", err)
			go reStart(60)
			continue
		}

		mhyc.CLI = cli

		thread(func() {
			mhyc.ListenMessage(mhyc.CTX, &mhyc.Pong{})
		})

		thread(func() {

			// role info
			info := mhyc.Receive.CreateChannel(&mhyc.S2CRoleInfo{})
			info.Call.Message(<-info.Wait())
			thread(func() {
				defer info.Close()
				for {
					select {
					case data := <-info.Wait():
						(&mhyc.S2CRoleInfo{}).Message(data)
					case <-mhyc.CTX.Done():
						return
					}
				}
			})
			// user bag
			mhyc.Receive.Action(cli.UserBag)
			ubag := mhyc.Receive.CreateChannel(&mhyc.S2CUserBag{})
			ubag.Call.Message(<-ubag.Wait())
			thread(func() {
				defer ubag.Close()
				for {
					select {
					case data := <-ubag.Wait():
						(&mhyc.S2CUserBag{}).Message(data)
					case <-mhyc.CTX.Done():
						return
					}
				}
			})
			//
			mhyc.Receive.Action(cli.LoginEnd)
			//
			thread(
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[0], func(data []byte) {
						listenAction[0].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[1], func(data []byte) {
						listenAction[1].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[2], func(data []byte) {
						listenAction[2].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[3], func(data []byte) {
						listenAction[3].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[4], func(data []byte) {
						listenAction[4].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[5], func(data []byte) {
						listenAction[5].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[6], func(data []byte) {
						listenAction[6].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[7], func(data []byte) {
						listenAction[7].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[8], func(data []byte) {
						listenAction[8].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[9], func(data []byte) {
						listenAction[9].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[10], func(data []byte) {
						listenAction[10].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[11], func(data []byte) {
						listenAction[11].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[12], func(data []byte) {
						listenAction[12].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[13], func(data []byte) {
						listenAction[13].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[14], func(data []byte) {
						listenAction[14].Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(mhyc.CTX, listenAction[15], func(data []byte) {
						listenAction[15].Message(data)
					})
				},
			)

			threadCtx(
				mhyc.Everyday,
				mhyc.Mail,
				mhyc.AFK,
				mhyc.StageFight,
				mhyc.FamilyJJC,
				mhyc.EnterAnimalPark,
				mhyc.XianDianXDSW,
				mhyc.XianDianSSSL,
				mhyc.XianDianXDXS,
				mhyc.XuanShangBoss,
				mhyc.BossPersonal,
				mhyc.BossVIP,
				mhyc.BossXYCM,
				mhyc.BossXYCMGo,
				mhyc.BossMulti,
				mhyc.BossHome,
				mhyc.BossXLD,
				mhyc.BossXSD,
				mhyc.BossXMD,
				mhyc.BossHLTJ,
				mhyc.BossBDJJ,
				mhyc.WorldBoss,
				mhyc.KuaFu,
				mhyc.FuBen,
				mhyc.HuoDongSBHS,
				mhyc.HuoDongBusiness,
				mhyc.HuoDongSSZN,
				mhyc.HuoDongXS,
				mhyc.HuoDongZJLDZ,
				mhyc.JJC,
				mhyc.WZZB,
				mhyc.Buy,
				mhyc.JXSC, //JXSC_Join_Times
				//mhyc.TJZC,
			)
			//go mhyc.ShenYu(ctx)
		})

		thread(func() {
			for {
				var message []byte
				if _, message, err = cli.Conn.ReadMessage(); err != nil {
					mhyc.CancelFunc()
					mhyc.Receive.Close()
					log.Printf("read: %v", err)
					return
				}
				var id uint16
				err = binary.Read(bytes.NewBuffer(message[2:4]), binary.BigEndian, &id)
				if err != nil {
					mhyc.CancelFunc()
					mhyc.Receive.Close()
					log.Printf("recv: %v", err)
					return
				}
				if id == 6 {
					_ = cli.Conn.Close()
					mhyc.CancelFunc()
					mhyc.Receive.Close()
					return
				}
				mhyc.Receive.Notify(id, message[4:])
			}
		})
		wg.Wait()
		go reStart(300)
	}

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
