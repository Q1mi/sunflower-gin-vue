package user

import (
	"errors"
	"sunflower-gin/api"
	v1 "sunflower-gin/api/user/v1"
	"sunflower-gin/internal/middleware"
	"sunflower-gin/internal/model"
	"sunflower-gin/internal/service/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// hanlder 层

// CreateHandler 创建用户接口
func CreateHandler(c *gin.Context) {
	// 1. 获取请求参数&校验参数
	var req v1.CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		// 请求参数有问题
		zap.L().Error("CreateHandler: ShouldBindJSON failed", zap.Error(err))
		api.ResponseError(c, api.CodeInvalidParam)
		return
	}
	zap.L().Sugar().Debugf("---> CreateHandler: %+v", req)
	// 2. 执行业务逻辑
	input := &model.CreateUserInput{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}
	output, err := user.Create(c, input)
	if err != nil {
		// 如果是用户已存在错误
		if errors.Is(err, user.ErrUserExist) {
			api.ResponseError(c, api.CodeUserExist)
			return
		}
		// 其它错误，统一返回服务繁忙
		api.ResponseError(c, api.CodeServerBusy)
		return
	}
	// 3. 返回响应,数据应该是定义好的想用结构体
	api.ResponseSuccess(c, &v1.CreateRes{
		UserId:   output.UserId,
		Username: output.Username,
	})
}

func ProfileHandler(c *gin.Context) {
	// 1. 获取请求参数&校验参数
	// 从上下文中获取 userID
	userID := c.Value(middleware.CtxKeyUserID).(int64)
	if userID == 0 {
		api.ResponseError(c, api.CodeNeedLogin)
		return
	}
	// 2. 执行业务逻辑
	output, err := user.GetProfile(c, userID)
	if err != nil {
		api.ResponseError(c, api.CodeServerBusy)
		return
	}
	// 3. 返回响应
	api.ResponseSuccess(c, &v1.MeRes{
		Username: output.Username,
		Email:    output.Email,
		Avatar:   output.Avatar,
	})
}
