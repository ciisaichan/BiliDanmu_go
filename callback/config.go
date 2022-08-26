package callback

import "flag"

var Conf = &Config{}

type Config struct {
	RoomID         uint
	MiddlewareMode bool
	Token          string
	PostUrl        string
	NewDmUrl       string
	NewGiftUrl     string
}

func init() {
	flag.UintVar(&Conf.RoomID, "RoomID", 0, "Live Room ID")
	flag.BoolVar(&Conf.MiddlewareMode, "MiddlewareMode", false, "Enable Middleware Mode")
	flag.StringVar(&Conf.Token, "Token", "", "Call-back Server Access Token")
	flag.StringVar(&Conf.PostUrl, "PostURL", "", "Server Base URL for Call-back")
	flag.StringVar(&Conf.NewDmUrl, "DanMuRouterPath", "danmu", "Server DanMu Router URL for Call-back")
	flag.StringVar(&Conf.NewGiftUrl, "GiftRouterPath", "gift", "Server Gift Router URL for Call-back")
}
