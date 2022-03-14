package mhyc

import (
	"bytes"
	"encoding/binary"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"net/http"
	"net/url"
)

func NewClient() (*Client, error) {
	param := url.Values{}
	param.Add("channel_id", "7108")
	param.Add("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDc4NDE2MjEsImlzcyI6IjE3MDM3MzI1MSJ9.c0KRjlajdbjHKiJvwwjFx6RcuCTeKjY-aAZ9UG9nVXs")
	param.Add("server_id", "300914")
	param.Add("area_id", "303265")
	param.Add("user_id", "170373251")
	param.Add("uuid", "h51645175082697")
	param.Add("sys_ver", "")
	param.Add("phone_model", "")
	param.Add("role_id", "53927069")
	param.Add("auto_create", "1")
	req, err := http.NewRequest(http.MethodGet, "https://cdns1.huanlingxiuxian.com/tz/login?"+param.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", UserAgent)
	return nil, nil
}

type Client struct {
	Conn *websocket.Conn
}

func (c *Client) Login(info *C2SLogin) error {
	body, err := proto.Marshal(info)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err = binary.Write(buf, binary.BigEndian, int32(1)); err != nil {
		return err
	}
	buf.Write(body)
	return c.Conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
}

func (c *Client) Ping() error {
	body, err := proto.Marshal(&Ping{})
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err = binary.Write(buf, binary.BigEndian, int32(22)); err != nil {
		return err
	}
	buf.Write(body)
	return c.Conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
}

func (c *Client) MailList() error {
	body, err := proto.Marshal(&Ping{})
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err = binary.Write(buf, binary.BigEndian, int32(440)); err != nil {
		return err
	}
	buf.Write(body)
	return c.Conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
}

func (c *Client) GetMailAttach(act *C2SGetMailAttach) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err = binary.Write(buf, binary.BigEndian, int32(444)); err != nil {
		return err
	}
	buf.Write(body)
	return c.Conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
}

func (c *Client) ActGiftNewReceive(act *C2SActGiftNewReceive) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err = binary.Write(buf, binary.BigEndian, int32(12011)); err != nil {
		return err
	}
	buf.Write(body)
	return c.Conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
}

func (c *Client) Respect(act *C2SRespect) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer([]byte{})
	if err = binary.Write(buf, binary.BigEndian, int32(13)); err != nil {
		return err
	}
	buf.Write(body)
	return c.Conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())
}
