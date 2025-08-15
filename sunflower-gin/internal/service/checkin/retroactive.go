package checkin

import (
	"context"
	"errors"
	"fmt"
	"sunflower-gin/internal/dao"
	"sunflower-gin/internal/dao/query"
	"sunflower-gin/internal/model"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 补签相关的业务逻辑

const (
	defaultRetroCostPoints = 100 // 补签消耗积分
)

var (
	ErrInvalidRetroDate    = errors.New("无效的补签日期")
	ErrRetroNoTimes        = errors.New("本月已经没有补签次数了")
	ErrRetroNoEnoughPoints = errors.New("积分不足，无法补签")
)

// Retroactive 补签逻辑
func Retroactive(ctx context.Context, userID int64, date time.Time) error {
	// 1. 补签日期的校验（涉及业务逻辑的参数有效校验）
	if err := checkRetroDate(ctx, userID, date); err != nil {
		return err
	}
	// 2. 执行补签逻辑
	// 2.1 在 Redis 中标记补签的日期， setbit 设置补签记录
	key := fmt.Sprintf(monthRetroKeyFormat, userID, date.Year(), int(date.Month()))
	offset := date.Day() - 1 // 0 base index
	err := dao.RedisClient.SetBit(ctx, key, int64(offset), 1).Err()
	if err != nil {
		zap.L().Error("setbit error", zap.Error(err))
		return err
	}
	// 2.2 补签消耗积分，签到增加积分，增加积分记录到数据库
	if err := retroWithTransaction(ctx, userID, date); err != nil {
		// 如果补签逻辑执行失败，需要回滚 Redis 中的标记
		if err := dao.RedisClient.SetBit(ctx, key, int64(offset), 0).Err(); err != nil {
			return fmt.Errorf("retroWithTransaction rollback retro bit error:%w", err)
		}
		return err
	}
	// 3. 发放可能存在的连续签到奖励
	return updateConsecutiveBonus(ctx, userID, date.Year(), int(date.Month()))
}

// checkRetroDate 校验补签日期是否合法
func checkRetroDate(ctx context.Context, userID int64, date time.Time) error {
	// 1. 补签日期不能是今天或者未来的日期
	// 2. 补签的日期只能是当前月份的
	now := time.Now()
	if date.Year() > now.Year() ||
		date.Month() != now.Month() ||
		(date.Year() == now.Year() && date.YearDay() >= now.YearDay()) {
		return ErrInvalidRetroDate
	}
	// 3. 补签的日期不能是已经签到或者补签的日期
	checkinBitmap, retroBitmap, err := getMonthBitmap(ctx, userID, date.Year(), int(date.Month()))
	if err != nil {
		zap.L().Error("getMonthBitmap error", zap.Error(err))
		return err
	}
	bitmap := checkinBitmap | retroBitmap // 1111111111111111111011111111111111
	// 根据补签的日期去取对应位数的二进制值
	// 获取当月的总天数
	days := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, -1).Day()
	if bitmap&(1<<uint(days-date.Day())) != 0 {
		return ErrInvalidRetroDate
	}
	// 4. 补签的次数不能超过限制（3次）
	// 统计 retroBitmap 里有几个二进制位是1
	count := 0
	for retroBitmap != 0 {
		retroBitmap &= (retroBitmap - 1) // 去掉最右边的二进制位，直到为0
		count++
	}
	if count >= 3 {
		return ErrRetroNoTimes
	}
	return nil
}

// 补签逻辑，涉及到事务的处理
func retroWithTransaction(ctx context.Context, userID int64, date time.Time) error {
	return query.Q.Transaction(func(tx *query.Query) error {
		// 1. 查询用户当前的积分，积分不够也不能补签
		var (
			upInst *model.UserPoint
			err    error
		)
		upInst, err = tx.UserPoint.WithContext(ctx).
			Where(tx.UserPoint.UserID.Eq(userID)).
			First()
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err // 非预期的数据库查询错误，需要返回
			}
			upInst = &model.UserPoint{ // 初始化一个新的 UserPoint 实例，默认0积分
				UserID: userID,
			}
		}
		if upInst.Points < defaultRetroCostPoints {
			return ErrRetroNoEnoughPoints
		}
		// 2. 扣除积分
		pointsChange := -defaultRetroCostPoints          // 扣除积分
		newPoints := upInst.Points + int64(pointsChange) // 当前积分值
		// 3. 增加积分记录流水
		retroCostRecord := &model.UserPointsTransaction{
			UserID:          userID,
			PointsChange:    int64(pointsChange),
			CurrentBalance:  newPoints,
			TransactionType: int32(PointsTransactionTypeRetroactive),
			Description:     fmt.Sprintf(pointsTransactionTypeDescMap[PointsTransactionTypeRetroactive], date.Format(time.DateOnly)),
		}
		if err := tx.WithContext(ctx).UserPointsTransaction.Create(retroCostRecord); err != nil {
			zap.L().Error("create retroCostRecord error", zap.Error(err))
			return err
		}
		// 4. 更新用户积分
		upInst.Points = newPoints
		if err := tx.UserPoint.WithContext(ctx).Save(upInst); err != nil {
			zap.L().Error("update user points error", zap.Error(err))
			return err
		}
		return nil
	})
}
