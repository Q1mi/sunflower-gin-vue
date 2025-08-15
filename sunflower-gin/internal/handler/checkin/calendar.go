package checkin

import (
	"sunflower-gin/api"
	v1 "sunflower-gin/api/checkin/v1"
	"sunflower-gin/internal/middleware"
	"sunflower-gin/internal/service/checkin"
	"time"

	"github.com/gin-gonic/gin"
)

func CalendarHandler(c *gin.Context) {
	// 1. 获取请求参数并校验 userID
	var req v1.CalendarReq
	if err := c.ShouldBindQuery(&req); err != nil {
		api.ResponseError(c, api.CodeInvalidParam)
		return
	}
	userID := c.Value(middleware.CtxKeyUserID).(int64)
	// 解析请求参数的年月
	t, err := time.Parse("2006-01", req.YearMonth)
	if err != nil {
		api.ResponseErrorWithMsg(c, api.CodeInvalidParam, err.Error())
		return
	}
	// 2. 调用 service 层处理业务
	output, err := checkin.MonthDetail(c, userID, t)
	if err != nil {
		api.ResponseErrorWithMsg(c, api.CodeServerBusy, err.Error())
		return
	}
	// 3. 返回响应
	api.ResponseSuccess(c, &v1.CalendarResp{
		Year:  t.Year(),
		Month: int(t.Month()),
		Detail: v1.DetailInfo{
			CheckedInDays:      output.CheckedInDays,
			RetroCheckedInDays: output.RetroCheckedInDays,
			IsCheckedInToday:   output.IsCheckedInToday,
			RemainRetroTimes:   output.RemainRetroTimes,
			ConsecutiveDays:    output.ConsecutiveDays,
		},
	})
}
