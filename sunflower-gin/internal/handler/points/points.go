package points

import (
	"sunflower-gin/api"
	v1 "sunflower-gin/api/points/v1"
	"sunflower-gin/internal/middleware"
	"sunflower-gin/internal/model"
	"sunflower-gin/internal/service/points"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit  = 10 // 默认分页大小
	maxLimit      = 50 // 最大分页大小
	defaultOffset = 0  // 默认偏移量
)

// SummaryHandler 获取积分信息
func SummaryHandler(c *gin.Context) {
	// 1. 获取当前用户信息
	userID := c.Value(middleware.CtxKeyUserID).(int64)
	if userID == 0 {
		api.ResponseError(c, api.CodeNeedLogin)
		return
	}
	// 2. 调用 service 层获取积分信息
	output, err := points.Summary(c, userID)
	if err != nil {
		api.ResponseErrorWithMsg(c, api.CodeServerBusy, err.Error())
		return
	}
	// 3. 返回积分信息
	api.ResponseSuccess(c, &v1.SummaryResp{Total: output.TotalPoint})
}

// RecordsHandler 获取积分记录
func RecordsHandler(c *gin.Context) {
	// 1. 获取当前用户信息和分页信息
	// c.GetQuery("limit")
	var req v1.RecordsReq
	if err := c.ShouldBind(&req); err != nil {
		api.ResponseError(c, api.CodeInvalidParam)
		return
	}
	userID := c.Value(middleware.CtxKeyUserID).(int64)
	if userID == 0 {
		api.ResponseError(c, api.CodeNeedLogin)
		return
	}
	// 分页参数校验
	if req.Limit <= 0 || req.Limit > maxLimit {
		req.Limit = defaultLimit
	}
	if req.Offset < 0 {
		req.Offset = defaultOffset
	}
	// 2. 调用 service 层获取积分记录
	output, err := points.Records(c, &model.RecordsInput{
		UserID: userID,
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		api.ResponseErrorWithMsg(c, api.CodeServerBusy, err.Error())
		return
	}
	// 3. 返回积分记录
	// 3.1 组装返回数据
	list := make([]*v1.RecordInfo, 0, len(output.List))
	for _, item := range output.List {
		list = append(list, &v1.RecordInfo{
			PointsChange:    item.PointsChange,
			TransactionType: item.TransactionType,
			Description:     item.Description,
			TransactionTime: item.TransactionTime,
		})
	}
	api.ResponseSuccess(c, &v1.RecordsResp{
		Total:   output.Total,
		HasMore: output.HasMore,
		List:    list,
	})
}
