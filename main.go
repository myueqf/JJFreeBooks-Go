package main

import (
	"JJFreeBooks/config"
	"JJFreeBooks/utils"
)

func main() {
	logger := utils.GetLogger()
	appConfig, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		panic(err)
	}
	logger.Info("配置加载成功")
	logger.Info("Token", "", appConfig.Token)
	logger.Info("Cron", "", appConfig.Cron)
}
