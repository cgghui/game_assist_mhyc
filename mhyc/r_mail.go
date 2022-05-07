package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func Mail(ctx context.Context) {
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "领取邮件附件")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			list := &S2CMailList{}
			// 普通邮件
			go func() {
				_ = CLI.MailList(DefineMailListOrdinary)
			}()
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, list, s3); err != nil {
				loop = 0
				next = RandMillisecond(6, 15)
				return
			}
			if len(list.MailList) > 0 {
				if list.MailList[0].IsReceive == 1 && list.MailList[0].IsRead == 1 {
					go func() {
						_ = CLI.GetMailAttach(DefineGetMailAttachOrdinary)
					}()
					if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetMailAttach{}, s3); err != nil {
						loop = 0
						next = RandMillisecond(6, 15)
						return
					}
				}
				count := len(list.MailList)
				i := 0
				am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
					mail := list.MailList[i]
					if mail.AttachInfo == nil && mail.AttachData == nil {
						go func() {
							_ = CLI.ReadMail(&C2SReadMail{MailId: mail.MailId, MailType: mail.MailType})
						}()
						_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CReadMail{}, s3)
					}
					i++
					if i >= count {
						return 0, 0
					}
					return ms10, 0
				})
				go func() {
					_ = CLI.DelMail(DefineDelMailOrdinary)
				}()
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CDelMail{}, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 15)
					return
				}
			}
			// 活动邮件
			go func() {
				_ = CLI.MailList(DefineMailListActivity)
			}()
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, list, s3); err != nil {
				loop = 0
				next = RandMillisecond(6, 15)
				return
			}
			if len(list.MailList) > 0 {
				if list.MailList[0].IsReceive == 1 && list.MailList[0].IsRead == 1 {
					go func() {
						_ = CLI.GetMailAttach(DefineGetMailAttachActivity)
					}()
					_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetMailAttach{}, s3)
				}
				// 忽略 阅读无附件的活动邮件
				go func() {
					_ = CLI.DelMail(DefineDelMailActivity)
				}()
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CDelMail{}, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 15)
					return
				}
			}
			return 0, RandMillisecond(600, 1800)
		})
	}
	// 监听新邮件
	go ListenMessageCall(ctx, &S2CNewMail{}, func(_ []byte) {
		go f()
	})
	// 定时打开邮件
	t := time.NewTimer(ms10)
	for {
		select {
		case <-t.C:
			t.Reset(f())
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
	log.Printf("[S][GetMailAttach] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
	return
}

////////////////////////////////////////////////////////////

func (x *S2CDelMail) ID() uint16 {
	return 458
}

func (x *S2CDelMail) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][DelMail] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
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
