package mhyc

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mozillazg/go-pinyin"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestDefine(t *testing.T) {
	t.Log(DefineGetMailAttachOrdinary)
	t.Log(DefineGetMailAttachActivity)
	t.Log(DefineGiftRechargeEveryDay)
	t.Log(DefineRespectL)
	t.Log(DefineRespectG)
}

func TestS(t *testing.T) {
	body, _ := proto.Marshal(&Ping{})
	buf := bytes.NewBuffer([]byte{})
	_ = binary.Write(buf, binary.BigEndian, int16(RandInt64(0, 65536)))
	_ = binary.Write(buf, binary.BigEndian, int16(22))
	buf.Write(body)
	ping := buf.Bytes()
	m := &sync.Mutex{}
	for {
		go func() {
			conn, _, err := websocket.DefaultDialer.Dial("wss://allws.huanlingxiuxian.com/300914:30042", nil)
			if err != nil {
				t.Error(err)
			}
			m.Lock()
			_ = conn.WriteMessage(websocket.BinaryMessage, ping)
			m.Unlock()
		}()
	}
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)
	//<-interrupt
}

func TestDeBinaryCode(t *testing.T) {
	code, err := base64.StdEncoding.DecodeString("ApJp4hJQCOwhEAkYAyI8CJ252xkSDOacrOiHquWFt+i2sxifwRIgiurqDii8tesOMLKxsbECOA5AAkhHUAZYB2ADagVTMzI2NXABKLKxsbECOL0GQCUamgYInbnbGRIHCG8Qw8imBBIFCHcQq00SBAh5EAESBAh9EA4SBQiIARABEgkIigEQzKCTkQYSBwiUARChwRISBgiVARDlARIGCMsBEMEZEgcIzwEQn8ESEgYIogIQvxkSBQilAxAJEgYI9wMQxhkSBQjYBBABEgYI3QQQ+AgSBgjeBBDXLxIGCN8EENkJEgYI4wQQ5yASBgjxBBC6FxIHCPIEEMuPRBIHCPkEEJGYRBIHCPoEEPufRBIGCPsEEOkHEgYI/AQQ6gcSCAiBBRCK6uoOEggIggUQvLXrDhIICIMFENiE7A4SBwiGBRCwkEQSBgiHBRCnChIGCIgFEKkKEgYIiQUQpAgSBQiKBRAEEgYIjQUQ2zYSBgiOBRDcNhIGCI8FEPIuEgYIkAUQ8hUSBQiRBRADEgcIkgUQxMdbEgYIlQUQxDcSBgiWBRDgNhIGCJcFEKg4EgUImAUQAhIHCJkFEJvXWxIFCJoFEAISBwicBRDN3GESBwidBRCV3mESBwieBRDd32ESBwifBRCl4WESBwigBRDt4mESBwihBRC15GESBwiiBRD95WESBwijBRDF52ESBwikBRCN6WESBwilBRDV6mESBQimBRACEgUIrAUQAhIHCK0FEKHhZxIGCL0GEIYHEgYIvgYQ4wMSBQjABhAJEgUIwgYQChIFCMQGEAoSBQjTBhAIEgkIj04QsrGxsQISBgj4VRCgAhIFCKVYEDYSBgi+WxDBGRIFCL9iEAESBQiPaxAGEgYImWsQxj4SBgibaxCmPxIGCJ1rEIpAEgUIzGwQRxIFCM1sEAYSBQjObBAHEgUIz2wQAxIFCMF1EAEaEAhxEgzmnKzoh6rlhbfotrMaCgikAhIFUzMyNjUaDAilAhIHUzMwMzI2MxoSCKYCEg3lvqHliZEzMjY15Yy6GhMI9AMSDjU0MTYwNTg0XzE4NzkwGggI9gMSA+m+mRoUCLUFEg81MzkyNzA2OV8xNDA5ODYaEQjxVhIM5p+g5qqs5p6c5Ya7GhEIv1sSDOS6uueUn+aXoOW4uBoKCMFbEgVTMzI2NQ==") //EWFnIQiluQI=
	if err != nil {
		t.Error(err)
	}
	var ma = &S2CCreateTeam{}
	if err = proto.Unmarshal(code[4:], ma); err != nil {
		t.Error(err)
	}
	// type:8 id:385
	t.Log(ma) // 13936 8
}

func TestDeBinaryCode2(t *testing.T) {
	code, err := base64.StdEncoding.DecodeString("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") //EWFnIQiluQI=
	if err != nil {
		t.Error(err)
	}
	var ma = &S2CTeamInfo{}
	if err = proto.Unmarshal(code[4:], ma); err != nil {
		t.Error(err)
	}
	// type:8 id:385
	t.Log(ma) // 13936 8
}

func TestJson2PB(t *testing.T) {
	s, _ := ioutil.ReadFile("D:\\go\\game_assist_mhyc\\mhyc\\data.json")
	var ret map[string]map[string]map[string]interface{}
	err := json.Unmarshal(s, &ret)
	if err != nil {
		t.Error(err)
	}
	var rr = make([]string, 0, 0)
	for k, v := range ret {
		var ss = make([]string, 0, 0)
		for fn, vx := range v["fields"] {
			vv := vx.(map[string]interface{})
			var rule string
			if _, ok := vv["rule"]; ok {
				rule = vv["rule"].(string) + " "
			}
			ss = append(ss, fmt.Sprintf("\t%s%v %v = %v;", rule, vv["type"], fn, vv["id"]))
		}
		rr = append(rr, "message "+k+"{\n"+strings.Join(ss, "\n")+"\n}")
	}
	err = ioutil.WriteFile("D:\\go\\game_assist_mhyc\\mhyc\\sss.proto", []byte(strings.Join(rr, "\n")), 0666)
	fmt.Println(ret)
}

// 11 250
func TestMPQ(t *testing.T) {
	re, _ := LoadDataRes("D:\\go\\game_assist_mhyc\\mhyc\\cfg_3_data.txt")
	for ii, i := range re.Index {
		fmt.Println(ii, " ", i.Name)
		_ = re.GetData(&re.Index[ii])
		ioutil.WriteFile("D:\\go\\game_assist_mhyc\\mhyc\\cfg_3\\"+i.Name, re.Index[ii].Data, 0666)
	}
	fmt.Println()
}

func TestC2P(t *testing.T) {
	f1 := "C:\\Users\\admin\\Desktop\\新建文本文档 (2).txt"
	f2 := "C:\\Users\\admin\\Desktop\\新建文本文档 (3).txt"
	content, err := ioutil.ReadFile(f1)
	if err != nil {
		t.Error(err)
	}
	var w *os.File
	if w, err = os.Create(f2); err != nil {
		t.Error(err)
	}
	defer func() { _ = w.Close() }()
	arg := pinyin.NewArgs()
	content = bytes.ReplaceAll(content, []byte("\r"), []byte{})
	for _, line := range bytes.Split(content, []byte("\n")) {
		str := string(line)
		py := strings.Join(pinyin.LazyPinyin(str, arg), "")
		if py == "" {
			py = str
		}
		_, _ = w.WriteString(strings.ToLower(py) + "\n")
	}
}
