package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"backend/database"
	"backend/models"
)

// 生成访客令牌的密钥
var visitorSecretKey = []byte("visitor_secret_key")

// VerifyVisitor 验证访客身份
func VerifyVisitor(c *gin.Context) {
	var req models.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 验证类型有效性
	if req.VerificationType != "name" && req.VerificationType != "email" &&
		req.VerificationType != "phone" && req.VerificationType != "password" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的验证类型",
		})
		return
	}

	// 检查值是否为空
	if req.Value == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "验证值不能为空",
		})
		return
	}

	// 根据验证类型进行验证
	switch req.VerificationType {
	case "password":
		// 验证访客密码
		var accessKey string
		err := database.DB.QueryRow(
			"SELECT access_key FROM visitor_access WHERE access_type = ? AND value = ?",
			req.VerificationType, req.Value,
		).Scan(&accessKey)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, models.APIResponse{
					Success: false,
					Message: "密码验证失败",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "验证过程中发生错误",
			})
			return
		}

		// 生成访客令牌
		token, err := generateVisitorToken(accessKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "生成访客令牌失败",
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Message: "密码验证成功",
			Data: models.VerificationResponse{
				Success: true,
				Token:   token,
			},
		})

	case "name", "email", "phone":
		// 从个人信息中验证
		var value string
		var field string

		switch req.VerificationType {
		case "name":
			field = "name"
		case "email":
			field = "email"
		case "phone":
			field = "phone"
		}

		// 构建查询
		query := "SELECT " + field + " FROM profile WHERE " + field + " = ? LIMIT 1"
		err := database.DB.QueryRow(query, req.Value).Scan(&value)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, models.APIResponse{
					Success: false,
					Message: "验证失败，未找到匹配信息",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "验证过程中发生错误",
			})
			return
		}

		// 生成访客令牌
		token, err := generateVisitorToken(req.VerificationType + "_" + value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "生成访客令牌失败",
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Message: "验证成功",
			Data: models.VerificationResponse{
				Success: true,
				Token:   token,
			},
		})

	default:
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "不支持的验证类型",
		})
	}
}

// 生成访客JWT令牌
func generateVisitorToken(accessKey string) (string, error) {
	// 创建令牌
	token := jwt.New(jwt.SigningMethodHS256)

	// 设置声明
	claims := token.Claims.(jwt.MapClaims)
	claims["access_key"] = accessKey
	claims["type"] = "visitor"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 24小时有效期

	// 签名令牌
	tokenString, err := token.SignedString(visitorSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VisitorAuthMiddleware 访客身份验证中间件
func VisitorAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "未授权访问",
			})
			return
		}

		// 移除"Bearer "前缀
		tokenStr := authHeader[7:]

		// 解析令牌
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return visitorSecretKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "无效的访客令牌",
			})
			return
		}

		// 验证令牌有效性
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["type"] != "visitor" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, models.APIResponse{
					Success: false,
					Message: "无效的访客令牌类型",
				})
				return
			}
			// 设置访客标识
			c.Set("visitorAccessKey", claims["access_key"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "无效的访客令牌",
			})
			return
		}
	}
}

// ManageVisitorAccess 管理访客密码
func ManageVisitorAccess(c *gin.Context) {
	// 验证管理员权限
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		log.Printf("尝试访问管理功能但角色不是admin: %v", role)
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "需要管理员权限",
		})
		return
	}

	log.Printf("管理员准备获取访客密码列表，用户角色: %v", role)

	// 获取所有访问记录
	rows, err := database.DB.Query("SELECT id, access_type, value, access_key, created_at FROM visitor_access")
	if err != nil {
		log.Printf("查询访客密码数据库错误: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "获取访客密码记录失败",
		})
		return
	}
	defer rows.Close()

	var accessList []map[string]interface{}
	for rows.Next() {
		var id int
		var accessType, value, accessKey string
		var createdAt time.Time

		err := rows.Scan(&id, &accessType, &value, &accessKey, &createdAt)
		if err != nil {
			log.Printf("扫描访客密码行数据错误: %v", err)
			continue
		}

		accessList = append(accessList, map[string]interface{}{
			"id":          id,
			"access_type": accessType,
			"value":       value,
			"access_key":  accessKey,
			"created_at":  createdAt,
		})
	}

	log.Printf("成功获取访客密码列表，共 %d 条记录", len(accessList))
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "获取访客密码记录成功",
		Data:    accessList,
	})
}

// AddVisitorAccess 添加访客密码
func AddVisitorAccess(c *gin.Context) {
	// 在demo分支中，禁用访客密码添加功能
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success: false,
		Message: "访客密码添加功能已禁用",
	})
	return

	// 以下是原始代码，在demo分支中不会执行
	// 验证管理员权限
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		log.Printf("尝试访问管理功能但角色不是admin: %v", role)
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "需要管理员权限",
		})
		return
	}

	log.Printf("管理员准备添加访客密码")

	var accessData struct {
		AccessType string `json:"access_type" binding:"required"`
		Value      string `json:"value" binding:"required"`
		AccessKey  string `json:"access_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&accessData); err != nil {
		log.Printf("解析访客密码请求数据错误: %v", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 验证类型有效性
	if accessData.AccessType != "password" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "目前仅支持密码类型的访客验证",
		})
		return
	}

	// 插入数据
	_, err := database.DB.Exec(
		"INSERT INTO visitor_access (access_type, value, access_key) VALUES (?, ?, ?)",
		accessData.AccessType, accessData.Value, accessData.AccessKey,
	)

	if err != nil {
		log.Printf("添加访客密码到数据库失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "添加访客密码失败: " + err.Error(),
		})
		return
	}

	log.Printf("成功添加访客密码: %s", accessData.Value)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "访客密码添加成功",
	})
}

// DeleteVisitorAccess 删除访客密码
func DeleteVisitorAccess(c *gin.Context) {
	// 在demo分支中，禁用访客密码删除功能
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success: false,
		Message: "访客密码删除功能已禁用",
	})
	return

	// 以下是原始代码，在demo分支中不会执行
	// 验证管理员权限
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		log.Printf("尝试访问管理功能但角色不是admin: %v", role)
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "需要管理员权限",
		})
		return
	}

	idStr := c.Param("id")
	log.Printf("管理员准备删除访客密码，ID: %s", idStr)

	// 删除记录
	result, err := database.DB.Exec("DELETE FROM visitor_access WHERE id = ?", idStr)
	if err != nil {
		log.Printf("删除访客密码数据库操作失败: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "删除访客密码失败: " + err.Error(),
		})
		return
	}

	// 检查是否删除了记录
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("未找到ID为 %s 的访客密码记录", idStr)
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "未找到指定的访客密码记录",
		})
		return
	}

	log.Printf("成功删除ID为 %s 的访客密码", idStr)
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "访客密码删除成功",
	})
}
