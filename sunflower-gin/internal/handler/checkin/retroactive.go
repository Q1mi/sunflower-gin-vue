package checkin

import (
	"sunflower-gin/api"
	v1 "sunflower-gin/api/checkin/v1"
	"sunflower-gin/internal/middleware"
	"sunflower-gin/internal/service/checkin"
	"time"

	"github.com/gin-gonic/gin"
)

// RetroactiveHandler 补签接口
func RetroactiveHandler(c *gin.Context) {
	// 1. 获取请求参数和当前用户
	var req v1.RetroReq
	if err := c.ShouldBindJSON(&req); err != nil {
		api.ResponseError(c, api.CodeInvalidParam)
		return
	}
	userID := c.Value(middleware.CtxKeyUserID).(int64)
	if userID == 0 {
		api.ResponseError(c, api.CodeNeedLogin)
		return
	}
	// 校验日期格式（最基本的格式校验）
	t, err := time.Parse(time.DateOnly, req.Date)
	if err != nil {
		api.ResponseError(c, api.CodeInvalidParam)
		return
	}
	// 2. 调用service层补签逻辑
	if err := checkin.Retroactive(c, userID, t); err != nil {
		api.ResponseErrorWithMsg(c, api.CodeServerBusy, err.Error())
		return
	}
	// 3. 返回响应
	api.ResponseSuccess(c, &v1.RetroResp{})
}
