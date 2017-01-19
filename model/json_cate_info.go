package model

type CateInfoJson struct {
	Error int         `json:"error"`
	Data  []*CateInfo `json:"data"`
}

type CateInfo struct {
	CateID    string `json:"cate_id"`
	GameName  string `json:"game_name"`
	ShortName string `json:"short_name"`
	GameURL   string `json:"game_url"`
	GameSrc   string `json:"game_src"`
	GameIcon  string `json:"game_icon"`
}
