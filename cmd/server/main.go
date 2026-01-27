package main

import (
	"charonoms/internal/infrastructure/config"
	"charonoms/internal/infrastructure/logger"
	"charonoms/internal/infrastructure/persistence/mysql"
	"charonoms/internal/interfaces/http/router"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg, err := config.Load("./config/config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init(cfg.Logger); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("CharonOMS starting...")

	// 初始化数据库
	if err := mysql.Init(cfg.Database); err != nil {
		logger.Fatal("Failed to init database", zap.Error(err))
	}
	defer func() {
		if err := mysql.Close(); err != nil {
			logger.Error("Failed to close database", zap.Error(err))
		}
	}()

	// 设置路由
	r := router.SetupRouter(cfg)

	// 启动服务器
	addr := ":" + cfg.Server.Port
	logger.Info("Server started",
		zap.String("addr", addr),
		zap.String("mode", cfg.Server.Mode),
	)

	// 优雅关闭
	go func() {
		if err := r.Run(addr); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down...")
}
