package mhyc

const UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.12(0x18000c28) NetType/WIFI Language/zh_CN"

var DefineGetMailAttachAll = &C2SGetMailAttach{MailId: -1, MailType: 0}    // 邮件附件 一键领取
var DefineGiftRechargeEveryDay = &C2SActGiftNewReceive{Gid: 311, Aid: 301} // 充值->1元秒杀->每日礼
var DefineRespect = &C2SRespect{Type: 1}                                   // 排名—>跨服榜->膜拜
