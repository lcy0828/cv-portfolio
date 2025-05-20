package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"backend/database"
	"backend/handlers"
)

func main() {
	// 初始化数据库
	err := initDatabase()
	if err != nil {
		log.Printf("数据库初始化失败: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin路由
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 提供静态文件(Admin UI)
	r.StaticFS("/admin", http.Dir("./backend/admin"))

	// 提供前端静态文件
	r.StaticFS("/", http.Dir("./public"))

	// API路由组
	api := r.Group("/api")
	{
		// 访客验证接口 - 无需任何验证
		api.POST("/verify", handlers.VerifyVisitor)

		// 登录接口 - 无需任何验证
		api.POST("/login", handlers.Login)

		// 临时诊断接口 - 仅用于开发调试
		api.GET("/debug/visitor-access", func(c *gin.Context) {
			log.Printf("正在执行访客密码表诊断")

			// 检查表是否存在
			var exists int
			err := database.DB.QueryRow("SELECT 1 FROM sqlite_master WHERE type='table' AND name='visitor_access'").Scan(&exists)
			if err != nil {
				log.Printf("检查visitor_access表是否存在失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   err.Error(),
					"message": "检查表存在性失败",
				})
				return
			}

			if exists != 1 {
				log.Printf("visitor_access表不存在！")
				c.JSON(http.StatusOK, gin.H{
					"table_exists": false,
				})
				return
			}

			// 查询表中的数据
			rows, err := database.DB.Query("SELECT id, access_type, value, access_key, created_at FROM visitor_access")
			if err != nil {
				log.Printf("查询visitor_access表失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   err.Error(),
					"message": "查询表失败",
				})
				return
			}
			defer rows.Close()

			var records []map[string]interface{}
			for rows.Next() {
				var id int
				var accessType, value, accessKey string
				var createdAt time.Time

				err := rows.Scan(&id, &accessType, &value, &accessKey, &createdAt)
				if err != nil {
					log.Printf("扫描行数据失败: %v", err)
					continue
				}

				records = append(records, map[string]interface{}{
					"id":          id,
					"access_type": accessType,
					"value":       value,
					"access_key":  accessKey,
					"created_at":  createdAt,
				})
			}

			log.Printf("找到 %d 条访客密码记录", len(records))
			c.JSON(http.StatusOK, gin.H{
				"table_exists": true,
				"records":      records,
				"count":        len(records),
			})
		})

		// 临时诊断接口 - 重置访客密码表
		api.GET("/debug/reset-visitor-access", func(c *gin.Context) {
			log.Printf("重置访客密码表")

			// 删除表
			_, err := database.DB.Exec("DROP TABLE IF EXISTS visitor_access")
			if err != nil {
				log.Printf("删除visitor_access表失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   err.Error(),
					"success": false,
					"message": "删除表失败",
				})
				return
			}
			log.Printf("已删除visitor_access表")

			// 重新创建表
			_, err = database.DB.Exec(`
			CREATE TABLE visitor_access (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				access_key VARCHAR(255) NOT NULL,
				access_type VARCHAR(50) NOT NULL,
				value VARCHAR(255) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				UNIQUE(access_type, value)
			);`)
			if err != nil {
				log.Printf("重新创建visitor_access表失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   err.Error(),
					"success": false,
					"message": "重新创建表失败",
				})
				return
			}
			log.Printf("已重新创建visitor_access表")

			// 插入默认数据
			_, err = database.DB.Exec(`
			INSERT INTO visitor_access (access_type, value, access_key) 
			VALUES ('password', 'default_password', 'visitor_key');`)
			if err != nil {
				log.Printf("插入默认访客密码记录失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   err.Error(),
					"success": false,
					"message": "插入默认数据失败",
				})
				return
			}

			log.Printf("成功重置访客密码表并添加默认记录")
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "重置访客密码表成功",
			})
		})

		// 受保护的管理接口 - 仅需要管理员身份验证
		admin := api.Group("/admin")
		admin.Use(handlers.AuthMiddleware())
		{
			// 个人信息接口 - 修改需要认证
			admin.GET("/profile", handlers.GetProfile)
			admin.PUT("/profile", handlers.UpdateProfile)

			// 技能接口
			admin.GET("/skills", handlers.GetSkills)
			admin.GET("/skills/:id", handlers.GetSkill)
			admin.POST("/skills", handlers.CreateSkill)
			admin.PUT("/skills/:id", handlers.UpdateSkill)
			admin.DELETE("/skills/:id", handlers.DeleteSkill)

			// 工作经历接口
			admin.GET("/experiences", handlers.GetExperiences)
			admin.GET("/experiences/:id", handlers.GetExperience)
			admin.POST("/experiences", handlers.CreateExperience)
			admin.PUT("/experiences/:id", handlers.UpdateExperience)
			admin.DELETE("/experiences/:id", handlers.DeleteExperience)

			// 项目经验接口
			admin.GET("/projects", handlers.GetProjects)
			admin.GET("/projects/:id", handlers.GetProject)
			admin.POST("/projects", handlers.CreateProject)
			admin.PUT("/projects/:id", handlers.UpdateProject)
			admin.DELETE("/projects/:id", handlers.DeleteProject)

			// 证书接口
			admin.GET("/certificates", handlers.GetCertificates)
			admin.GET("/certificates/:id", handlers.GetCertificate)
			admin.POST("/certificates", handlers.CreateCertificate)
			admin.PUT("/certificates/:id", handlers.UpdateCertificate)
			admin.DELETE("/certificates/:id", handlers.DeleteCertificate)

			// 设置接口
			admin.PUT("/settings/password", handlers.ChangePassword)

			// 访客密码管理
			admin.GET("/visitor/access", handlers.ManageVisitorAccess)
			admin.POST("/visitor/access", handlers.AddVisitorAccess)
			admin.DELETE("/visitor/access/:id", handlers.DeleteVisitorAccess)
		}

		// 需要访客验证的接口 - 只提供GET请求访问
		visitor := api.Group("/")
		visitor.Use(handlers.VisitorAuthMiddleware())
		{
			// 个人信息接口 - 仅GET需要访客验证
			visitor.GET("/profile", handlers.GetProfile)

			// 技能接口 - 仅GET需要访客验证
			visitor.GET("/skills", handlers.GetSkills)
			visitor.GET("/skills/:id", handlers.GetSkill)

			// 工作经历接口 - 仅GET需要访客验证
			visitor.GET("/experiences", handlers.GetExperiences)
			visitor.GET("/experiences/:id", handlers.GetExperience)

			// 项目经验接口 - 仅GET需要访客验证
			visitor.GET("/projects", handlers.GetProjects)
			visitor.GET("/projects/:id", handlers.GetProject)

			// 证书接口 - 仅GET需要访客验证
			visitor.GET("/certificates", handlers.GetCertificates)
			visitor.GET("/certificates/:id", handlers.GetCertificate)
		}
	}

	// 获取端口，默认为8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务器
	log.Printf("后端API服务已启动，运行在 http://localhost:%s", port)
	r.Run(":" + port)
}

// 初始化数据库
func initDatabase() error {
	// 调用数据库初始化函数
	return database.SetupDatabase()
}
