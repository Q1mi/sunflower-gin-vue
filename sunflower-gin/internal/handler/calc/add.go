package calc

import (
	"github.com/gin-gonic/gin"
	"sunflower-gin/api"
	v1 "sunflower-gin/api/calc/v1"
)

// AddHandler 加法处理函数
func AddHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req v1.Req
		if err := c.ShouldBindJSON(&req); err != nil {
			api.ResponseError(c, api.CodeInvalidParam)
			return
		}
		ret := v1.Resp{
			Result: req.X + req.Y,
		}
		api.ResponseSuccess(c, ret)
	}
}
