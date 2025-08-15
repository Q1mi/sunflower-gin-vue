package checkin

import (
	"sunflower-gin/api"
	v1 "sunflower-gin/api/checkin/v1"
	"sunflower-gin/internal/middleware"
	"sunflower-gin/internal/service/checkin"

	"github.com/gin-gonic/gin"
)

// DailyHandler 每日签到接口处理函数
func DailyHandler(c *gin.Context) {
	// 1. 获取请求参数 userID
	userID := c.Value(middleware.CtxKeyUserID).(int64)
	if userID == 0 {
		api.ResponseError(c, api.CodeNeedLogin)
		return
	}
	// 2. 调用 service 层处理业务
	err := checkin.Daily(c, userID)
	if err != nil {
		api.ResponseErrorWithMsg(c, api.CodeServerBusy, err.Error())
		return
	}
	// 3. 返回响应
	api.ResponseSuccess(c, &v1.DailyResp{})
}
