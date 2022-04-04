package mhyc

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type CfgErr struct {
	Id  int32
	MSG string
}

var TagMsg = make(map[int32]string)
var tagMsgMutex = &sync.Mutex{}

func init() {
	data, err := ioutil.ReadFile(DataRoot + "\\cfg_2\\Cfg_Err.json")
	if err != nil {
		panic(err)
	}
	var ret []CfgErr
	if err = json.Unmarshal(data, &ret); err != nil {
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
