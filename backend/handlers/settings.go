package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"backend/database"
	"backend/models"
)

// 密码修改请求结构
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword 修改用户密码
func ChangePassword(c *gin.Context) {
	// 在demo分支中，禁用密码修改功能
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success: false,
		Message: "密码修改功能已禁用",
	})
	return

	// 以下是原始代码，在demo分支中不会执行
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "未授权的访问",
		})
		return
	}

	var req PasswordChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 获取用户信息
	var storedPassword string
	err := database.DB.QueryRow("SELECT password FROM users WHERE id = ?", userID).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "用户不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "获取用户信息失败",
		})
		return
	}

	// 验证当前密码
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.CurrentPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "当前密码不正确",
		})
		return
	}

	// 生成新密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "密码哈希生成失败",
		})
		return
	}

	// 更新数据库中的密码
	_, err = database.DB.Exec("UPDATE users SET password = ? WHERE id = ?", string(hashedPassword), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "更新密码失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "密码修改成功",
	})
}
