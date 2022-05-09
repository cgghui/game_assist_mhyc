package mhyc

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	channelID = "7108"
	serverID  = "300908"
	areaID    = "303263"
	areaID64  = int64(303263)
	UUID      = "h51645175082697"
)

type GameSession struct {
	sessionID string
	channelID string
}

// BaseUserinfo 包含微信用户信息
type BaseUserinfo struct {
	WeChatName string `json:"wechaname"`
	UnionID    string `json:"unionid"`
	Sex        string `json:"sex"`
	RealOpenID string `json:"real_openid"`
	Province   string `json:"province"`
	Portrait   string `json:"portrait"`
	OpenID     string `json:"openid"`
	IsJumper   int    `json:"is_jumper"`
	City       string `json:"city"`
}

func NewClient(sid, cid string) (*Client, error) {
	s := &GameSession{sessionID: sid, channelID: cid}
	sign, err := s.sign()
	if err != nil {
		return nil, err
	}
	var u *BaseUserinfo
	if u, err = s.getUserinfo(sign); err != nil {
		return nil, err
	}
	return s.login(u)
}

func (s *GameSession) login(i *BaseUserinfo) (*Client, error) {
	var req *http.Request
	var err error
	body := "sdk=shengye&app_id=8&channel_id=" + channelID + "&uid=" + i.OpenID + "&token=" + s.sessionID + "&ext=&imei=" + UUID + "&phone_model=&sys_ver=&platform=h5&system=%E6%9C%AA%E7%9F%A5"
	req, err = http.NewRequest(http.MethodPost, "https://sdk.tianzongyouxi.com/v1/sdk/login", strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	type Body struct {
		Msg    string `json:"msg"`
		Status int    `json:"status"`
		SDK    string `json:"sdk"`
		Data   Client `json:"data"`
	}
	var ret Body
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}
	if ret.Msg == "请求成功" && ret.Status == 200 {
		return &ret.Data, nil
	}
	return nil, errors.New("login error: " + ret.Msg)
}

func (s *GameSession) getUserinfo(sign string) (*BaseUserinfo, error) {
	var req *http.Request
	var err error
	req, err = http.NewRequest(http.MethodGet, "https://docater1.cn/index.php?g=Home&m=GameOauth&a=get_userinfo&channel_id="+s.channelID+"&userToken="+s.sessionID+"&sign="+sign, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UserAgent)
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	type Body struct {
		Info     string       `json:"info"`
		Status   int          `json:"status"`
		UserInfo BaseUserinfo `json:"userinfo"`
	}
	var ret Body
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}
	if ret.Info == "success" && ret.Status == 1001 {
		return &ret.UserInfo, nil
	}
	return nil, errors.New("getUserinfo error: " + ret.Info)
}

// sign 获取游戏的根本信息
func (s *GameSession) sign() (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://sdk.tianzongyouxi.com/v1/sdk/ext/shengye/sign/8/"+channelID+"?token="+s.sessionID, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", UserAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	type Body struct {
		Data   string `json:"data"`
		Status int    `json:"status"`
	}
	var ret Body
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return "", err
	}
	if ret.Status != http.StatusOK {
		return "", errors.New("sign error: " + strconv.Itoa(ret.Status))
	}
	return ret.Data, nil
}
