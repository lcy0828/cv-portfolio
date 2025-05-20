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

// GetSkills 获取所有技能分类及技能
func GetSkills(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, name, description, icon
		FROM skill_categories
		ORDER BY id`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "获取技能分类失败: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	categories := []models.SkillCategory{}
	for rows.Next() {
		var category models.SkillCategory
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.Icon); err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "解析技能分类数据失败: " + err.Error(),
			})
			return
		}

		// 获取该分类下的所有技能
		skillRows, err := database.DB.Query(`
			SELECT id, name, level, description, tags
			FROM skills
			WHERE category_id = ?
			ORDER BY id`, category.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "获取技能数据失败: " + err.Error(),
			})
			return
		}
		defer skillRows.Close()

		category.Skills = []models.Skill{}
		for skillRows.Next() {
			var skill models.Skill
			var tagsJSON string
			if err := skillRows.Scan(&skill.ID, &skill.Name, &skill.Level, &skill.Description, &tagsJSON); err != nil {
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Success: false,
					Message: "解析技能数据失败: " + err.Error(),
				})
				return
			}

			// 解析JSON格式的标签
			var tags []string
			if tagsJSON != "" {
				if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
					c.JSON(http.StatusInternalServerError, models.APIResponse{
						Success: false,
						Message: "解析技能标签失败: " + err.Error(),
					})
					return
				}
			}
			skill.Tags = tags
			skill.CategoryID = category.ID
			category.Skills = append(category.Skills, skill)
		}

		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "获取技能数据成功",
		Data:    categories,
	})
}

// GetSkill 获取单个技能
func GetSkill(c *gin.Context) {
	id := c.Param("id")
	skillID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的技能ID",
		})
		return
	}

	var skill models.Skill
	var tagsJSON string
	err = database.DB.QueryRow(`
		SELECT id, category_id, name, level, description, tags
		FROM skills
		WHERE id = ?`, skillID).Scan(
		&skill.ID, &skill.CategoryID, &skill.Name, &skill.Level, &skill.Description, &tagsJSON)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "未找到指定技能",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "获取技能数据失败: " + err.Error(),
			})
		}
		return
	}

	// 解析JSON格式的标签
	var tags []string
	if tagsJSON != "" {
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "解析技能标签失败: " + err.Error(),
			})
			return
		}
	}
	skill.Tags = tags

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "获取技能数据成功",
		Data:    skill,
	})
}

// CreateSkill 创建技能
func CreateSkill(c *gin.Context) {
	// 在demo分支中，禁用技能创建功能
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success: false,
		Message: "技能创建功能已禁用",
	})
	return

	// 以下是原始代码，在demo分支中不会执行
	var skill models.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 将标签转换为JSON
	tagsJSON, err := json.Marshal(skill.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "转换标签数据失败: " + err.Error(),
		})
		return
	}

	result, err := database.DB.Exec(`
		INSERT INTO skills (category_id, name, level, description, tags)
		VALUES (?, ?, ?, ?, ?)`,
		skill.CategoryID, skill.Name, skill.Level, skill.Description, string(tagsJSON))

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "创建技能失败: " + err.Error(),
		})
		return
	}

	// 获取新创建的技能ID
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "获取新技能ID失败: " + err.Error(),
		})
		return
	}
	skill.ID = int(id)

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "技能创建成功",
		Data:    skill,
	})
}

// UpdateSkill 更新技能
func UpdateSkill(c *gin.Context) {
	// 在demo分支中，禁用技能更新功能
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success: false,
		Message: "技能更新功能已禁用",
	})
	return

	// 以下是原始代码，在demo分支中不会执行
	id := c.Param("id")
	skillID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的技能ID",
		})
		return
	}

	var skill models.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的请求数据: " + err.Error(),
		})
		return
	}

	// 将标签转换为JSON
	tagsJSON, err := json.Marshal(skill.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "转换标签数据失败: " + err.Error(),
		})
		return
	}

	_, err = database.DB.Exec(`
		UPDATE skills 
		SET category_id = ?, name = ?, level = ?, description = ?, tags = ?
		WHERE id = ?`,
		skill.CategoryID, skill.Name, skill.Level, skill.Description, string(tagsJSON), skillID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "更新技能失败: " + err.Error(),
		})
		return
	}

	// 设置技能ID
	skill.ID = skillID

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "技能更新成功",
		Data:    skill,
	})
}

// DeleteSkill 删除技能
func DeleteSkill(c *gin.Context) {
	// 在demo分支中，禁用技能删除功能
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success: false,
		Message: "技能删除功能已禁用",
	})
	return

	// 以下是原始代码，在demo分支中不会执行
	id := c.Param("id")
	skillID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的技能ID",
		})
		return
	}

	result, err := database.DB.Exec("DELETE FROM skills WHERE id = ?", skillID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "删除技能失败: " + err.Error(),
		})
		return
	}

	// 检查是否删除了记录
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "未找到指定技能",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "技能删除成功",
	})
}
