package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func init() {
	PCK[441] = &S2CMailList{}
	PCK[445] = &S2CGetMailAttach{}
	PCK[446] = &S2CNewMail{}
	PCK[448] = &S2CDelMail{}
}

var mailListThread = make(chan interface{})
var mailListAction = make(chan struct{})

// Mail 邮件
func (c *Connect) Mail() {
	go func() {
		mailListAction <- struct{}{}
		t := time.NewTimer(time.Minute)
		for range t.C {
			mailListAction <- struct{}{}
			t.Reset(time.Minute)
		}
	}()
	for range mailListAction {
		// 普通邮件
		_ = c.mailList(DefineMailListOrdinary)
		if list := (<-mailListThread).(*S2CMailList); len(list.MailList) > 0 {
			if list.MailList[0].IsReceive == 1 && list.MailList[0].IsRead == 1 {
				_ = c.getMailAttach(DefineGetMailAttachOrdinary)
				<-mailListThread
			}
			_ = c.delMail(DefineDelMailOrdinary)
			<-mailListThread
		}
		// 活动邮件
		_ = c.mailList(DefineMailListActivity)
		if list := (<-mailListThread).(*S2CMailList); len(list.MailList) > 0 {
			if list.MailList[0].IsReceive == 1 && list.MailList[0].IsRead == 1 {
				_ = c.getMailAttach(DefineGetMailAttachActivity)
				<-mailListThread
			}
			_ = c.delMail(DefineDelMailActivity)
			<-mailListThread
		}
	}
}

// mailList 邮件列表
func (c *Connect) mailList(act *C2SMailList) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(440, body)
}

// getMailAttach 领取邮件附件
func (c *Connect) getMailAttach(act *C2SGetMailAttach) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(444, body)
}

// delMail 删除已读邮件
func (c *Connect) delMail(act *C2SDelMail) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(457, body)
}

func (x *S2CMailList) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][MailList] list=%v", x.MailList)
	mailListThread <- x
	return
}

func (x *S2CGetMailAttach) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetMailAttach] tag=%v", x.Tag)
	mailListThread <- x
	return
}

func (x *S2CDelMail) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][DelMail] tag=%v", x.Tag)
	mailListThread <- x
	return
}

func (x *S2CNewMail) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][NewMail]")
	mailListThread <- struct{}{}
	return
}
