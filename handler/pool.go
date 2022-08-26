package handler

import (
	"fmt"
	"log"
	"strings"

	"biliDanMu/callback"
	"biliDanMu/models"
)

// Pool`s fields map CMD value
type Pool struct {
	callbackClient  *callback.Client
	UserMsg         chan []byte
	UserGift        chan []byte
	UserEnter       chan []byte
	UserGuard       chan []byte
	MsgUncompressed chan []byte
	UserEntry       chan []byte
}

func NewPool() *Pool {
	return &Pool{
		callbackClient:  callback.NewClient(),
		UserMsg:         make(chan []byte, 10),
		UserGift:        make(chan []byte, 10),
		UserEnter:       make(chan []byte, 10),
		MsgUncompressed: make(chan []byte, 10),
		UserEntry:       make(chan []byte, 10),
		UserGuard:       make(chan []byte, 10),
	}
}

func (pool *Pool) Handle() {
	for {
		select {
		case uc := <-pool.MsgUncompressed:
			// 目前只处理未压缩数据的关注数变化信息
			if cmd := models.Json.Get(uc, "cmd").ToString(); models.CMD(cmd) == models.CMDRoomRealTimeMessageUpdate {
				fans := models.Json.Get(uc, "data", "fans").ToInt()
				fmt.Println("当前房间关注数变动：", fans)
			}
		case src := <-pool.UserMsg:
			m := models.NewDanmu()
			m.GetDanmuMsg(src)
			fmt.Printf("%d-%s | %d-%s: %s\n", m.MedalLevel, m.MedalName, m.Ulevel, m.Uname, m.Text)
		case src := <-pool.UserGift:
			g := models.NewGift()
			g.GetGiftMsg(src)
			fmt.Printf("%s %s 价值 %d 的 %s\n", g.UUname, g.Action, g.Price, g.GiftName)
		case src := <-pool.UserEnter:
			name := models.Json.Get(src, "data", "uname").ToString()
			fmt.Printf("欢迎VIP %s 进入直播间\n", name)
		case src := <-pool.UserGuard:
			name := models.Json.Get(src, "data", "username").ToString()
			fmt.Printf("欢迎房管 %s 进入直播间\n", name)
		case src := <-pool.UserEntry:
			cw := models.Json.Get(src, "data", "copy_writing").ToString()
			cw = strings.Replace(cw, "<%", "", 1)
			cw = strings.Replace(cw, "%>", "", 1)
			fmt.Printf("%s\n", cw)
		}
	}
}

func (pool *Pool) HandleWithCallback() {
	for {
		select {
		case uc := <-pool.MsgUncompressed:
			// 目前只处理未压缩数据的关注数变化信息
			if cmd := models.Json.Get(uc, "cmd").ToString(); models.CMD(cmd) == models.CMDRoomRealTimeMessageUpdate {
				fans := models.Json.Get(uc, "data", "fans").ToInt()
				fmt.Println("当前房间关注数变动：", fans)
			}
		case src := <-pool.UserMsg:
			m := models.NewDanmu()
			m.GetDanmuMsg(src)
			ctx, err := callback.NewContextWithDanMu(m)
			if err != nil {
				log.Println("装填弹幕数据错误", err)
			}
			pool.callbackClient.CallBack(ctx)
			fmt.Printf("%d-%s | %d-%s: %s\n", m.MedalLevel, m.MedalName, m.Ulevel, m.Uname, m.Text)
		case src := <-pool.UserGift:
			g := models.NewGift()
			g.GetGiftMsg(src)
			ctx, err := callback.NewContextWithGift(g)
			if err != nil {
				log.Println("装填礼物数据错误", err)
			}
			pool.callbackClient.CallBack(ctx)
			fmt.Printf("%s %s 价值 %d 的 %s\n", g.UUname, g.Action, g.Price, g.GiftName)
		case src := <-pool.UserEnter:
			name := models.Json.Get(src, "data", "uname").ToString()
			fmt.Printf("欢迎VIP %s 进入直播间\n", name)
		case src := <-pool.UserGuard:
			name := models.Json.Get(src, "data", "username").ToString()
			fmt.Printf("欢迎房管 %s 进入直播间\n", name)
		case src := <-pool.UserEntry:
			cw := models.Json.Get(src, "data", "copy_writing").ToString()
			cw = strings.Replace(cw, "<%", "", 1)
			cw = strings.Replace(cw, "%>", "", 1)
			fmt.Printf("%s\n", cw)
		}
	}
}
