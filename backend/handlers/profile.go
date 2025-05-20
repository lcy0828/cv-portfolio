package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"backend/database"
	"backend/models"
)

// GetProfile 获取个人信息
func GetProfile(c *gin.Context) {
	var profile models.Profile
	var resumeFileURL sql.NullString // 使用sql.NullString处理可能为NULL的字段

	err := database.DB.QueryRow(`
		SELECT id, name, title, avatar, email, phone, location, introduction, 
		years_of_exp, education, job_status, philosophy, last_updated, resume_file_url 
		FROM profile LIMIT 1`).Scan(
		&profile.ID, &profile.Name, &profile.Title, &profile.Avatar,
		&profile.Email, &profile.Phone, &profile.Location, &profile.Introduction,
		&profile.YearsOfExp, &profile.Education, &profile.JobStatus, &profile.Philosophy,
		&profile.LastUpdated, &resumeFileURL)

	if err != nil {
		if err == sql.ErrNoRows {
			// 没有找到记录，返回空的profile对象
			c.JSON(http.StatusOK, models.APIResponse{
				Success: true,
				Message: "获取个人信息成功，但没有数据",
				Data:    models.Profile{},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "获取个人信息失败: " + err.Error(),
		})
		return
	}

	// 处理可能为NULL的字段
	if resumeFileURL.Valid {
		profile.ResumeFileURL = resumeFileURL.String
	} else {
		profile.ResumeFileURL = ""
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "获取个人信息成功",
		Data:    profile,
	})
}

// UpdateProfile 更新个人信息
func UpdateProfile(c *gin.Context) {
	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 更新最后修改时间
	profile.LastUpdated = time.Now()

	// 检查是否存在记录
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM profile").Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "查询失败: " + err.Error(),
		})
		return
	}

	var result sql.Result
	if count == 0 {
		// 没有记录，执行插入操作
		result, err = database.DB.Exec(`
			INSERT INTO profile 
			(name, title, avatar, email, phone, location, introduction, 
			years_of_exp, education, job_status, philosophy, last_updated, resume_file_url)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			profile.Name, profile.Title, profile.Avatar, profile.Email,
			profile.Phone, profile.Location, profile.Introduction, profile.YearsOfExp,
			profile.Education, profile.JobStatus, profile.Philosophy,
			profile.LastUpdated, profile.ResumeFileURL)
	} else {
		// 有记录，执行更新操作
		result, err = database.DB.Exec(`
			UPDATE profile SET 
			name = ?, title = ?, avatar = ?, email = ?, phone = ?, location = ?, 
			introduction = ?, years_of_exp = ?, education = ?, job_status = ?, 
			philosophy = ?, last_updated = ?, resume_file_url = ?
			WHERE id = ?`,
			profile.Name, profile.Title, profile.Avatar, profile.Email,
			profile.Phone, profile.Location, profile.Introduction, profile.YearsOfExp,
			profile.Education, profile.JobStatus, profile.Philosophy,
			profile.LastUpdated, profile.ResumeFileURL, profile.ID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "更新个人信息失败: " + err.Error(),
		})
		return
	}

	if count == 0 {
		// 获取新插入记录的ID
		id, _ := result.LastInsertId()
		profile.ID = int(id)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "个人信息更新成功",
		Data:    profile,
	})
}
