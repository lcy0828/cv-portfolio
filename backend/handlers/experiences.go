package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/database"
	"backend/models"
)

// GetExperiences 获取所有工作经历
func GetExperiences(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, period, title, company, location, color, icon, 
		responsibilities, achievements, technologies, sort_order
		FROM experiences
		ORDER BY sort_order, id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "获取工作经历失败: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	experiences := []models.Experience{}
	for rows.Next() {
		var exp models.Experience
		var responsibilitiesJSON, achievementsJSON, technologiesJSON string

		if err := rows.Scan(
			&exp.ID, &exp.Period, &exp.Title, &exp.Company, &exp.Location,
			&exp.Color, &exp.Icon, &responsibilitiesJSON, &achievementsJSON,
			&technologiesJSON, &exp.SortOrder); err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "解析工作经历数据失败: " + err.Error(),
			})
			return
		}

		// 解析JSON数据
		if responsibilitiesJSON != "" {
			if err := json.Unmarshal([]byte(responsibilitiesJSON), &exp.Responsibilities); err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "解析工作职责数据失败: " + err.Error(),
				})
				return
			}
		}

		if achievementsJSON != "" {
			if err := json.Unmarshal([]byte(achievementsJSON), &exp.Achievements); err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "解析工作成就数据失败: " + err.Error(),
				})
				return
			}
		}

		if technologiesJSON != "" {
			if err := json.Unmarshal([]byte(technologiesJSON), &exp.Technologies); err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "解析使用技术数据失败: " + err.Error(),
				})
				return
			}
		}

		experiences = append(experiences, exp)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "获取工作经历成功",
		Data:    experiences,
	})
}

// GetExperience 获取单个工作经历
func GetExperience(c *gin.Context) {
	id := c.Param("id")
	expID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的工作经历ID",
		})
		return
	}

	var exp models.Experience
	var responsibilitiesJSON, achievementsJSON, technologiesJSON string

	err = database.DB.QueryRow(`
		SELECT id, period, title, company, location, color, icon, 
		responsibilities, achievements, technologies, sort_order
		FROM experiences
		WHERE id = ?`, expID).Scan(
		&exp.ID, &exp.Period, &exp.Title, &exp.Company, &exp.Location,
		&exp.Color, &exp.Icon, &responsibilitiesJSON, &achievementsJSON,
		&technologiesJSON, &exp.SortOrder)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "未找到指定工作经历",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "获取工作经历数据失败: " + err.Error(),
			})
		}
		return
	}

	// 解析JSON数据
	if responsibilitiesJSON != "" {
		if err := json.Unmarshal([]byte(responsibilitiesJSON), &exp.Responsibilities); err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "解析工作职责数据失败: " + err.Error(),
			})
			return
		}
	}

	if achievementsJSON != "" {
		if err := json.Unmarshal([]byte(achievementsJSON), &exp.Achievements); err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "解析工作成就数据失败: " + err.Error(),
			})
			return
		}
	}

	if technologiesJSON != "" {
		if err := json.Unmarshal([]byte(technologiesJSON), &exp.Technologies); err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "解析使用技术数据失败: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "获取工作经历成功",
		Data:    exp,
	})
}

// CreateExperience 创建工作经历
func CreateExperience(c *gin.Context) {
	// 在demo分支中，禁用工作经历创建功能
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success: false,
		Message: "工作经历创建功能已禁用",
	})
	return

	// 以下是原始代码，在demo分支中不会执行
	var exp models.Experience
	if err := c.ShouldBindJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 将切片数据转换为JSON
	responsibilitiesJSON, err := json.Marshal(exp.Responsibilities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "转换工作职责数据失败: " + err.Error(),
		})
		return
	}

	achievementsJSON, err := json.Marshal(exp.Achievements)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "转换工作成就数据失败: " + err.Error(),
		})
		return
	}

	technologiesJSON, err := json.Marshal(exp.Technologies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "转换使用技术数据失败: " + err.Error(),
		})
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO experiences 
		(period, title, company, location, color, icon, responsibilities, achievements, technologies, sort_order)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		exp.Period, exp.Title, exp.Company, exp.Location, exp.Color, exp.Icon,
		string(responsibilitiesJSON), string(achievementsJSON), string(technologiesJSON), exp.SortOrder)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "创建工作经历失败: " + err.Error(),
		})
		return
	}

	// 获取新创建的记录ID
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "获取新工作经历ID失败: " + err.Error(),
		})
		return
	}
	exp.ID = int(id)

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "工作经历创建成功",
		Data:    exp,
	})
}

// UpdateExperience 更新工作经历
func UpdateExperience(c *gin.Context) {
	// 在demo分支中，禁用工作经历更新功能
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success: false,
		Message: "工作经历更新功能已禁用",
	})
	return

	// 以下是原始代码，在demo分支中不会执行
	id := c.Param("id")
	expID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的工作经历ID",
		})
		return
	}

	var exp models.Experience
	if err := c.ShouldBindJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 将切片数据转换为JSON
	responsibilitiesJSON, err := json.Marshal(exp.Responsibilities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "转换工作职责数据失败: " + err.Error(),
		})
		return
	}

	achievementsJSON, err := json.Marshal(exp.Achievements)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "转换工作成就数据失败: " + err.Error(),
		})
		return
	}

	technologiesJSON, err := json.Marshal(exp.Technologies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "转换使用技术数据失败: " + err.Error(),
		})
		return
	}

	_, err = database.DB.Exec(`
		UPDATE experiences SET 
		period = ?, title = ?, company = ?, location = ?, color = ?, icon = ?, 
		responsibilities = ?, achievements = ?, technologies = ?, sort_order = ?
		WHERE id = ?`,
		exp.Period, exp.Title, exp.Company, exp.Location, exp.Color, exp.Icon,
		string(responsibilitiesJSON), string(achievementsJSON), string(technologiesJSON),
		exp.SortOrder, expID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "更新工作经历失败: " + err.Error(),
		})
		return
	}

	// 设置记录ID
	exp.ID = expID

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "工作经历更新成功",
		Data:    exp,
	})
}

// DeleteExperience 删除工作经历
func DeleteExperience(c *gin.Context) {
	// 在demo分支中，禁用工作经历删除功能
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success: false,
		Message: "工作经历删除功能已禁用",
	})
	return

	// 以下是原始代码，在demo分支中不会执行
	id := c.Param("id")
	expID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的工作经历ID",
		})
		return
	}

	result, err := database.DB.Exec("DELETE FROM experiences WHERE id = ?", expID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "删除工作经历失败: " + err.Error(),
		})
		return
	}

	// 检查是否删除了记录
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "未找到指定工作经历",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "工作经历删除成功",
	})
}
