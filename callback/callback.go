package callback

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"biliDanMu/models"
)

type Client struct {
	client *http.Client
}

func (c *Client) CallBack(ctx *Context) {
	// 拼接回调地址
	u, err := url.JoinPath(Conf.PostUrl, ctx.router)
	if err != nil {
		log.Println("回调请求装填错误", err)
	}

	// 构建回调数据
	request, err := http.NewRequest(http.MethodPost, u, ctx.payload)
	if err != nil {
		log.Println("回调数据装填错误", err)
	}
	request.Header.Add("Authorization", "Bearer "+Conf.Token)
	request.Header.Add("Content-Type", "application/json")

	// 执行回调请求
	_, err = c.client.Do(request)
	if err != nil {
		log.Println("回调请求发送错误", err)
	}
}

func NewClient() *Client {
	return &Client{client: &http.Client{}}
}

type Context struct {
	payload io.Reader
	router  string
}

func NewContextWithGift(m *models.Gift) (ctx *Context, err error) {
	payload, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return &Context{payload: bytes.NewReader(payload), router: Conf.NewGiftUrl}, nil
}

func NewContextWithDanMu(m *models.DanMuMsg) (*Context, error) {
	payload, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return &Context{payload: bytes.NewReader(payload), router: Conf.NewDmUrl}, nil
}
