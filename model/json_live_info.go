package model

type LiveInfoJson struct {
	Error int         `json:"error"`
	Data  []*LiveInfo `json:"data"`
}

type LiveInfo struct {
	RoomID          string `json:"room_id"`
	RoomSrc         string `json:"room_src"`
	VerticalSrc     string `json:"vertical_src"`
	IsVertical      int    `json:"isVertical"`
	CateID          string `json:"cate_id"`
	RoomName        string `json:"room_name"`
	ShowStatus      string `json:"show_status"`
	Subject         string `json:"subject"`
	ShowTime        string `json:"show_time"`
	OwnerUID        string `json:"owner_uid"`
	SpecificCatalog string `json:"specific_catalog"`
	SpecificStatus  string `json:"specific_status"`
	VodQuality      string `json:"vod_quality"`
	Nickname        string `json:"nickname"`
	Online          int    `json:"online"`
	URL             string `json:"url"`
	GameURL         string `json:"game_url"`
	GameName        string `json:"game_name"`
	ChildID         string `json:"child_id"`
	Avatar          string `json:"avatar"`
	AvatarMid       string `json:"avatar_mid"`
	AvatarSmall     string `json:"avatar_small"`
	JumpURL         string `json:"jumpUrl"`
	IsHide          int    `json:"isHide,omitempty"`
	Fans            string `json:"fans"`
	Ranktype        int    `json:"ranktype"`
	AnchorCity      string `json:"anchor_city"`
}
