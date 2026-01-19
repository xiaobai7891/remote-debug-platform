package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"remote-debug-platform/middleware"
	"remote-debug-platform/models"
	"remote-debug-platform/services"
	"remote-debug-platform/utils"
)

// GetDevices 获取设备列表
func GetDevices(c *gin.Context) {
	userID := middleware.GetUserID(c)

	devices, err := models.GetUserDevices(userID)
	if err != nil {
		utils.InternalError(c, "获取设备列表失败")
		return
	}

	// 附加端口信息
	result := make([]gin.H, 0)
	for _, d := range devices {
		port, _ := models.GetPortByID(d.PortID)
		portNum := 0
		if port != nil {
			portNum = port.Port
		}
		result = append(result, gin.H{
			"id":              d.ID,
			"device_id":       d.DeviceID,
			"device_name":     d.DeviceName,
			"device_info":     d.DeviceInfo,
			"last_ping":       d.LastPing,
			"status":          d.Status,
			"connected_at":    d.ConnectedAt,
			"disconnected_at": d.DisconnectedAt,
			"port":            portNum,
			"port_id":         d.PortID,
		})
	}

	utils.Success(c, result)
}

// GetDevice 获取设备详情
func GetDevice(c *gin.Context) {
	userID := middleware.GetUserID(c)
	deviceID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	device, err := models.GetDeviceByID(deviceID)
	if err != nil {
		utils.NotFound(c, "设备不存在")
		return
	}

	// 检查权限
	if device.UserID != userID && middleware.GetRole(c) != "admin" {
		utils.Forbidden(c, "无权访问")
		return
	}

	// 获取端口信息
	port, _ := models.GetPortByID(device.PortID)

	utils.Success(c, gin.H{
		"device": device,
		"port":   port,
	})
}

// GetDeviceLogs 获取设备日志
func GetDeviceLogs(c *gin.Context) {
	userID := middleware.GetUserID(c)
	deviceID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	device, err := models.GetDeviceByID(deviceID)
	if err != nil {
		utils.NotFound(c, "设备不存在")
		return
	}

	// 检查权限
	if device.UserID != userID && middleware.GetRole(c) != "admin" {
		utils.Forbidden(c, "无权访问")
		return
	}

	logs, err := models.GetDeviceLogs(deviceID, limit)
	if err != nil {
		utils.InternalError(c, "获取日志失败")
		return
	}

	utils.Success(c, logs)
}

type ExecRequest struct {
	Code string `json:"code" binding:"required"`
}

// ExecCode 远程执行代码
func ExecCode(c *gin.Context) {
	userID := middleware.GetUserID(c)
	deviceID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req ExecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	device, err := models.GetDeviceByID(deviceID)
	if err != nil {
		utils.NotFound(c, "设备不存在")
		return
	}

	// 检查权限
	if device.UserID != userID && middleware.GetRole(c) != "admin" {
		utils.Forbidden(c, "无权操作")
		return
	}

	// 检查设备在线状态
	if device.Status != "online" {
		utils.Error(c, 3001, "设备不在线")
		return
	}

	// 获取端口信息
	port, err := models.GetPortByID(device.PortID)
	if err != nil {
		utils.InternalError(c, "端口不存在")
		return
	}

	// 发送执行命令
	err = services.GetPortManager().SendToDevice(port.Port, device.DeviceID, gin.H{
		"type": "exec",
		"code": req.Code,
	})
	if err != nil {
		utils.Error(c, 3002, "发送命令失败: "+err.Error())
		return
	}

	// 记录日志
	models.AddDeviceLog(deviceID, "info", "执行代码: "+req.Code[:min(len(req.Code), 100)])

	utils.Success(c, nil)
}

// Screenshot 远程截图
func Screenshot(c *gin.Context) {
	userID := middleware.GetUserID(c)
	deviceID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	device, err := models.GetDeviceByID(deviceID)
	if err != nil {
		utils.NotFound(c, "设备不存在")
		return
	}

	// 检查权限
	if device.UserID != userID && middleware.GetRole(c) != "admin" {
		utils.Forbidden(c, "无权操作")
		return
	}

	// 检查设备在线状态
	if device.Status != "online" {
		utils.Error(c, 3001, "设备不在线")
		return
	}

	// 获取端口信息
	port, err := models.GetPortByID(device.PortID)
	if err != nil {
		utils.InternalError(c, "端口不存在")
		return
	}

	// 发送截图命令
	err = services.GetPortManager().SendToDevice(port.Port, device.DeviceID, gin.H{
		"type": "screenshot",
	})
	if err != nil {
		utils.Error(c, 3002, "发送命令失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"message": "截图命令已发送，请稍后查看",
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
