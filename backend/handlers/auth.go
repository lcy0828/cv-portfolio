package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"backend/database"
	"backend/models"
)

// JWT密钥，实际应用中应从环境变量或配置文件读取
var jwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "your-secret-key"))

// 获取环境变量，不存在则使用默认值
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// 自定义JWT声明结构
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// Login 用户登录处理
func Login(c *gin.Context) {
	var loginReq models.LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的登录信息",
		})
		return
	}

	// 查询用户
	var user models.User
	var hashedPassword string
	err := database.DB.QueryRow("SELECT id, username, password, role FROM users WHERE username = ?", loginReq.Username).Scan(
		&user.ID, &user.Username, &hashedPassword, &user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "用户名或密码错误",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "登录处理失败",
		})
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginReq.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "用户名或密码错误",
		})
		return
	}

	// 生成JWT令牌
	expirationTime := time.Now().Add(24 * time.Hour) // 令牌有效期1天
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "resume-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "生成令牌失败",
		})
		return
	}

	// 返回登录成功和令牌
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "登录成功",
		Data: models.LoginResponse{
			Token: tokenString,
		},
	})
}

// AuthMiddleware 身份验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("管理员验证中间件处理请求: %s %s", c.Request.Method, c.Request.URL.Path)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("请求缺少Authorization头")
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "需要身份验证",
			})
			c.Abort()
			return
		}

		// Bearer Token格式: "Bearer {token}"
		tokenString := authHeader[7:] // 跳过"Bearer "前缀

		// 解析令牌
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			log.Printf("令牌解析错误: %v", err)
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "无效的令牌",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// 将用户信息存入上下文
			log.Printf("令牌验证成功，用户角色: %s", claims.Role)
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
			c.Next()
		} else {
			log.Printf("令牌有效性检查失败")
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "无效的令牌",
			})
			c.Abort()
		}
	}
}
