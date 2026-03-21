package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")     // 执行目录在根目录时
	v.AddConfigPath("../../config") // 执行目录在 cmd/server 时

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config failed:%w", err)
	}
	//fmt.Println("Viper 读取到的端口:", v.GetInt("app.port"))

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config failed:%w", err)
	}

	return &cfg, nil
}
