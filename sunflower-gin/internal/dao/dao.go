package dao

import (
	"context"
	"fmt"
	"time"

	"sunflower-gin/internal/dao/query"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

// MustInitMySQL 初始化 MySQL 连接
func MustInitMySQL(cfg *viper.Viper) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.GetString("mysql.user"),
		cfg.GetString("mysql.password"),
		cfg.GetString("mysql.host"),
		cfg.GetString("mysql.port"),
		cfg.GetString("mysql.dbname"),
	)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("connect db fail: %w", err))
	}
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.GetInt("mysql.max_idle_conns"))
	sqlDB.SetMaxOpenConns(cfg.GetInt("mysql.max_open_conns"))
	sqlDB.SetConnMaxLifetime(cfg.GetDuration("mysql.max_lifetime"))

	query.SetDefault(db) // 指定 query 包使用的默认数据库连接
}

// MustInitRedis 初始化 Redis 连接
func MustInitRedis(conf *viper.Viper) {
	addr := fmt.Sprintf("%s:%d", conf.GetString("redis.host"), conf.GetInt("redis.port"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.GetString("redis.password"),
		DB:       conf.GetInt("redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Errorf("init redis failed, err:%w", err))
	}
	RedisClient = rdb
}
