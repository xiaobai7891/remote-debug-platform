package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"remote-debug-platform/config"
	"remote-debug-platform/models"
	"remote-debug-platform/services"
	"remote-debug-platform/utils"
)

// GetAdminStats 获取系统统计
func GetAdminStats(c *gin.Context) {
	// 用户数
	var userCount int
	models.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&userCount)

	// 端口数
	var portCount int
	models.DB.QueryRow(`SELECT COUNT(*) FROM ports`).Scan(&portCount)

	// 活跃端口数
	var activePortCount int
	models.DB.QueryRow(`SELECT COUNT(*) FROM ports WHERE status = 'active'`).Scan(&activePortCount)

	// 在线设备数
	onlineDeviceCount, _ := models.GetAllOnlineDeviceCount()

	// 收入统计
	totalRecharge, _ := models.GetTotalRecharge()
	totalConsume, _ := models.GetTotalConsume()
	todayRecharge, _ := models.GetTodayRecharge()
	todayConsume, _ := models.GetTodayConsume()

	utils.Success(c, gin.H{
		"user_count":          userCount,
		"port_count":          portCount,
		"active_port_count":   activePortCount,
		"online_device_count": onlineDeviceCount,
		"total_recharge":      totalRecharge,
		"total_consume":       totalConsume,
		"today_recharge":      todayRecharge,
		"today_consume":       todayConsume,
	})
}

// GetAdminUsers 获取用户列表
func GetAdminUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var users []*models.User
	var total int
	var err error

	if keyword != "" {
		users, total, err = models.SearchUsers(keyword, page, pageSize)
	} else {
		users, total, err = models.GetAllUsers(page, pageSize)
	}

	if err != nil {
		utils.InternalError(c, "获取用户列表失败")
		return
	}

	utils.SuccessWithPage(c, users, total, page, pageSize)
}

type RechargeRequest struct {
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description"`
}

// RechargeUser 用户充值
func RechargeUser(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	description := req.Description
	if description == "" {
		description = "管理员充值"
	}

	if err := user.AddBalance(req.Amount, description); err != nil {
		utils.InternalError(c, "充值失败")
		return
	}

	// 重新获取用户信息
	user, _ = models.GetUserByID(userID)
	utils.Success(c, user)
}

type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active disabled"`
}

// UpdateUserStatus 更新用户状态
func UpdateUserStatus(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	user.Status = req.Status
	if err := user.Update(); err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.Success(c, user)
}

type UpdateUserRequest struct {
	Email      string  `json:"email"`
	MaxPorts   int     `json:"max_ports"`
	MaxDevices int     `json:"max_devices"`
	Role       string  `json:"role"`
	Balance    float64 `json:"balance"`
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.MaxPorts > 0 {
		user.MaxPorts = req.MaxPorts
	}
	if req.MaxDevices > 0 {
		user.MaxDevices = req.MaxDevices
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := user.Update(); err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.Success(c, user)
}

// GetAdminPorts 获取全部端口
func GetAdminPorts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	ports, total, err := models.GetAllPortsWithUser(page, pageSize)
	if err != nil {
		utils.InternalError(c, "获取端口列表失败")
		return
	}

	// 附加在线设备数
	result := make([]gin.H, 0)
	for _, p := range ports {
		devices, _ := models.GetOnlineDevices(p.ID)
		result = append(result, gin.H{
			"id":           p.ID,
			"user_id":      p.UserID,
			"username":     p.Username,
			"port":         p.Port,
			"protocol":     p.Protocol,
			"name":         p.Name,
			"token":        p.Token,
			"status":       p.Status,
			"expire_at":    p.ExpireAt,
			"created_at":   p.CreatedAt,
			"device_count": len(devices),
		})
	}

	utils.SuccessWithPage(c, result, total, page, pageSize)
}

// GetAdminDevices 获取全部设备
func GetAdminDevices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	devices, total, err := models.GetAllDevicesWithPort(page, pageSize, status)
	if err != nil {
		utils.InternalError(c, "获取设备列表失败")
		return
	}

	utils.SuccessWithPage(c, devices, total, page, pageSize)
}

// GetAdminOrders 获取全部订单
func GetAdminOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	orderType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, total, err := models.GetAllOrdersWithUser(page, pageSize, orderType)
	if err != nil {
		utils.InternalError(c, "获取订单列表失败")
		return
	}

	utils.SuccessWithPage(c, orders, total, page, pageSize)
}

// GetAdminConfig 获取系统配置
func GetAdminConfig(c *gin.Context) {
	cfg := config.Get()
	utils.Success(c, gin.H{
		"port_pool": gin.H{
			"min": cfg.PortPool.Min,
			"max": cfg.PortPool.Max,
		},
		"pricing": gin.H{
			"port_per_day": cfg.Pricing.PortPerDay,
		},
	})
}

type UpdateConfigRequest struct {
	PortPoolMin  int     `json:"port_pool_min"`
	PortPoolMax  int     `json:"port_pool_max"`
	PortPerDay   float64 `json:"port_per_day"`
}

// UpdateAdminConfig 更新系统配置
func UpdateAdminConfig(c *gin.Context) {
	var req UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	cfg := config.Get()
	if req.PortPoolMin > 0 {
		cfg.PortPool.Min = req.PortPoolMin
	}
	if req.PortPoolMax > 0 {
		cfg.PortPool.Max = req.PortPoolMax
	}
	if req.PortPerDay > 0 {
		cfg.Pricing.PortPerDay = req.PortPerDay
	}

	utils.Success(c, gin.H{
		"message": "配置已更新（内存中）",
	})
}

// ForceStopPort 强制停止端口
func ForceStopPort(c *gin.Context) {
	portID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	port, err := models.GetPortByID(portID)
	if err != nil {
		utils.NotFound(c, "端口不存在")
		return
	}

	// 停止端口监听
	services.GetPortManager().StopPortListener(port.Port)

	// 更新状态
	port.Status = "disabled"
	port.Update()

	// 设备下线
	models.SetPortDevicesOffline(portID)

	utils.Success(c, nil)
}

// ForceStartPort 强制启动端口
func ForceStartPort(c *gin.Context) {
	portID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	port, err := models.GetPortByID(portID)
	if err != nil {
		utils.NotFound(c, "端口不存在")
		return
	}

	// 启动端口监听
	services.GetPortManager().StartPortListener(port)

	// 更新状态
	port.Status = "active"
	port.Update()

	utils.Success(c, nil)
}
