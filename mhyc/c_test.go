package mhyc

import (
	"bytes"
	"encoding/base64"
	"github.com/mozillazg/go-pinyin"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDefine(t *testing.T) {
	t.Log(DefineGetMailAttachAll)
	t.Log(DefineGiftRechargeEveryDay)
	t.Log(DefineRespect)
}

func TestDeBinaryCode(t *testing.T) {
	code, err := base64.StdEncoding.DecodeString("AY8ADQgB")
	if err != nil {
		t.Error(err)
	}
	var ma C2SRespect
	if err = proto.Unmarshal(code[4:], &ma); err != nil {
		t.Error(err)
	}
	t.Log(ma)
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
