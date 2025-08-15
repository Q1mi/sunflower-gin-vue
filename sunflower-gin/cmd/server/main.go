package main

import (
	"context"
	"flag"
	"fmt"

	"sunflower-gin/internal/conf"
	"sunflower-gin/internal/dao"
	"sunflower-gin/internal/server"
	"sunflower-gin/internal/task"
	"sunflower-gin/pkg/jwt"
	"sunflower-gin/pkg/logging"
	"sunflower-gin/pkg/snowflake"
)

var confPath = flag.String("conf", "./config/config.yaml", "配置文件路径")

func main() {
	// 加载配置
	flag.Parse()
	cfg := conf.Load(*confPath)

	// 初始化日志
	logger, err := logging.NewLogger(cfg)
	if err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer logger.Sync()

	dao.MustInitMySQL(cfg)  // 初始化 MySQL 连接
	dao.MustInitRedis(cfg)  // 初始化 Redis
	jwt.MustInit(cfg)       // 初始化 jwt
	snowflake.MustInit(cfg) // 初始化 snowflake

	// 初始化路由
	r := server.SetupRoutes(cfg)

	// 定时任务
	c := task.MustInit(context.Background())
	defer c.Stop()

	// 启动服务
	err = r.Run(fmt.Sprintf(":%d", cfg.GetInt("server.port")))
	if err != nil {
		panic(err)
	}
}
