package auth

import (
	"context"
	"errors"
	"sunflower-gin/internal/dao/query"
	"sunflower-gin/internal/model"
	"sunflower-gin/pkg/jwt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// 认证相关接口

// Login 登录接口
func Login(ctx context.Context, input *model.LoginInput) (*model.LoginOutput, error) {
	// 1. 登录校验
	// query.Userinfo.WithContext(ctx).
	// Where(query.Userinfo.Username.Eq(input.Username)).
	// Where(query.Userinfo.Password.Eq(input.Password)).
	// Find()
	userInst, err := query.Userinfo.WithContext(ctx).
		Where(query.Userinfo.Username.Eq(input.Username)).
		First()
	if err != nil {
		zap.L().Error("Login: query user failed", zap.Error(err))
		return nil, errors.New("用户名或密码错误")
	}
	// userInst.Password  // 加密之后的 password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(userInst.Password), []byte(input.Password)); err != nil {
		zap.L().Error("Login: password compare failed", zap.Error(err))
		return nil, errors.New("用户名或密码错误")
	}
	// 2. 如果登录成功，生成token
	// 2.1 生成access token
	accessToken, err := jwt.GenAccessToken(userInst.UserID, userInst.Username)
	if err != nil {
		zap.L().Error("Login: generate access token failed", zap.Error(err))
		return nil, errors.New("生成accessToken失败")
	}
	// 2.2 生成refresh token
	refreshToken, err := jwt.GenRefreshToken(userInst.UserID, userInst.Username)
	if err != nil {
		zap.L().Error("Login: generate refresh token failed", zap.Error(err))
		return nil, errors.New("生成refreshToken失败")
	}
	// 3. 返回token
	return &model.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
