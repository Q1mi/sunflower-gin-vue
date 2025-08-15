package auth

import (
	"context"
	"sunflower-gin/internal/dao/query"
	"sunflower-gin/internal/model"
	"sunflower-gin/pkg/jwt"

	"go.uber.org/zap"
)

func RefreshToken(ctx context.Context, token string) (*model.RefreshTokenOutput, error) {
	// 1. 校验refreshToken是否有效
	claims, err := jwt.ParseRefreshToken(token)
	if err != nil {
		zap.L().Error("refreshToken校验失败", zap.Error(err))
		return nil, err
	}
	// 2. 解析得到userID
	userId := claims.UserId
	// 3. 根据userID查询用户信息
	userInst, err := query.Userinfo.WithContext(ctx).
		Where(query.Userinfo.UserID.Eq(userId)).
		First()
	if err != nil {
		zap.L().Error("根据userID查询用户信息失败", zap.Error(err))
		return nil, err
	}
	// 4. 生成新的accessToken和refreshToken
	accessToken, err := jwt.GenAccessToken(userInst.UserID, userInst.Username)
	if err != nil {
		zap.L().Error("生成新的accessToken失败", zap.Error(err))
		return nil, err
	}
	refreshToken, err := jwt.GenRefreshToken(userInst.UserID, userInst.Username)
	if err != nil {
		zap.L().Error("生成新的refreshToken失败", zap.Error(err))
		return nil, err
	}
	// 5. 返回新的token
	return &model.RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
