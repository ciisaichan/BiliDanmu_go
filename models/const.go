package models

import jsoniter "github.com/json-iterator/go"

const (
	RealID      = "http://api.live.bilibili.com/room/v1/Room/room_init" // params: id=xxx
	DanMuServer = "broadcastlv.chat.bilibili.com:443"
	KeyUrl      = "https://api.live.bilibili.com/room/v1/Danmu/getConf"                 // params: room_id=xxx&platform=pc&player=web
	RoomInfoUrl = "https://api.live.bilibili.com/xlive/web-room/v1/index/getInfoByRoom" // params: room_id=xxx
)

var (
	Json = jsoniter.ConfigCompatibleWithStandardLibrary
)
