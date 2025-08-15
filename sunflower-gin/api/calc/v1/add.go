package v1

// Req 加法请求结构体
type Req struct {
	X int `json:"x" binding:"required"`
	Y int `json:"y" binding:"required"`
}

// Resp 加法响应结构体
type Resp struct {
	Result int `json:"result"`
}
