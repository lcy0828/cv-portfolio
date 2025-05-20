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

// GetProjects 获取所有项目经验
func GetProjects(c *gin.Context) {
	rows, err := database.DB.Query(`SELECT id, title, category, description, image, demo_link, repo_link, 
	                               show_architecture, metrics, key_points, tech_stack, sort_order 
	                               FROM projects ORDER BY sort_order ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "获取项目列表失败",
		})
		return
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		var metricsJSON, keyPointsJSON, techStackJSON string

		err := rows.Scan(
			&p.ID, &p.Title, &p.Category, &p.Description, &p.Image, &p.DemoLink, &p.RepoLink,
			&p.ShowArchitecture, &metricsJSON, &keyPointsJSON, &techStackJSON, &p.SortOrder,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "读取项目数据失败",
			})
			return
		}

		// 解析JSON字段
		var metrics []models.Metric
		if err := json.Unmarshal([]byte(metricsJSON), &metrics); err == nil {
			p.Metrics = metrics
		}

		var keyPoints []string
		if err := json.Unmarshal([]byte(keyPointsJSON), &keyPoints); err == nil {
			p.KeyPoints = keyPoints
		}

		var techStack []string
		if err := json.Unmarshal([]byte(techStackJSON), &techStack); err == nil {
			p.TechStack = techStack
		}

		projects = append(projects, p)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "获取项目列表成功",
		Data:    projects,
	})
}

// GetProject 获取单个项目详情
func GetProject(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的项目ID",
		})
		return
	}

	var p models.Project
	var metricsJSON, keyPointsJSON, techStackJSON string

	err = database.DB.QueryRow(`SELECT id, title, category, description, image, demo_link, repo_link, 
	                           show_architecture, metrics, key_points, tech_stack, sort_order 
	                           FROM projects WHERE id = ?`, idInt).Scan(
		&p.ID, &p.Title, &p.Category, &p.Description, &p.Image, &p.DemoLink, &p.RepoLink,
		&p.ShowArchitecture, &metricsJSON, &keyPointsJSON, &techStackJSON, &p.SortOrder,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "项目不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "获取项目详情失败",
		})
		return
	}

	// 解析JSON字段
	var metrics []models.Metric
	if err := json.Unmarshal([]byte(metricsJSON), &metrics); err == nil {
		p.Metrics = metrics
	}

	var keyPoints []string
	if err := json.Unmarshal([]byte(keyPointsJSON), &keyPoints); err == nil {
		p.KeyPoints = keyPoints
	}

	var techStack []string
	if err := json.Unmarshal([]byte(techStackJSON), &techStack); err == nil {
		p.TechStack = techStack
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "获取项目详情成功",
		Data:    p,
	})
}

// CreateProject 创建新项目
func CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的项目数据: " + err.Error(),
		})
		return
	}

	// 转换为JSON存储
	metricsJSON, _ := json.Marshal(project.Metrics)
	keyPointsJSON, _ := json.Marshal(project.KeyPoints)
	techStackJSON, _ := json.Marshal(project.TechStack)

	result, err := database.DB.Exec(
		`INSERT INTO projects 
		(title, category, description, image, demo_link, repo_link, show_architecture, metrics, key_points, tech_stack, sort_order) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		project.Title, project.Category, project.Description, project.Image,
		project.DemoLink, project.RepoLink, project.ShowArchitecture,
		string(metricsJSON), string(keyPointsJSON), string(techStackJSON), project.SortOrder)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "创建项目失败: " + err.Error(),
		})
		return
	}

	// 获取新插入记录的ID
	id, _ := result.LastInsertId()
	project.ID = int(id)

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "创建项目成功",
		Data:    project,
	})
}

// UpdateProject 更新项目信息
func UpdateProject(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的项目ID",
		})
		return
	}

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的项目数据: " + err.Error(),
		})
		return
	}

	// 检查记录是否存在
	var exists int
	err = database.DB.QueryRow("SELECT 1 FROM projects WHERE id = ?", idInt).Scan(&exists)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "项目不存在",
		})
		return
	}

	// 转换为JSON存储
	metricsJSON, _ := json.Marshal(project.Metrics)
	keyPointsJSON, _ := json.Marshal(project.KeyPoints)
	techStackJSON, _ := json.Marshal(project.TechStack)

	_, err = database.DB.Exec(
		`UPDATE projects SET 
		title=?, category=?, description=?, image=?, demo_link=?, repo_link=?, show_architecture=?, 
		metrics=?, key_points=?, tech_stack=?, sort_order=? 
		WHERE id=?`,
		project.Title, project.Category, project.Description, project.Image,
		project.DemoLink, project.RepoLink, project.ShowArchitecture,
		string(metricsJSON), string(keyPointsJSON), string(techStackJSON), project.SortOrder, idInt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "更新项目失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "更新项目成功",
		Data:    project,
	})
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的项目ID",
		})
		return
	}

	// 检查记录是否存在
	var exists int
	err = database.DB.QueryRow("SELECT 1 FROM projects WHERE id = ?", idInt).Scan(&exists)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "项目不存在",
		})
		return
	}

	_, err = database.DB.Exec("DELETE FROM projects WHERE id = ?", idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "删除项目失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "删除项目成功",
	})
}
