package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"remote-debug-platform/config"
	"remote-debug-platform/handlers"
	"remote-debug-platform/middleware"
	"remote-debug-platform/models"
	"remote-debug-platform/services"
	"remote-debug-platform/utils"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Printf("加载配置文件失败: %v, 使用默认配置", err)
		config.SetDefault()
		cfg = config.Get()
	}

	// 初始化数据库
	if err := models.InitDB(cfg.Database.Path); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer models.CloseDB()

	// 初始化管理员账号
	initAdmin()

	// 初始化端口管理器
	pm := services.GetPortManager()
	if err := pm.InitFromDB(); err != nil {
		log.Printf("恢复端口监听失败: %v", err)
	}

	// 启动定时任务
	scheduler := services.NewScheduler()
	scheduler.Start()
	defer scheduler.Stop()

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())

	// API 路由
	api := r.Group("/api")
	{
		// 认证相关 - 无需登录
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
		}

		// 用户相关 - 需要登录
		user := api.Group("/user")
		user.Use(middleware.JWTAuth())
		{
			user.GET("/profile", handlers.GetProfile)
			user.PUT("/profile", handlers.UpdateProfile)
			user.PUT("/password", handlers.UpdatePassword)
			user.GET("/orders", handlers.GetOrders)
			user.GET("/dashboard", handlers.GetDashboard)
		}

		// 端口管理 - 需要登录
		ports := api.Group("/ports")
		ports.Use(middleware.JWTAuth())
		{
			ports.GET("", handlers.GetPorts)
			ports.POST("", handlers.CreatePort)
			ports.GET("/:id", handlers.GetPort)
			ports.PUT("/:id", handlers.UpdatePort)
			ports.DELETE("/:id", handlers.DeletePort)
			ports.POST("/:id/renew", handlers.RenewPort)
			ports.POST("/:id/change", handlers.ChangePort)
			ports.POST("/:id/regenerate-token", handlers.RegenerateToken)
		}

		// 设备管理 - 需要登录
		devices := api.Group("/devices")
		devices.Use(middleware.JWTAuth())
		{
			devices.GET("", handlers.GetDevices)
			devices.GET("/:id", handlers.GetDevice)
			devices.GET("/:id/logs", handlers.GetDeviceLogs)
			devices.POST("/:id/exec", handlers.ExecCode)
			devices.POST("/:id/screenshot", handlers.Screenshot)
		}

		// 管理后台 - 需要管理员权限
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
		{
			admin.GET("/stats", handlers.GetAdminStats)
			admin.GET("/users", handlers.GetAdminUsers)
			admin.POST("/users/:id/recharge", handlers.RechargeUser)
			admin.PUT("/users/:id/status", handlers.UpdateUserStatus)
			admin.PUT("/users/:id", handlers.UpdateUser)
			admin.GET("/ports", handlers.GetAdminPorts)
			admin.POST("/ports/:id/stop", handlers.ForceStopPort)
			admin.POST("/ports/:id/start", handlers.ForceStartPort)
			admin.GET("/devices", handlers.GetAdminDevices)
			admin.GET("/orders", handlers.GetAdminOrders)
			admin.GET("/config", handlers.GetAdminConfig)
			admin.PUT("/config", handlers.UpdateAdminConfig)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		utils.Success(c, gin.H{"status": "ok"})
	})

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("正在关闭服务...")
		scheduler.Stop()
		models.SetAllDevicesOffline()
		models.CloseDB()
		os.Exit(0)
	}()

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.HTTPPort)
	log.Printf("服务启动在 http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}

// initAdmin 初始化管理员账号
func initAdmin() {
	cfg := config.Get()

	// 检查管理员是否存在
	exists, _ := models.UserExists(cfg.Admin.Username)
	if exists {
		return
	}

	// 创建管理员
	hashedPassword, _ := utils.HashPassword(cfg.Admin.Password)
	user, err := models.CreateUser(cfg.Admin.Username, hashedPassword, "admin@localhost")
	if err != nil {
		log.Printf("创建管理员账号失败: %v", err)
		return
	}

	// 设置为管理员
	user.Role = "admin"
	user.Balance = 9999
	user.MaxPorts = 100
	user.MaxDevices = 1000
	user.Update()

	log.Printf("管理员账号已创建: %s / %s", cfg.Admin.Username, cfg.Admin.Password)
}
