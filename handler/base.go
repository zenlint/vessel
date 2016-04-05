package handler

// BaseResp base responce for all responce
type BaseResp struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

const (
	// CodeSuccess :
	CodeSuccess = 200
	// CodeError :
	CodeError = 500

	// MsgSuccess :
	MsgSuccess = "success"
	// MsgError :
	MsgError = "error"
)
