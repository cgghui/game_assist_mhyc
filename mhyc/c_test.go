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
	code, err := base64.StdEncoding.DecodeString("BFdYLwgF")
	if err != nil {
		t.Error(err)
	}
	var ma C2SClimbingTowerFight
	if err = proto.Unmarshal(code[4:], &ma); err != nil {
		t.Error(err)
	}
	t.Log(ma)
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
