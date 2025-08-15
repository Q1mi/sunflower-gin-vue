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

// 每日签到业务逻辑

const (
	yearSignKeyFormat   = "user:checkins:daily:%d:%d"      // user:checkins:daily:123213131:2025
	monthRetroKeyFormat = "user:checkins:retro:%d:%d:%02d" // user:checkins:retro:123213131:2025:01
)

const (
	defaultDailyPoints    = 1 // 默认每日签到积分
	maxRetroTimesPerMonth = 3 // 每月最多补签次数
)

// 积分变更记录表的交易类型
type PointsTransactionType int32

const (
	PointsTransactionTypeDaily       PointsTransactionType = 1 // 每日签到 1
	PointsTransactionTypeConsecutive PointsTransactionType = 2 // 连续签到 2
	PointsTransactionTypeRetroactive PointsTransactionType = 3 // 补签 3
)

// 积分变更记录表的交易类型对应的描述信息
var pointsTransactionTypeDescMap = map[PointsTransactionType]string{
	PointsTransactionTypeDaily:       "每日签到奖励",
	PointsTransactionTypeConsecutive: "连续签到奖励",
	PointsTransactionTypeRetroactive: "补签%s消耗",
}

// 定义连续签到的奖励类型和描述
type ConsecutiveBonusType int32

const (
	consecutiveBonus3  ConsecutiveBonusType = 1 // 连续签到3天
	consecutiveBonus7  ConsecutiveBonusType = 2 // 连续签到7天
	consecutiveBonus15 ConsecutiveBonusType = 3 // 连续签到15天
	consecutiveBonus30 ConsecutiveBonusType = 4 // 月度满签
)

var consecutiveBonusNameMap = map[ConsecutiveBonusType]string{
	consecutiveBonus3:  "连续签到3天奖励",
	consecutiveBonus7:  "连续签到7天奖励",
	consecutiveBonus15: "连续签到15天奖励",
	consecutiveBonus30: "月度满签奖励",
}

// 连续签到奖励的触发规则
type consecutiveBonusRule struct {
	BonusType   ConsecutiveBonusType // 奖励类型
	TriggerDays int                  // 需要连续签到多少天才能触发这个规则
	Points      int64                // 发放的积分数量
}

var consecutiveBonusRuleList = []consecutiveBonusRule{
	{TriggerDays: 3, Points: 5, BonusType: consecutiveBonus3},
	{TriggerDays: 7, Points: 10, BonusType: consecutiveBonus7},
	{TriggerDays: 15, Points: 20, BonusType: consecutiveBonus15},
	{TriggerDays: 28, Points: 100, BonusType: consecutiveBonus30},
}

var (
	ErrCheckedIn = errors.New("今日已签到") // 已签到
)

// Daily 每日签到处理函数
func Daily(ctx context.Context, userID int64) error {
	// setbit key offset 1
	now := time.Now()
	year := now.Year()
	key := fmt.Sprintf(yearSignKeyFormat, userID, year)
	// 1. 获取今天是今年的第几天，算出 offset
	// now.YearDay() // 今天是今年的第几天
	offset := now.YearDay() - 1 // offset 从 0 开始
	// 2、Redis 中执行 setbit 操作
	zap.L().Sugar().Debugf("--> daily setbit key: %s, offset: %d", key, offset)
	ret, err := dao.RedisClient.SetBit(ctx, key, int64(offset), 1).Result()
	if err != nil {
		zap.L().Error("daily setbit error", zap.Error(err))
		return err
	}
	if ret == 1 {
		// 已签到
		return ErrCheckedIn
	}
	// 3. 发放每日签到积分
	err = addPoints(ctx, &model.AddPointInput{
		UserID:      userID,
		PointAmount: defaultDailyPoints,
		Type:        int32(PointsTransactionTypeDaily),
		Desc:        pointsTransactionTypeDescMap[PointsTransactionTypeDaily],
	})
	if err != nil {
		zap.L().Error("addPoints error", zap.Error(err))
		return err
	}
	// 4. 发放连续签到奖励
	// 1 1 0 1 1 0 1
	return updateConsecutiveBonus(ctx, userID, year, int(now.Month()))
}

// updateConsecutiveBonus 更新连续签到奖励
func updateConsecutiveBonus(ctx context.Context, userID int64, year, month int) error {
	// 1. 获取当前的连续签到天数
	checkinBitmap, retroBitmap, err := getMonthBitmap(ctx, userID, year, month)
	if err != nil {
		zap.L().Error("getMonthBitmap error", zap.Error(err))
		return err
	}
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	dayNum := lastOfMonth.Day()
	bitmap := checkinBitmap | retroBitmap
	maxConsecutive, err := calcMonthConsecutiveDays(ctx, bitmap, dayNum)
	if err != nil {
		zap.L().Error("calcMonthConsecutiveDays error", zap.Error(err))
		return err
	}
	// 2. 计算连续签到的奖励积分
	// 2.1 先查询用户当月领取了哪些连续签到奖励
	bonusLogList, err := query.UserMonthlyBonusLog.WithContext(ctx).
		Where(query.UserMonthlyBonusLog.UserID.Eq(userID)).
		Where(query.UserMonthlyBonusLog.YearMonth.Eq(fmt.Sprintf("%d%02d", year, month))).
		Find()
	if err != nil {
		zap.L().Error("query user_monthly_bonus_log error", zap.Error(err))
		return err
	}
	bonusLogMap := make(map[ConsecutiveBonusType]bool, len(bonusLogList))
	for _, v := range bonusLogList {
		bonusLogMap[ConsecutiveBonusType(v.BonusType)] = true
	}
	for _, rule := range consecutiveBonusRuleList {
		if maxConsecutive >= rule.TriggerDays && !bonusLogMap[rule.BonusType] {
			// 2.1.1 发放连续签到奖励积分
			// 更新 user_points 表 和 user_points_transactions 表，
			err := addPoints(ctx, &model.AddPointInput{
				UserID:      userID,
				PointAmount: rule.Points,
				Type:        int32(PointsTransactionTypeConsecutive),
				Desc:        consecutiveBonusNameMap[rule.BonusType],
			})
			if err != nil {
				zap.L().Error("[NEED_HANDLE] updateConsecutiveBonus addPoints error", zap.Error(err))
				return err
			}
			// 并且记录连续签到奖励日志表 bonus_log 表
			err = query.UserMonthlyBonusLog.WithContext(ctx).
				Create(&model.UserMonthlyBonusLog{
					UserID:      userID,
					YearMonth:   fmt.Sprintf("%d%02d", year, month),
					BonusType:   int32(rule.BonusType),
					Description: consecutiveBonusNameMap[rule.BonusType],
				})
			if err != nil {
				zap.L().Error("[NEED_HANDLE] updateConsecutiveBonus create user_monthly_bonus_log error", zap.Error(err))
				continue
			}
		}
	}
	return nil
}

