package response

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Cost string      `json:"cost"`
}

type Response struct {
	httpCode int
	result   *Result
}
