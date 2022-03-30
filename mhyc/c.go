package mhyc

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
}

type Connect struct {
	Code      int             `json:"code"`
	AccountID int64           `json:"account_id"`
	UserID    int64           `json:"user_id"`
	Sign      string          `json:"sign"`
	Timestamp int64           `json:"timestamp"`
	IP        string          `json:"ip"`
	Msg       string          `json:"msg"`
	Conn      *websocket.Conn `json:"-"`
	m         *sync.Mutex
}

func (c *Client) Connect(ctx context.Context) (*Connect, error) {
	param := url.Values{}
	param.Add("channel_id", channelID)
	param.Add("token", c.Token)
	param.Add("server_id", serverID)
	param.Add("area_id", areaID)
	param.Add("user_id", strconv.Itoa(c.UserID))
	param.Add("uuid", UUID)
	param.Add("sys_ver", "")
	param.Add("phone_model", "")
	param.Add("auto_create", "1")
	req, err := http.NewRequest(http.MethodGet, "https://cdns1.huanlingxiuxian.com/tz/login?"+param.Encode(), nil)
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
	var ret Connect
	if err = json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err
	}
	arr := strings.Split(ret.IP, "|")
	ret.Conn, _, err = websocket.DefaultDialer.Dial("wss://"+arr[len(arr)-1], nil)
	if err != nil {
		return nil, err
	}
	ret.m = &sync.Mutex{}
	err = ret.login(&C2SLogin{
		AreaId:       areaID64,
		AccountId:    ret.AccountID,
		Token:        ret.Sign,
		UserId:       ret.UserID,
		Fcm:          2,
		LoginPf:      "h5",
		CheckWordUrl: "",
		CodeVersion:  30670,
		ExcelVersion: 30613,
	})
	if err != nil {
		return nil, err
	}
	go func() {
		t := time.NewTicker(3 * time.Second)
		for {
			select {
			case <-t.C:
				_ = ret.Ping()
				log.Printf("[C][Ping]")
			case <-ctx.Done():
				_ = ret.Close()
				return
			}
		}
	}()
	return &ret, nil
}

func (c *Connect) Close() error {
	return c.Conn.Close()
}

func (c *Connect) login(info *C2SLogin) error {
	body, err := proto.Marshal(info)
	if err != nil {
		return err
	}
	return c.send(1, body)
}

func (c *Connect) Ping() error {
	body, err := proto.Marshal(&Ping{})
	if err != nil {
		return err
	}
	return c.send(22, body)
}

func (c *Connect) GetBattlePrize() error {
	body, err := proto.Marshal(&C2SGetBattlePrize{})
	if err != nil {
		return err
	}
	return c.send(523, body)
}

func (c *Connect) FightContinue() error {
	body, err := proto.Marshal(&C2SFightContinue{})
	if err != nil {
		return err
	}
	return c.send(135, body)
}

func (c *Connect) CheckFight() error {
	body, err := proto.Marshal(&C2SCheckFight{})
	if err != nil {
		return err
	}
	return c.send(65, body)
}

// RoutePath 主页 走来走去
func (c *Connect) RoutePath() error {
	body, err := proto.Marshal(&C2SRoutePath{
		MapId: 33,
		FX:    int32(RandInt64(1, 255)),
		FY:    int32(RandInt64(1, 255)),
		TX:    int32(RandInt64(1, 255)),
		TY:    int32(RandInt64(1, 255)),
	})
	if err != nil {
		return err
	}
	return c.send(154, body)
}

func (c *Connect) ChangeMap(m *C2SChangeMap) error {
	body, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return c.send(50, body)
}

func (c *Connect) HuanLingList() error {
	body, err := proto.Marshal(&C2SHuanLingList{})
	if err != nil {
		return err
	}
	return c.send(27151, body)
}

func (c *Connect) GetActTimestamp(act *C2SGetActTimestamp) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(1541, body)
}

