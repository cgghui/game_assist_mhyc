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
		_ = Receive.Wait(list, s3)
		if len(list.MailList) > 0 {
			if list.MailList[0].IsReceive == 1 && list.MailList[0].IsRead == 1 {
				go func() {
					_ = CLI.GetMailAttach(DefineGetMailAttachOrdinary)
				}()
				_ = Receive.Wait(&S2CGetMailAttach{}, s3)
			}
			go func() {
				_ = CLI.DelMail(DefineDelMailOrdinary)
			}()
			_ = Receive.Wait(&S2CDelMail{}, s3)
		}
		// 活动邮件
		go func() {
			_ = CLI.MailList(DefineMailListActivity)
		}()
		_ = Receive.Wait(list, s3)
		if len(list.MailList) > 0 {
			if list.MailList[0].IsReceive == 1 && list.MailList[0].IsRead == 1 {
				go func() {
					_ = CLI.GetMailAttach(DefineGetMailAttachActivity)
				}()
				_ = Receive.Wait(&S2CGetMailAttach{}, s3)
			}
			go func() {
				_ = CLI.DelMail(DefineDelMailActivity)
			}()
			_ = Receive.Wait(&S2CDelMail{}, s3)
		}
	}
	// 监听新邮件
	go func() {
		channel := Receive.CreateChannel(&S2CNewMail{})
		for range channel.Wait() {
			f()
		}
	}()
	// 定时打开邮件
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

////////////////////////////////////////////////////////////

func (x *S2CMailList) ID() uint16 {
	return 441
}

func (x *S2CMailList) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][MailList] list=%v", x.MailList)
	return
}

////////////////////////////////////////////////////////////

func (x *S2CGetMailAttach) ID() uint16 {
	return 445
}

func (x *S2CGetMailAttach) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetMailAttach] tag=%v", x.Tag)
	return
}

////////////////////////////////////////////////////////////

func (x *S2CDelMail) ID() uint16 {
	return 458
}

func (x *S2CDelMail) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][DelMail] tag=%v", x.Tag)
	return
}

////////////////////////////////////////////////////////////

func (x *S2CNewMail) ID() uint16 {
	return 446
}

func (x *S2CNewMail) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][NewMail]")
	return
}
