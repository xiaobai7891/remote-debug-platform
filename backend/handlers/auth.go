package handlers

import (
	"github.com/gin-gonic/gin"
	"remote-debug-platform/models"
	"remote-debug-platform/utils"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查用户名是否存在
	exists, err := models.UserExists(req.Username)
	if err != nil {
		utils.InternalError(c, "系统错误")
		return
	}
	if exists {
		utils.Error(c, 1001, "用户名已存在")
		return
	}

	// 检查邮箱是否存在
	if req.Email != "" {
		exists, err = models.EmailExists(req.Email)
		if err != nil {
			utils.InternalError(c, "系统错误")
			return
		}
		if exists {
			utils.Error(c, 1002, "邮箱已被使用")
			return
		}
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalError(c, "密码加密失败")
		return
	}

	// 创建用户
	user, err := models.CreateUser(req.Username, hashedPassword, req.Email)
	if err != nil {
		utils.InternalError(c, "创建用户失败: "+err.Error())
		return
	}

	// 生成 token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.InternalError(c, "生成 token 失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 查找用户
	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		utils.Error(c, 1003, "用户名或密码错误")
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		utils.Error(c, 1003, "用户名或密码错误")
		return
	}

	// 检查状态
	if user.Status != "active" {
		utils.Forbidden(c, "账号已被禁用")
		return
	}

	// 生成 token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.InternalError(c, "生成 token 失败")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	// JWT 是无状态的，客户端删除 token 即可
	utils.Success(c, nil)
}
