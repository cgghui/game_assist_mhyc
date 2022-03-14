package main

import (
	"bytes"
	"encoding/binary"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"study/mhyc"
	"time"
)

func main() {

	var err error

	cli := &mhyc.Client{}

	cli.Conn, _, err = websocket.DefaultDialer.Dial("wss://allws.huanlingxiuxian.com/300914:30042", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func() {
		_ = cli.Conn.Close()
	}()

	err = cli.Login(&mhyc.C2SLogin{
		AreaId:       303265,
		AccountId:    52097426,
		Token:        "0332b6df-e353-4ee4-92bc-d0bd1017dbe7",
		UserId:       53927069,
		Fcm:          2,
		LoginPf:      "h5",
		CheckWordUrl: "",
		CodeVersion:  30283,
		ExcelVersion: 29812,
	})
	if err != nil {
		log.Fatal("login:", err)
	}

	go func() {
		s := 3 * time.Second
		t := time.NewTimer(s)
		for range t.C {
			_ = cli.Ping()
			t.Reset(s)
		}
	}()

	go func() {
		time.Sleep(time.Second)
		_ = cli.GetMailAttach(mhyc.DefineGetMailAttachAll)
		time.Sleep(time.Second)
		_ = cli.ActGiftNewReceive(mhyc.DefineGiftRechargeEveryDay)
		time.Sleep(time.Second)
		_ = cli.Respect(mhyc.DefineRespect)
	}()

	//go func() {
	//	t := time.NewTimer(time.Second)
	//	for range t.C {
	//		_ = cli.MailList()
	//		t.Reset(time.Second)
	//	}
	//}()

	go func() {
		for {
			var message []byte
			if _, message, err = cli.Conn.ReadMessage(); err != nil {
				log.Println("read:", err)
				return
			}
			var id int16
			err = binary.Read(bytes.NewBuffer(message[2:4]), binary.BigEndian, &id)
			if err != nil {
				log.Printf("recv: %v", err)
				continue
			}
			if _, ok := mhyc.PCK[id]; !ok {
				log.Printf("recv: id[%d] manage func non-existent", id)
				continue
			}
			go func(id int16, message []byte) {
				mhyc.PCK[id].Message(message)
			}(id, message[4:])
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
