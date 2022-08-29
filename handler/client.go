package handler

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"

	"biliDanMu/callback"
	"biliDanMu/models"
)

// Client instance
type Client struct {
	Room      RoomInfo     `json:"room"`
	Request   *RequestInfo `json:"request"`
	conn      *websocket.Conn
	Connected bool `json:"connected"`
}

// RoomInfo basic information of the live room
type RoomInfo struct {
	RoomId     uint32 `json:"room_id"`
	UpUid      uint32 `json:"up_uid"`
	Title      string `json:"title"`
	Online     uint32 `json:"online"`
	Tags       string `json:"tags"`
	LiveStatus bool   `json:"live_status"`
	LockStatus bool   `json:"lock_status"`
}

// RequestInfo data on handshake packets
type RequestInfo struct {
	Uid       uint8  `json:"uid"`
	Roomid    uint32 `json:"roomid"`
	Protover  uint8  `json:"protover"`
	Platform  string `json:"platform"`
	Clientver string `json:"clientver"`
	Type      uint8  `json:"type"`
	Key       string `json:"key"`
}

// NewRequestInfo return initialized structure
func NewRequestInfo(roomid uint32) *RequestInfo {
	t := GetToken(roomid)
	return &RequestInfo{
		Uid:       0,
		Roomid:    roomid,
		Protover:  2,
		Platform:  "web",
		Clientver: "1.10.2",
		Type:      2,
		Key:       t,
	}
}

// NewClient return a new websocket client
func NewClient(roomid uint32) (c *Client, err error) {
	return &Client{
		Room:      GetRoomInfo(roomid),
		Request:   NewRequestInfo(roomid),
		conn:      nil,
		Connected: false,
	}, nil
}

func (c *Client) Start() (err error) {
	u := url.URL{Scheme: "wss", Host: models.DanMuServer, Path: "/sub"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		fmt.Println("net.Dial err: ", err)
		return err
	}
	c.conn = conn

	fmt.Println("当前直播间状态：", c.Room.LiveStatus)

	fmt.Println("连接弹幕服务器 ", models.DanMuServer, " 成功，正在发送握手包...")
	r, err := json.Marshal(c.Request)

	if err != nil {
		fmt.Println("marshal err ,", err)
		return
	}
	if err = c.SendPackage(0, 16, 1, 7, 1, r); err != nil {
		fmt.Println("SendPackage err,", err)
		return
	}
	go c.HeartBeat()
	c.ReceiveMsg()
	return
}

func (c *Client) SendPackage(packetlen uint32, magic uint16, ver uint16, typeID uint32, param uint32, data []byte) (err error) {
	packetHead := new(bytes.Buffer)

	if packetlen == 0 {
		packetlen = uint32(len(data) + 16)
	}
	var pdata = []interface{}{
		packetlen,
		magic,
		ver,
		typeID,
		param,
	}

	// 将包的头部信息以大端序方式写入字节数组
	for _, v := range pdata {
		if err = binary.Write(packetHead, binary.BigEndian, v); err != nil {
			fmt.Println("binary.Write err: ", err)
			return
		}
	}

	// 将包内数据部分追加到数据包内
	sendData := append(packetHead.Bytes(), data...)

	// fmt.Println("本次发包消息为：", sendData)

	if err = c.conn.WriteMessage(websocket.BinaryMessage, sendData); err != nil {
		fmt.Println("c.conn.Write err: ", err)
		return
	}

	return
}

func (c *Client) ReceiveMsg() {
	pool := NewPool()
	if callback.Conf.MiddlewareMode {
		go pool.HandleWithCallback()
	} else {
		go pool.Handle()
	}

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("ReadMsg err :", err)
			continue
		}

		switch msg[11] {
		case 8:
			fmt.Println("握手包收发完毕，连接成功")
			c.Connected = true
		case 3:
			onlineNow := ByteArrToDecimal(msg[16:])
			if uint32(onlineNow) != c.Room.Online {
				c.Room.Online = uint32(onlineNow)
				fmt.Println("当前房间人气变动：", uint32(onlineNow))
			}
		case 5:
			if inflated, err := ZlibInflate(msg[16:]); err != nil {
				// 代表是未压缩数据
				pool.MsgUncompressed <- msg[16:]
			} else {
				for len(inflated) > 0 {
					l := ByteArrToDecimal(inflated[:4])
					c := models.Json.Get(inflated[16:l], "cmd").ToString()
					switch models.CMD(c) {
					case models.CMDDanmuMsg:
						pool.UserMsg <- inflated[16:l]
					case models.CMDSendGift:
						pool.UserGift <- inflated[16:l]
					case models.CMDWELCOME:
						pool.UserEnter <- inflated[16:l]
					case models.CMDWelcomeGuard:
						pool.UserGuard <- inflated[16:l]
					case models.CMDEntry:
						pool.UserEntry <- inflated[16:l]
					}
					inflated = inflated[l:]
				}
			}
		}
	}
}

func (c *Client) HeartBeat() {
	for {
		if c.Connected {
			obj := []byte("5b6f626a656374204f626a6563745d")
			if err := c.SendPackage(31, 16, 1, 2, 1, obj); err != nil {
				log.Println("heart beat err: ", err)
				continue
			}
			time.Sleep(30 * time.Second)
		}
	}
}
