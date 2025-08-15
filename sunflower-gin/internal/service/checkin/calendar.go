package checkin

import (
	"context"
	"fmt"
	"sunflower-gin/internal/dao"
	"sunflower-gin/internal/model"
	"time"

	"go.uber.org/zap"
)

// 签到日历

func MonthDetail(ctx context.Context, userID int64, t time.Time) (*model.MonthDetailOutput, error) {
	// 1. 取出当月所有签到记录和补签记录
	checkinBitmap, retroBitmap, err := getMonthBitmap(ctx, userID, t.Year(), int(t.Month()))
	if err != nil {
		zap.L().Error("getMonthBitmap error", zap.Error(err))
		return nil, err
	}
	zap.L().Sugar().Debugf("checkinBitmap: %031b, retroBitmap: %031b\n", checkinBitmap, retroBitmap)
	// 当月多少天
	firstOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	dayNum := lastOfMonth.Day()
	checkinDays := parseBitmap2Days(checkinBitmap, dayNum)
	retroDays := parseBitmap2Days(retroBitmap, dayNum)

	// 2. 计算连续签到天数
	bitmap := checkinBitmap | retroBitmap
	zap.L().Sugar().Debugf("bitmap: %031b\n", bitmap)
	maxConsecutive, err := calcMonthConsecutiveDays(ctx, bitmap, dayNum)
	if err != nil {
		zap.L().Error("calcMonthConsecutiveDays error", zap.Error(err))
		return nil, err
	}
	// 3. 计算剩余补签次数
	remainRetroTimes := maxRetroTimesPerMonth - len(retroDays)
	// 4. 计算当天是否已签到
	now := time.Now()
	isCheckedToday := checkinBitmap&(1<<(dayNum-now.Day())) != 0
	// 5. 返回
	return &model.MonthDetailOutput{
		CheckedInDays:      checkinDays,
		RetroCheckedInDays: retroDays,
		ConsecutiveDays:    maxConsecutive,
		RemainRetroTimes:   remainRetroTimes,
		IsCheckedInToday:   isCheckedToday,
	}, nil
}

// IsCheckedToday 判断今天是否已签到
func IsCheckedToday(ctx context.Context, userID int64) (bool, error) {
	now := time.Now()
	year := now.Year()
	key := fmt.Sprintf(yearSignKeyFormat, userID, year)
	dayOffset := now.YearDay() - 1 // 偏移量从0开始
	value, err := dao.RedisClient.GetBit(ctx, key, int64(dayOffset)).Result()
	if err != nil {
		zap.L().Error("getBit error", zap.Error(err))
		return false, err
	}
	return value == 1, nil
}

// parseBitmap2Days 解析 bitmap 到日期列表
func parseBitmap2Days(bitmap uint64, dayNum int) []int {
	// 111100000001111111  --> [1,2,3,4,12,13,14,...]
	days := make([]int, 0)
	for i := range dayNum {
		if bitmap&(1<<(dayNum-1-i)) != 0 {
			// 当天已签到
			days = append(days, i+1)
		}
	}
	return days
}
