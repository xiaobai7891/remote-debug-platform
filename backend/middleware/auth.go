package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"remote-debug-platform/models"
	"remote-debug-platform/utils"
)

const (
	ContextUserIDKey   = "user_id"
	ContextUsernameKey = "username"
	ContextRoleKey     = "role"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "token 无效或已过期")
			c.Abort()
			return
		}

		// 检查用户状态
		user, err := models.GetUserByID(claims.UserID)
		if err != nil {
			utils.Unauthorized(c, "用户不存在")
			c.Abort()
			return
		}

		if user.Status != "active" {
			utils.Forbidden(c, "账号已被禁用")
			c.Abort()
			return
		}

		// 设置上下文
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Set(ContextRoleKey, claims.Role)

		c.Next()
	}
}

// AdminAuth 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextRoleKey)
		if !exists || role != "admin" {
			utils.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get(ContextUserIDKey)
	if !exists {
		return 0
	}
	return userID.(int64)
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	username, exists := c.Get(ContextUsernameKey)
	if !exists {
		return ""
	}
	return username.(string)
}

// GetRole 从上下文获取角色
func GetRole(c *gin.Context) string {
	role, exists := c.Get(ContextRoleKey)
	if !exists {
		return ""
	}
	return role.(string)
}
