package mhyc

import (
	"encoding/json"
	"sync"
)

type CfgErr struct {
	Id  int32
	MSG string
}

var TagMsg = make(map[int32]string)
var tagMsgMutex = &sync.Mutex{}

func init() {
	var ret []CfgErr
	if err := json.Unmarshal(cfg2Err, &ret); err != nil {
		panic(err)
	}
	for _, r := range ret {
		TagMsg[r.Id] = r.MSG
	}
}

func GetTagMsg(tag int32) string {
	tagMsgMutex.Lock()
	defer tagMsgMutex.Unlock()
	msg, _ := TagMsg[tag]
	return msg
}
