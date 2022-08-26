package main

import (
	"biliDanMu/callback"
	"biliDanMu/handler"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

func init() {
	flag.Parse()
}

func main() {
	var roomid uint32

	wg := &sync.WaitGroup{}
	wg.Add(1)

	if callback.Conf.RoomID == 0 {
		fmt.Print("请输入房间号，长短 ID 均可：")
		_, err := fmt.Scanf("%d", &roomid)
		if err != nil {
			log.Println("房间号输入错误，请退出重新输入！")
			return
		}

		input := bufio.NewScanner(os.Stdin)
		for {
			input.Scan()
		}
	} else {
		if callback.Conf.PostUrl == "" || callback.Conf.NewDmUrl == "" || callback.Conf.NewGiftUrl == "" {
			log.Println("参数校验错误")
			return
		}

		roomid = uint32(callback.Conf.RoomID)
	}

	// 兼容房间号短 ID
	if roomid >= 100 && roomid < 1000 {
		r, err := handler.GetRealRoomID(int(roomid))
		if err != nil {
			log.Println("房间号输入错误，请退出重新输入！")
			return
		}
		roomid = r
	}

	c, err := handler.NewClient(roomid)
	if err != nil {
		log.Println("models.NewClient err: ", err)
		return
	}
	if err = c.Start(); err != nil {
		log.Println("c.Start err :", err)
		return
	}

	go func() {
		escape := ""
		_, _ = fmt.Scanf("%v\n", &escape)
		if escape == "exit" {
			wg.Done()
		}
	}()

	wg.Wait()
}
