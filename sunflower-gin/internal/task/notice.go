package task

import (
	"context"
	"fmt"
	"sunflower-gin/internal/dao"
	"time"

	_ "embed"
)

const (
	yearSignKeyFormat = "user:checkins:daily:%d:%d" // user:checkins:daily:12131321421312:2025
)

//go:embed remind.lua
var remindScript string

// CheckAndNotify 检查签到并发送通知
func CheckAndNotify(ctx context.Context, remindThreshold int) error {
	fmt.Println("start check and notify...")
	// 1. 获取所有符合条件的用户
	// 可以通过扫 MySQL 用户表找到最近几天有登录过的（如果有现成的用户登录时间记录可以直接拿来用）
	// 或者可以在用户签到的时候记录一个 ZSet, userID:签到时间戳（如果用户量多需要拆分 Key）
	userIDs := []uint64{25016147980058993}

	// 2. 加载 lua script
	sha, err := dao.RedisClient.ScriptLoad(ctx, remindScript).Result()
	if err != nil {
		fmt.Printf("ScriptLoad err: %v\n", err)
		return err
	}
	// 3. 遍历判断每个用户
	now := time.Now()
	for _, userID := range userIDs {
		key := fmt.Sprintf(yearSignKeyFormat, userID, time.Now().Year())

		// 计算当前日偏移量(当年第几天)
		dayOfYearOffset := now.YearDay() - 1
		fmt.Printf("key: %s, dayOfYearOffset: %d\n", key, dayOfYearOffset)
		// 执行LUA脚本
		result, err := dao.RedisClient.EvalSha(ctx, sha, []string{key}, dayOfYearOffset, remindThreshold).Int()
		fmt.Printf("result: %d, err: %v\n", result, err)
		if err != nil {
			return err
		}

		if result == 1 {
			fmt.Printf("用户%d需要发送断签提醒\n", userID)
			// 发送到消息队列，执行后续的推送逻辑（APP Push 或 短信等）
		}
	}
	return nil
}
