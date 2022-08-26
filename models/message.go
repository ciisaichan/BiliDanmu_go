package models

type DanMuMsg struct {
	UID        uint32 `json:"uid"`
	Uname      string `json:"uname"`
	Ulevel     uint32 `json:"ulevel"`
	Text       string `json:"text"`
	MedalLevel uint32 `json:"medal_level"`
	MedalName  string `json:"medal_name"`
}

func NewDanmu() *DanMuMsg {
	return &DanMuMsg{
		UID:        0,
		Uname:      "",
		Ulevel:     0,
		Text:       "",
		MedalLevel: 0,
		MedalName:  "无勋章",
	}
}

type Gift struct {
	UUname   string `json:"u_uname"`
	Action   string `json:"action"`
	Price    uint32 `json:"price"`
	GiftName string `json:"gift_name"`
}

func NewGift() *Gift {
	return &Gift{
		UUname:   "",
		Action:   "",
		Price:    0,
		GiftName: "",
	}
}

type WelCome struct {
}

type Notice struct {
}

type CMD string

var (
	CMDDanmuMsg                  CMD = "DANMU_MSG"                     // 普通弹幕信息
	CMDSendGift                  CMD = "SEND_GIFT"                     // 普通的礼物，不包含礼物连击
	CMDWELCOME                   CMD = "WELCOME"                       // 欢迎VIP
	CMDWelcomeGuard              CMD = "WELCOME_GUARD"                 // 欢迎房管
	CMDEntry                     CMD = "ENTRY_EFFECT"                  // 欢迎舰长等头衔
	CMDRoomRealTimeMessageUpdate CMD = "ROOM_REAL_TIME_MESSAGE_UPDATE" // 房间关注数变动
)

func (d *DanMuMsg) GetDanmuMsg(source []byte) {
	d.UID = Json.Get(source, "info", 2, 0).ToUint32()
	d.Uname = Json.Get(source, "info", 2, 1).ToString()
	d.Ulevel = Json.Get(source, "info", 4, 0).ToUint32()
	d.Text = Json.Get(source, "info", 1).ToString()
	d.MedalName = Json.Get(source, "info", 3, 1).ToString()
	if d.MedalName == "" {
		d.MedalName = "无勋章"
	}
	d.MedalLevel = Json.Get(source, "info", 3, 0).ToUint32()
	return
}

func (g *Gift) GetGiftMsg(source []byte) {
	g.UUname = Json.Get(source, "data", "uname").ToString()
	g.Action = Json.Get(source, "data", "action").ToString()
	nums := Json.Get(source, "data", "num").ToUint32()
	g.Price = Json.Get(source, "data", "price").ToUint32() * nums
	g.GiftName = Json.Get(source, "data", "giftName").ToString()
}
