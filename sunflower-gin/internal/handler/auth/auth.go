package auth

import (
	"sunflower-gin/api"
	v1 "sunflower-gin/api/auth/v1"
	"sunflower-gin/internal/model"
	"sunflower-gin/internal/service/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoginHandler 登录处理函数
func LoginHandler(c *gin.Context) {
	// 1. 获取请求参数并进行参数校验
	var req v1.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("参数校验失败", zap.Error(err))
		api.ResponseError(c, api.CodeInvalidParam)
		return
	}
	// 2. 调用用户登录服务
	output, err := auth.Login(c, &model.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		zap.L().Error("用户登录失败", zap.Error(err))
		api.ResponseErrorWithMsg(c, api.CodeInvalidPassword, err.Error())
		return
	}
	// 3. 拼装响应数据并返回
	api.ResponseSuccess(c, v1.LoginResp{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	})
}

func RefreshHandler(c *gin.Context) {
	// 1. 获取请求参数并进行参数校验
	var req v1.RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("参数校验失败", zap.Error(err))
		api.ResponseError(c, api.CodeInvalidParam)
		return
	}
	// 2. 调用service层刷新token
	output, err := auth.RefreshToken(c, req.RefreshToken)
	if err != nil {
		zap.L().Error("刷新token失败", zap.Error(err))
		api.ResponseErrorWithMsg(c, api.CodeInvalidToken, err.Error())
		return
	}
	// 3. 拼装响应数据并返回
	api.ResponseSuccess(c, v1.RefreshResp{
		AccessToken:  output.AccessToken, // 新的accessToken
		RefreshToken: output.RefreshToken,
	})
}
