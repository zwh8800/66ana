package model

type ReqPager struct {
	Offset int64 `json:"offset" form:"offset" query:"offset"`
	Limit  int64 `json:"limit" form:"limit" query:"limit"`
}
