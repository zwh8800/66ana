package model

const (
	CodeOK                  = "OK"
	CodeInternalServerError = "InternalServerError"
)

const (
	MessageOK                  = "操作成功"
	MessageInternalServerError = "操作失败，请稍后重试"
)

type ApiGeneral struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewApiOKResponse() *ApiGeneral {
	return NewApiResponse(CodeOK, MessageOK)
}

func NewApiResponse(code string, message string) *ApiGeneral {
	return &ApiGeneral{
		Code:    code,
		Message: message,
	}
}
