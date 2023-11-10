package result

type ResponseSuccessBean struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

type HooksApiResponseSuccesBean struct {
	Data any `json:"data,omitempty"`
}

type IndexApiResponseSuccesBean struct {
	Data []byte `json:"data,omitempty"`
}

type NullJson struct{}

func Success(data any) *ResponseSuccessBean {
	return &ResponseSuccessBean{200, "success", data}
}

func HooksSuccess(data any) *HooksApiResponseSuccesBean {
	return &HooksApiResponseSuccesBean{data}
}

type ResponseErrorBean struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func Error(errCode int64, errMsg string) *ResponseErrorBean {
	return &ResponseErrorBean{errCode, errMsg}
}