// calcMonthConsecutiveDays 计算本月连续签到天数
func calcMonthConsecutiveDays(ctx context.Context, bitmap uint64, dayNum int) (int, error) {
	// 3. 取并集之后计算本月连续签到天数
	// bitmap := checkinBitmap | retroBitmap
	// 计算连续签到天数
	maxCount := 0
	currCount := 0
	for i := range dayNum {
		// 1000100111111111100000000000
		// 000001
		if bitmap&(1<<(dayNum-1-i)) != 0 {
			currCount++
			if currCount > maxCount {
				maxCount = currCount
			}
		} else {
			currCount = 0
		}
	}
	if currCount > maxCount { // 循环结束再比较一次
		maxCount = currCount
	}
	return maxCount, nil
}

// getMonthBitmap 获取当月签到记录和补签记录的 bitmap
func getMonthBitmap(ctx context.Context, userID int64, year, month int) (uint64, uint64, error) {
	// 1. 获取月度正常签到数据
	// 1.1 获取当月的天数 --> 我拿到当月的最后一天是几号就知道当月有多少天了
	// 取当月的第一天，然后+1个月，再减去一天，就能得到当月的最后一天是几号
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	dayNum := lastOfMonth.Day()
	offset := firstOfMonth.YearDay() - 1 // 当月第一天的是一年的第几天-1 =  offset
	// 取到年度签到记录 key
	key := fmt.Sprintf(yearSignKeyFormat, userID, year)
	// 从 年度签到数据中取出当月的签到记录
	bitWidthType := fmt.Sprintf("u%d", dayNum) // u31 表示无符号 31 位整数
	zap.L().Sugar().Debugf("key:%s bitWidthType:%s offset:%d\n", key, bitWidthType, offset)
	values, err := dao.RedisClient.BitField(ctx, key, "GET", bitWidthType, offset).Result()
	if err != nil {
		zap.L().Error("获取用户签到记录失败", zap.Error(err))
		return 0, 0, err
	}
	zap.L().Sugar().Debugf("checkin values:%#v\n", values)
	if len(values) == 0 {
		values = []int64{0}
	}
	checkinBitmap := uint64(values[0])
	// 取 月度 补签数据
	retroKey := fmt.Sprintf(monthRetroKeyFormat, userID, year, month)
	retroValues, err := dao.RedisClient.BitField(ctx, retroKey, "GET", bitWidthType, "#0").Result()
	if err != nil {
		zap.L().Error("获取用户补签记录失败", zap.Error(err))
		return 0, 0, err
	}
	zap.L().Sugar().Debugf("retro values:%#v\n", values)
	if len(retroValues) == 0 { // 用户当月可能没有补签记录
		retroValues = []int64{0}
	}
	retroBitmap := uint64(retroValues[0])
	return checkinBitmap, retroBitmap, nil
}

func addPoints(ctx context.Context, input *model.AddPointInput) error {
	// 需要分别更新 user_points 表和 user_points_transactions 表
	// 3.1 查询 user_points 表
	userPoint, err := query.UserPoint.WithContext(ctx).
		Where(query.UserPoint.UserID.Eq(input.UserID)).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.L().Error("query user_points error", zap.Error(err))
		return err
	}
	if userPoint == nil || userPoint.ID == 0 {
		// 没有查到记录，说明是第一次来签到
		userPoint = &model.UserPoint{
			UserID: input.UserID,
		}
	}
	// 更新积分
	userPoint.Points = userPoint.Points + input.PointAmount
	userPoint.PointsTotal = userPoint.PointsTotal + input.PointAmount
	// 在事务中更新 user_points 表和 user_points_transactions 表
	err = query.Q.Transaction(func(tx *query.Query) error {
		// 更新 user_points 表
		if err := tx.UserPoint.WithContext(ctx).Save(userPoint); err != nil {
			zap.L().Error("tx save user_points error", zap.Error(err))
			return err
		}
		// 更新 user_points_transactions 表
		if err := tx.UserPointsTransaction.WithContext(ctx).
			Create(&model.UserPointsTransaction{
				UserID:          input.UserID,
				PointsChange:    input.PointAmount,
				CurrentBalance:  userPoint.Points,
				TransactionType: input.Type,
				Description:     input.Desc,
			}); err != nil {
			zap.L().Error("tx create user_points_transactions error", zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		zap.L().Error("tx commit failed", zap.Error(err))
		return err
	}
	return nil
}
