package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"remote-debug-platform/config"
	"remote-debug-platform/middleware"
	"remote-debug-platform/models"
	"remote-debug-platform/services"
	"remote-debug-platform/utils"
)

type CreatePortRequest struct {
	Protocol string `json:"protocol"` // autojs, lazyman, easyclick
	Name     string `json:"name"`
	Days     int    `json:"days"` // 购买天数
}

type UpdatePortRequest struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
}

type RenewPortRequest struct {
	Days int `json:"days" binding:"required,min=1"`
}

// GetPorts 获取端口列表
func GetPorts(c *gin.Context) {
	userID := middleware.GetUserID(c)

	ports, err := models.GetUserPorts(userID)
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

	utils.Success(c, result)
}

// CreatePort 申请新端口
func CreatePort(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreatePortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 默认值
	if req.Protocol == "" {
		req.Protocol = "autojs"
	}
	if req.Days < 1 {
		req.Days = 30
	}

	// 检查用户端口数量限制
	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	portCount, _ := models.GetUserPortCount(userID)
	if portCount >= user.MaxPorts {
		utils.Error(c, 2001, "已达到端口数量上限")
		return
	}

	// 计算费用
	cfg := config.Get()
	cost := cfg.Pricing.PortPerDay * float64(req.Days)

	// 检查余额
	if user.Balance < cost {
		utils.Error(c, 2002, "余额不足")
		return
	}

	// 分配端口
	portNum, err := services.GetPortManager().AllocatePort()
	if err != nil {
		utils.Error(c, 2003, "无可用端口")
		return
	}

	// 创建端口记录
	port, err := models.CreatePort(userID, portNum, req.Protocol, req.Days)
	if err != nil {
		services.GetPortManager().ReleasePort(portNum)
		utils.InternalError(c, "创建端口失败")
		return
	}

	// 扣费
	err = user.AddBalance(-cost, "申请端口 "+strconv.Itoa(portNum)+" "+strconv.Itoa(req.Days)+"天")
	if err != nil {
		utils.InternalError(c, "扣费失败")
		return
	}

	// 更新端口名称
	if req.Name != "" {
		port.Name = req.Name
		port.Update()
	}

	// 启动端口监听
	services.GetPortManager().StartPortListener(port)

	utils.Success(c, port)
}

// GetPort 获取端口详情
func GetPort(c *gin.Context) {
	userID := middleware.GetUserID(c)
	portID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	port, err := models.GetPortByID(portID)
	if err != nil {
		utils.NotFound(c, "端口不存在")
		return
	}

	// 检查权限
	if port.UserID != userID && middleware.GetRole(c) != "admin" {
		utils.Forbidden(c, "无权访问")
		return
	}

	// 获取设备列表
	devices, _ := models.GetPortDevices(portID)

	utils.Success(c, gin.H{
		"port":    port,
		"devices": devices,
	})
}

// UpdatePort 更新端口配置
func UpdatePort(c *gin.Context) {
	userID := middleware.GetUserID(c)
	portID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req UpdatePortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	port, err := models.GetPortByID(portID)
	if err != nil {
		utils.NotFound(c, "端口不存在")
		return
	}

	// 检查权限
	if port.UserID != userID {
		utils.Forbidden(c, "无权操作")
		return
	}

	if req.Name != "" {
		port.Name = req.Name
	}
	if req.Protocol != "" {
		port.Protocol = req.Protocol
	}

	if err := port.Update(); err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.Success(c, port)
}

// DeletePort 释放端口
func DeletePort(c *gin.Context) {
	userID := middleware.GetUserID(c)
	portID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	port, err := models.GetPortByID(portID)
	if err != nil {
		utils.NotFound(c, "端口不存在")
		return
	}

	// 检查权限
	if port.UserID != userID && middleware.GetRole(c) != "admin" {
		utils.Forbidden(c, "无权操作")
		return
	}

	// 停止端口监听
	services.GetPortManager().StopPortListener(port.Port)

	// 删除端口
	if err := port.Delete(); err != nil {
		utils.InternalError(c, "删除失败")
		return
	}

	// 释放端口号
	services.GetPortManager().ReleasePort(port.Port)

	utils.Success(c, nil)
}

// RenewPort 续费端口
func RenewPort(c *gin.Context) {
	userID := middleware.GetUserID(c)
	portID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req RenewPortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	port, err := models.GetPortByID(portID)
	if err != nil {
		utils.NotFound(c, "端口不存在")
		return
	}

	// 检查权限
	if port.UserID != userID {
		utils.Forbidden(c, "无权操作")
		return
	}

	// 计算费用
	cfg := config.Get()
	cost := cfg.Pricing.PortPerDay * float64(req.Days)

	// 检查余额
	user, _ := models.GetUserByID(userID)
	if user.Balance < cost {
		utils.Error(c, 2002, "余额不足")
		return
	}

	// 扣费
	err = user.AddBalance(-cost, "续费端口 "+strconv.Itoa(port.Port)+" "+strconv.Itoa(req.Days)+"天")
	if err != nil {
		utils.InternalError(c, "扣费失败")
		return
	}

	// 续期
	if err := port.Renew(req.Days); err != nil {
		utils.InternalError(c, "续费失败")
		return
	}

	// 如果端口已过期停止，重新启动
	if port.Status == "expired" {
		services.GetPortManager().StartPortListener(port)
	}

	utils.Success(c, port)
}

// ChangePort 更换端口
func ChangePort(c *gin.Context) {
	userID := middleware.GetUserID(c)
	portID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	port, err := models.GetPortByID(portID)
	if err != nil {
		utils.NotFound(c, "端口不存在")
		return
	}

	// 检查权限
	if port.UserID != userID {
		utils.Forbidden(c, "无权操作")
		return
	}

	// 停止旧端口
	services.GetPortManager().StopPortListener(port.Port)
	oldPort := port.Port

	// 分配新端口
	newPortNum, err := services.GetPortManager().AllocatePort()
	if err != nil {
		// 重新启动旧端口
		services.GetPortManager().StartPortListener(port)
		utils.Error(c, 2003, "无可用端口")
		return
	}

	// 更新数据库
	_, err = models.DB.Exec(`UPDATE ports SET port = ? WHERE id = ?`, newPortNum, port.ID)
	if err != nil {
		services.GetPortManager().ReleasePort(newPortNum)
		services.GetPortManager().StartPortListener(port)
		utils.InternalError(c, "更换失败")
		return
	}

	// 释放旧端口
	services.GetPortManager().ReleasePort(oldPort)

	// 更新 port 对象
	port.Port = newPortNum

	// 启动新端口
	services.GetPortManager().StartPortListener(port)

	// 生成新 token
	port.RegenerateToken()

	utils.Success(c, port)
}

// RegenerateToken 重新生成 Token
func RegenerateToken(c *gin.Context) {
	userID := middleware.GetUserID(c)
	portID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	port, err := models.GetPortByID(portID)
	if err != nil {
		utils.NotFound(c, "端口不存在")
		return
	}

	// 检查权限
	if port.UserID != userID {
		utils.Forbidden(c, "无权操作")
		return
	}

	if err := port.RegenerateToken(); err != nil {
		utils.InternalError(c, "重新生成失败")
		return
	}

	// 重新获取
	port, _ = models.GetPortByID(portID)

	utils.Success(c, port)
}
