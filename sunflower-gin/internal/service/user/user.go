package user

import (
	"context"
	"errors"
	"sunflower-gin/internal/dao/query"
	"sunflower-gin/internal/model"
	"sunflower-gin/pkg/snowflake"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultAvatar = "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8dXNlciUyMHByb2ZpbGV8ZW58MHx8MHx8fDA=&auto=format&fit=crop&w=50&q=60"
)

var (
	ErrUserExist = errors.New("用户名已存在")
)

// 业务逻辑层

// Create 创建用户
func Create(ctx context.Context, input *model.CreateUserInput) (*model.CreateUserOutput, error) {
	// 1. 判断用户名是否存在
	// 传统的GORM查询
	// dao.DB.WithContext(ctx).Model(&model.Userinfo{}).Where("username = ?", input.Username)
	// GORM GEN 生成代码的查询
	count, err := query.Userinfo.WithContext(ctx).
		Where(query.Userinfo.Username.Eq(input.Username)).
		Count()
	if err != nil {
		zap.L().Error("Create: query userinfo failed", zap.Error(err))
		return nil, err
	}
	if count > 0 {
		return nil, ErrUserExist
	}
	// 2. 创建用户
	// 传统的密码加密是加点盐算个md5 之类的，
	// 进阶一点的做法是使用bcrypt库进行密码加密,及时相同的密码，每次生成的hash值都不一样
	// 密码加密 bcrypt 加密后的密码长度固定为 60 字符。
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("Create: bcrypt generate password failed", zap.Error(err))
		return nil, err
	}
	zap.L().Sugar().Debugf("--> hashPwd:%s\n", string(hashedPwd))
	// 使用雪花算法生成唯一的 user id
	uid, err := snowflake.NextID()
	if err != nil {
		zap.L().Error("Create: generate snowflake id failed", zap.Error(err))
		return nil, err
	}
	user := &model.Userinfo{
		UserID:   uid, // 每次注册，系统生成一个唯一的userid
		Username: input.Username,
		Password: string(hashedPwd), // 密码需要加密
		Email:    input.Email,
		Avatar:   defaultAvatar,
	}
	if err := query.Userinfo.WithContext(ctx).Create(user); err != nil {
		zap.L().Error("Create: create userinfo failed", zap.Error(err))
		return nil, err
	}
	// 3. 返回结果
	return &model.CreateUserOutput{
		UserId:   user.UserID,
		Username: user.Username,
	}, nil
}

// GetProfile 获取用户信息
func GetProfile(ctx context.Context, userID int64) (*model.UserProfileOutput, error) {
	// 根据userID查库
	userInst, err := query.Userinfo.WithContext(ctx).
		Where(query.Userinfo.UserID.Eq(userID)).
		First()
	if err != nil {
		zap.L().Error("GetProfile: query userinfo failed", zap.Error(err))
		return nil, err
	}
	// 封装返回结果
	return &model.UserProfileOutput{
		UserId:   userInst.UserID,
		Username: userInst.Username,
		Email:    userInst.Email,
		Avatar:   userInst.Avatar,
	}, nil
}
