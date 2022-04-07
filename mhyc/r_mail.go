package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func Mail(ctx context.Context) {
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
			for _, mail := range list.MailList {
				if mail.AttachInfo == nil && mail.AttachData == nil {
					go func() {
						_ = CLI.ReadMail(&C2SReadMail{MailId: mail.MailId, MailType: mail.MailType})
					}()
					_ = Receive.Wait(&S2CReadMail{}, s3)
				}
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
	go ListenMessageCall(ctx, &S2CNewMail{}, func(_ []byte) {
		f()
	})
	// 定时打开邮件
	t := time.NewTimer(ms10)
	for {
		select {
		case <-t.C:
			f()
			t.Reset(RandMillisecond(300, 600)) // 5 ~ 10 分钟
		case <-ctx.Done():
			return
		}
	}
}

// MailList 邮件列表
func (c *Connect) MailList(act *C2SMailList) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][MailList] mail_type=%v", act.MailType)
	return c.send(440, body)
}

// GetMailAttach 领取邮件附件
func (c *Connect) GetMailAttach(act *C2SGetMailAttach) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][GetMailAttach] mail_id=%v mail_type=%v", act.MailId, act.MailType)
	return c.send(444, body)
}

// ReadMail 读邮件
func (c *Connect) ReadMail(act *C2SReadMail) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][ReadMail] mail_id=%v mail_type=%v", act.MailId, act.MailType)
	return c.send(448, body)
}

// DelMail 删除已读邮件
func (c *Connect) DelMail(act *C2SDelMail) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][DelMail] mail_id=%v mail_type=%v", act.MailId, act.MailType)
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

////////////////////////////////////////////////////////////

func (x *S2CReadMail) ID() uint16 {
	return 449
}

func (x *S2CReadMail) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ReadMail]")
	return
}
