package task

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

// 定时任务

/*
github.com/robfig/cron/v3 CRON 表达式格式

Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
*/

func MustInit(ctx context.Context) *cron.Cron {
	tz, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}
	c := cron.New(cron.WithLocation(tz))
	// 添加定时任务
	c.AddFunc("25 20 * * *", func() { CheckAndNotify(ctx, 2) })
	c.Start()
	return c
}
