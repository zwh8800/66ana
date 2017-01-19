package model

import "encoding/json"

type RoomInfoJson struct {
	Error int             `json:"error"`
	Data  json.RawMessage `json:"data"`
}

type RoomGiftInfo struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Pc    float64 `json:"pc"`
	Gx    int     `json:"gx"`
	Desc  string  `json:"desc"`
	Intro string  `json:"intro"`
	Mimg  string  `json:"mimg"`
	Himg  string  `json:"himg"`
}

type RoomInfo struct {
	RoomID      string          `json:"room_id"`
	RoomThumb   string          `json:"room_thumb"`
	CateID      string          `json:"cate_id"`
	CateName    string          `json:"cate_name"`
	RoomName    string          `json:"room_name"`
	RoomStatus  string          `json:"room_status"`
	OwnerName   string          `json:"owner_name"`
	Avatar      string          `json:"avatar"`
	Online      int             `json:"online"`
	OwnerWeight string          `json:"owner_weight"`
	FansNum     string          `json:"fans_num"`
	StartTime   string          `json:"start_time"`
	Gift        []*RoomGiftInfo `json:"gift"`
}
