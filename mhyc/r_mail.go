package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func Mail() {
	f := func() {
		list := &S2CMailList{}
		// 普通邮件
		go func() {
			_ = CLI.MailList(DefineMailListOrdinary)
		}()
		_ = Receive.Wait(441, list, s3)
		if len(list.MailList) > 0 {
			if list.MailList[0].IsReceive == 1 && list.MailList[0].IsRead == 1 {
				go func() {
					_ = CLI.GetMailAttach(DefineGetMailAttachOrdinary)
				}()
				_ = Receive.Wait(445, &S2CGetMailAttach{}, s3)
			}
			go func() {
				_ = CLI.DelMail(DefineDelMailOrdinary)
			}()
			_ = Receive.Wait(448, &S2CDelMail{}, s3)
		}
		// 活动邮件
		go func() {
			_ = CLI.MailList(DefineMailListActivity)
		}()
		_ = Receive.Wait(441, list, s3)
		if len(list.MailList) > 0 {
			if list.MailList[0].IsReceive == 1 && list.MailList[0].IsRead == 1 {
				go func() {
					_ = CLI.GetMailAttach(DefineGetMailAttachActivity)
				}()
				_ = Receive.Wait(445, &S2CGetMailAttach{}, s3)
			}
			go func() {
				_ = CLI.DelMail(DefineDelMailActivity)
			}()
			_ = Receive.Wait(448, &S2CDelMail{}, s3)
		}
	}
	go func() {
		_ = Receive.Wait(446, &S2CNewMail{})
		f()
	}()
	t := time.NewTimer(ms10)
	for range t.C {
		f()
		t.Reset(RandMillisecond(300, 600)) // 5 ~ 10 分钟
	}
}

// MailList 邮件列表
func (c *Connect) MailList(act *C2SMailList) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(440, body)
}

// GetMailAttach 领取邮件附件
func (c *Connect) GetMailAttach(act *C2SGetMailAttach) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(444, body)
}

// DelMail 删除已读邮件
func (c *Connect) DelMail(act *C2SDelMail) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(457, body)
}

func (x *S2CMailList) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][MailList] list=%v", x.MailList)
	return
}

func (x *S2CGetMailAttach) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetMailAttach] tag=%v", x.Tag)
	return
}

func (x *S2CDelMail) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][DelMail] tag=%v", x.Tag)
	return
}

func (x *S2CNewMail) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][NewMail]")
	return
}
