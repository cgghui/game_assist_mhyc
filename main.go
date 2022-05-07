package main

// 仙宗 仙殿 仙宗悬赏

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
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

type Thread struct {
	f func(ctx context.Context)
}

func main() {

	var ctx context.Context
	var cancel context.CancelFunc

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
			for _, am := range mhyc.ActionRunningHistoryList {
				text = append(text, am.RunningTime+" "+am.Name+" "+am.TakeUpTime.String())
			}
			ctt.String(http.StatusOK, strings.Join(text, "\n"))
		})

		web.GET("/get_role_data", func(ctt *gin.Context) {
			info := mhyc.RoleInfo.GetAll()
			ctt.AsciiJSON(http.StatusOK, info)
		})

		s := &http.Server{
			Addr:         "127.0.0.1:9292",
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
			fmt.Printf("\r系统将在%d秒后，重新运行", r)
			r--
			if r == 0 {
				tm.Reset(time.Millisecond)
				fmt.Println("")
				return
			}
		}
	}

	for range tm.C {

		mhyc.Init()

		ctx, cancel = context.WithCancel(context.Background())

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
					f[i](ctx)
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
		if cli, err = session.Connect(ctx); err != nil {
			log.Printf("connect: %v", err)
			go reStart(60)
			continue
		}

		mhyc.CLI = cli

		thread(func() {

			thread(func() {
				mhyc.ListenMessage(ctx, &mhyc.Pong{})
			})

			// role info
			info := mhyc.Receive.CreateChannel(&mhyc.S2CRoleInfo{})
			info.Call.Message(<-info.Wait())
			thread(func() {
				defer info.Close()
				for {
					select {
					case data := <-info.Wait():
						(&mhyc.S2CRoleInfo{}).Message(data)
					case <-ctx.Done():
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
					case <-ctx.Done():
						return
					}
				}
			})
			//
			mhyc.Receive.Action(cli.LoginEnd)
			//
			thread(
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CBagChange{}, func(data []byte) {
						(&mhyc.S2CBagChange{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.ItemFly{}, func(data []byte) {
						(&mhyc.ItemFly{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CBattlefieldReport{}, func(data []byte) {
						(&mhyc.S2CBattlefieldReport{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CTeamInstanceGetReport{}, func(data []byte) {
						(&mhyc.S2CTeamInstanceGetReport{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CServerTime{}, func(data []byte) {
						(&mhyc.S2CServerTime{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CRedState{}, func(data []byte) {
						(&mhyc.S2CRedState{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CStartFight{}, func(data []byte) {
						(&mhyc.S2CStartFight{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CNotice{}, func(data []byte) {
						(&mhyc.S2CNotice{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CPlayerEnterMap{}, func(data []byte) {
						(&mhyc.S2CPlayerEnterMap{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CPlayerLeaveMap{}, func(data []byte) {
						(&mhyc.S2CPlayerLeaveMap{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CChangeMap{}, func(data []byte) {
						(&mhyc.S2CChangeMap{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CMonsterEnterMap{}, func(data []byte) {
						(&mhyc.S2CMonsterEnterMap{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CUpdateAmount{}, func(data []byte) {
						(&mhyc.S2CUpdateAmount{}).Message(data)
					})
				},
				func() {
					mhyc.ListenMessageCall(ctx, &mhyc.S2CWestExp{}, func(data []byte) {
						(&mhyc.S2CWestExp{}).Message(data)
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
				mhyc.JJC,
				mhyc.WZZB,
				mhyc.HuoDongXS,
				mhyc.HuoDongZJLDZ,
				mhyc.Buy,
				//mhyc.TJZC,
				//mhyc.JXSC,
			)
			//go mhyc.ShenYu(ctx)
		})

		thread(func() {
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
					mhyc.Receive.Close()
					log.Printf("recv: %v", err)
					continue
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
