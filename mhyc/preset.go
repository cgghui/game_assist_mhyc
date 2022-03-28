package mhyc

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

const DataRoot = "D:\\go\\game_assist_mhyc\\mhyc"

const UserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.12(0x18000c28) NetType/WIFI Language/zh_CN"

var DefineMailListOrdinary = &C2SMailList{MailType: 0} // 普通邮件列表
var DefineMailListActivity = &C2SMailList{MailType: 1} // 活动邮件列表

var DefineGetMailAttachOrdinary = &C2SGetMailAttach{MailId: -1, MailType: 0} // 普通邮件附件 一键领取
var DefineGetMailAttachActivity = &C2SGetMailAttach{MailId: -1, MailType: 1} // 活动邮件附件 一键领取

var DefineDelMailOrdinary = &C2SDelMail{MailId: 0, MailType: 0} // 普通邮件 删除已读
var DefineDelMailActivity = &C2SDelMail{MailId: 0, MailType: 1} // 活动邮件 删除已读

var DefineGiftRechargeEveryDay = &C2SActGiftNewReceive{Gid: 311, Aid: 301} // 充值->1元秒杀->每日礼

var DefineRespectL = &C2SRespect{Type: 0} // 排名—>本区榜->膜拜
var DefineRespectG = &C2SRespect{Type: 1} // 排名—>跨服榜->膜拜

var DefineVipDayGift = &C2SGetVipDayGift{} // SVIP 每日礼包

var DefineGetActTask11002 = &C2SGetActTask{ActId: 11002} // 每日福利 -> 在线奖励

var DefineLifeCardDayPrize = &C2SLifeCardDayPrize{} // 特权卡 -> 至尊卡

var DefineSign = &C2SSign{} // 每日签到

var DefineStageFight = &C2SStageFight{} // 闯关

var DefineGetHistoryTaskPrize = &C2SGetHistoryTaskPrize{TaskId: 0} // 完成主线任务 - 领取奖励

var DefineStageDraw = &C2SStageDraw{} // 幸运转盘

var DefineShopBuyFree = &C2SShopBuy{GoodsId: 11001, Num: 1} // 商城 - 每日免费领的商品

var SearchPet502 = &C2SSearchPet{ItemId: 502} // 寻找宠物 - 神兽号角
var SearchPet501 = &C2SSearchPet{ItemId: 501} // 寻找宠物 - 高级
var SearchPet500 = &C2SSearchPet{ItemId: 500} // 寻找宠物 - 初级

var DefineClimbingTowerEnter5 = &C2SClimbingTowerEnter{TowerType: 5} // 副本 - 爬塔 - 剑魂塔 - 进入
var DefineClimbingTowerFight5 = &C2SClimbingTowerFight{TowerType: 5} // 副本 - 爬塔 - 剑魂塔 - 战斗

var DefineChangeMap33 = &C2SChangeMap{MapId: 33} // 切换地图 主图

func RandInt64(min, max int64) int64 {
	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))
		return result.Int64() - i64Min
	}
	result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	return min + result.Int64()
}

func RandMillisecond(min, max int64) time.Duration {
	min *= 1000
	max *= 1000
	return time.Duration(RandInt64(min, max)) * time.Millisecond
}

// SelfWeekMonday 指定日期的周一
func SelfWeekMonday(tm time.Time) time.Time {
	offset := int(time.Monday - tm.Weekday())
	if offset > 0 {
		offset = -6
	}
	return time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
}

// Tomorrow 指定日期的明天
func Tomorrow(tm time.Time) time.Time {
	return time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
}

// TomorrowDuration 指定日期的明天
func TomorrowDuration(d time.Duration) time.Duration {
	now := time.Now()
	return Tomorrow(now).Add(d).Sub(now)
}
