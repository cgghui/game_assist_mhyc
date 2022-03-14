package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
)

var PCK = map[int16]Handle{
	2:     &S2CLogin{},
	3:     &S2CRoleInfo{},
	23:    &Pong{},
	51:    &S2CChangeMap{},
	12012: &S2CActGiftNewReceive{},
	14:    &S2CRespect{},
	1001:  &S2CServerTime{},
	403:   &S2CNewChatMsg{},
	441:   &S2CMailList{},
	445:   &S2CGetMailAttach{},
	446:   &S2CNewMail{},
}

type Handle interface {
	Message([]byte)
}

// Message S2CLogin 登录
func (x *S2CLogin) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [Login] %v", err)
		return
	}
	if x.UserId == 0 {
		log.Printf("recv: [Login] 无法登录，请更换【token】")
		return
	}
	log.Printf("recv: [Login] 登录成功 用户ID: %d", x.UserId)
	return
}

// Message S2CRoleInfo 角色信息
func (x *S2CRoleInfo) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [RoleInfo] %v", err)
		return
	}
	log.Println("RoleInfo")
	return
}

func (x *Pong) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [Pong] %v", err)
		return
	}
	log.Println("Pong")
	return
}

func (x *S2CChangeMap) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [ChangeMap] %v", err)
		return
	}
	log.Println("ChangeMap")
	return
}

func (x *S2CActGiftNewReceive) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [ActGiftNewReceive] %v", err)
		return
	}
	log.Println("ActGiftNewReceive")
	return
}

func (x *S2CRespect) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [Respect] %v", err)
		return
	}
	log.Println("Respect")
	return
}

func (x *S2CServerTime) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [ServerTime] %v", err)
		return
	}
	log.Printf("ServerTime: %d", x.T)
	return
}

func (x *S2CNewChatMsg) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [NewChatMsg] %v", err)
		return
	}
	log.Printf("NewChatMsg: [%s] %s", x.Chatmessage.SenderNick, x.Chatmessage.Content)
	return
}

func (x *S2CMailList) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [MailList] %v", err)
		return
	}
	log.Println("MailList")
	return
}

func (x *S2CGetMailAttach) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [GetMailAttach] %v", err)
		return
	}
	log.Printf("GetMailAttach: %d", x.Tag)
	return
}

func (x *S2CNewMail) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("recv: [NewMail] %v", err)
		return
	}
	log.Println("NewMail")
	return
}
