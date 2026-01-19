package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"remote-debug-platform/middleware"
	"remote-debug-platform/models"
	"remote-debug-platform/utils"
)

type UpdateProfileRequest struct {
	Email string `json:"email"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// GetProfile 获取个人信息
func GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 获取统计信息
	stats, _ := models.GetUserStats(userID)

	utils.Success(c, gin.H{
		"user":  user,
		"stats": stats,
	})
}

// UpdateProfile 更新个人信息
func UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 检查邮箱是否被其他人使用
	if req.Email != "" && req.Email != user.Email {
		exists, _ := models.EmailExists(req.Email)
		if exists {
			utils.Error(c, 1002, "邮箱已被使用")
			return
		}
	}

	user.Email = req.Email
	if err := user.Update(); err != nil {
		utils.InternalError(c, "更新失败")
		return
	}

	utils.Success(c, user)
}

// UpdatePassword 修改密码
func UpdatePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		utils.Error(c, 1004, "原密码错误")
		return
	}

	// 哈希新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.InternalError(c, "密码加密失败")
		return
	}

	if err := user.UpdatePassword(hashedPassword); err != nil {
		utils.InternalError(c, "更新密码失败")
		return
	}

	utils.Success(c, nil)
}

// GetOrders 获取消费记录
func GetOrders(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, total, err := models.GetUserOrders(userID, page, pageSize)
	if err != nil {
		utils.InternalError(c, "获取订单失败")
		return
	}

	utils.SuccessWithPage(c, orders, total, page, pageSize)
}

// GetDashboard 获取用户仪表盘数据
func GetDashboard(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := models.GetUserByID(userID)
	if err != nil {
		utils.NotFound(c, "用户不存在")
		return
	}

	// 获取统计数据
	stats, _ := models.GetUserStats(userID)

	// 获取端口列表
	ports, _ := models.GetUserPorts(userID)

	// 获取在线设备
	devices, _ := models.GetUserDevices(userID)
	onlineDevices := make([]*models.Device, 0)
	for _, d := range devices {
		if d.Status == "online" {
			onlineDevices = append(onlineDevices, d)
		}
	}

	utils.Success(c, gin.H{
		"user":           user,
		"stats":          stats,
		"ports":          ports,
		"online_devices": onlineDevices,
	})
}
