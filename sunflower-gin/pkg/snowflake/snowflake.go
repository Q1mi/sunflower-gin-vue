package snowflake

import (
	"fmt"
	"time"

	"github.com/sony/sonyflake/v2"
	"github.com/spf13/viper"
)

var node *sonyflake.Sonyflake

// MustInit 初始化 snowflake
func MustInit(viper *viper.Viper) {
	// 完成 *sonyflake.Sonyflake 的初始化
	// 1. 读取配置文件中的起始时间
	st, err := time.Parse(time.DateOnly, viper.GetString("snowflake.start_time"))
	if err != nil {
		panic(fmt.Errorf("parse start time failed, err:%w", err))
	}
	settings := sonyflake.Settings{
		StartTime: st,
		MachineID: func() (int, error) {
			return viper.GetInt("snowflake.machine_id"), nil
		},
		CheckMachineID: func(int) bool { return true },
	}
	node, err = sonyflake.New(settings)
	if err != nil {
		panic(fmt.Errorf("init sonyflake failed, err:%w", err))
	}
}

func NextID() (int64, error) {
	if node == nil {
		return 0, fmt.Errorf("node is nil")
	}
	return node.NextID()
}