// ShopBuy 商城购物
func (c *Connect) ShopBuy(goods *C2SShopBuy) error {
	body, err := proto.Marshal(goods)
	if err != nil {
		return err
	}
	return c.send(432, body)
}

// RealmTask 修仙 - 境界 任务
func (c *Connect) RealmTask() error {
	body, err := proto.Marshal(&C2SRealmTask{})
	if err != nil {
		return err
	}
	return c.send(22012, body)
}

func (c *Connect) GetPetAMergeInfo() error {
	body, err := proto.Marshal(&C2SGetPetAMergeInfo{})
	if err != nil {
		return err
	}
	return c.send(22730, body)
}

func (c *Connect) GetAllEquipData() error {
	body, err := proto.Marshal(&C2SGetAllEquipData{})
	if err != nil {
		return err
	}
	return c.send(27001, body)
}

func (c *Connect) PlayerPractice() error {
	body, err := proto.Marshal(&C2SPlayerPractice{})
	if err != nil {
		return err
	}
	return c.send(23101, body)
}

func (c *Connect) Beasts() error {
	body, err := proto.Marshal(&C2SBeasts{})
	if err != nil {
		return err
	}
	return c.send(25795, body)
}

func (c *Connect) GetEquipData(act *C2SGetEquipData) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(27003, body)
}

func (c *Connect) GetHeroList() error {
	body, err := proto.Marshal(&C2SGetHeroList{})
	if err != nil {
		return err
	}
	return c.send(27801, body)
}

func (c *Connect) GetAlienData() error {
	body, err := proto.Marshal(&C2SGetAlienData{})
	if err != nil {
		return err
	}
	return c.send(28601, body)
}

func (c *Connect) YJInfo() error {
	// 偃钾
	body, err := proto.Marshal(&C2SYJInfo{})
	if err != nil {
		return err
	}
	return c.send(52226, body)
}

func (c *Connect) SLGetData() error {
	body, err := proto.Marshal(&C2SSLGetData{})
	if err != nil {
		return err
	}
	return c.send(29503, body)
}

func (c *Connect) NewStory() error {
	body, err := proto.Marshal(&C2SNewStory{})
	if err != nil {
		return err
	}
	return c.send(36, body)
}

// LoginEnd ?
func (c *Connect) LoginEnd() error {
	body, err := proto.Marshal(&C2SLoginEnd{})
	if err != nil {
		return err
	}
	return c.send(29, body)
}

func (c *Connect) GetActXunBaoInfo(act *C2SGetActXunBaoInfo) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(11031, body)
}

// AllWeddingToken 仙缘信息
func (c *Connect) AllWeddingToken() error {
	body, err := proto.Marshal(&C2SAllWeddingToken{})
	if err != nil {
		return err
	}
	return c.send(22627, body)
}

func (c *Connect) WeddingInsInvite() error {
	body, err := proto.Marshal(&C2SWeddingInsInvite{})
	if err != nil {
		return err
	}
	return c.send(22627, body)
}

// WeddingInsFight 仙缘副本 - 战斗
func (c *Connect) WeddingInsFight() error {
	body, err := proto.Marshal(&C2SWeddingInsFight{})
	if err != nil {
		return err
	}
	return c.send(22629, body)
}

// ClimbingTowerEnter 副本 - 爬塔 - 进入
func (c *Connect) ClimbingTowerEnter(act *C2SClimbingTowerEnter) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(22571, body)
}

// ClimbingTowerFight 副本 - 爬塔 - 战斗
func (c *Connect) ClimbingTowerFight(act *C2SClimbingTowerFight) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(22575, body)
}

func (c *Connect) send(code int, body []byte) error {
	var err error
	idx := uint16(RandInt64(0, 65536))
	buf := bytes.NewBuffer([]byte{})
	if err = binary.Write(buf, binary.BigEndian, idx); err != nil {
		return err
	}
	if err = binary.Write(buf, binary.BigEndian, uint16(code)); err != nil {
		return err
	}
	buf.Write(body)
	c.m.Lock()
	defer c.m.Unlock()
	return c.Conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
}
