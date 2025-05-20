package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/database"
	"backend/models"
)

// GetCertificates 获取所有证书
func GetCertificates(c *gin.Context) {
	rows, err := database.DB.Query(`SELECT id, name, organization, date, description, icon, link, sort_order 
	                               FROM certificates ORDER BY sort_order ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "获取证书列表失败",
		})
		return
	}
	defer rows.Close()

	var certificates []models.Certificate
	for rows.Next() {
		var cert models.Certificate

		err := rows.Scan(
			&cert.ID, &cert.Name, &cert.Organization, &cert.Date, &cert.Description,
			&cert.Icon, &cert.Link, &cert.SortOrder,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "读取证书数据失败",
			})
			return
		}

		certificates = append(certificates, cert)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "获取证书列表成功",
		Data:    certificates,
	})
}

// GetCertificate 获取单个证书详情
func GetCertificate(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的证书ID",
		})
		return
	}

	var cert models.Certificate
	err = database.DB.QueryRow(`SELECT id, name, organization, date, description, icon, link, sort_order 
	                           FROM certificates WHERE id = ?`, idInt).Scan(
		&cert.ID, &cert.Name, &cert.Organization, &cert.Date, &cert.Description,
		&cert.Icon, &cert.Link, &cert.SortOrder,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "证书不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "获取证书详情失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "获取证书详情成功",
		Data:    cert,
	})
}

// CreateCertificate 创建新证书
func CreateCertificate(c *gin.Context) {
	// 在demo分支中，禁用证书创建功能
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success: false,
		Message: "证书创建功能已禁用",
	})
	return

	// 以下是原始代码，在demo分支中不会执行
	var certificate models.Certificate
	if err := c.ShouldBindJSON(&certificate); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的证书数据: " + err.Error(),
		})
		return
	}

	result, err := database.DB.Exec(
		`INSERT INTO certificates 
		(name, organization, date, description, icon, link, sort_order) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		certificate.Name, certificate.Organization, certificate.Date, certificate.Description,
		certificate.Icon, certificate.Link, certificate.SortOrder)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "创建证书失败: " + err.Error(),
		})
		return
	}

	// 获取新插入记录的ID
	id, _ := result.LastInsertId()
	certificate.ID = int(id)

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "创建证书成功",
		Data:    certificate,
	})
}

// UpdateCertificate 更新证书信息
func UpdateCertificate(c *gin.Context) {
	// 在demo分支中，禁用证书更新功能
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success: false,
		Message: "证书更新功能已禁用",
	})
	return

	// 以下是原始代码，在demo分支中不会执行
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的证书ID",
		})
		return
	}

	var certificate models.Certificate
	if err := c.ShouldBindJSON(&certificate); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的证书数据: " + err.Error(),
		})
		return
	}

	// 检查记录是否存在
	var exists int
	err = database.DB.QueryRow("SELECT 1 FROM certificates WHERE id = ?", idInt).Scan(&exists)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "证书不存在",
		})
		return
	}

	_, err = database.DB.Exec(
		`UPDATE certificates SET 
		name=?, organization=?, date=?, description=?, icon=?, link=?, sort_order=? 
		WHERE id=?`,
		certificate.Name, certificate.Organization, certificate.Date, certificate.Description,
		certificate.Icon, certificate.Link, certificate.SortOrder, idInt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "更新证书失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "更新证书成功",
		Data:    certificate,
	})
}

// DeleteCertificate 删除证书
func DeleteCertificate(c *gin.Context) {
	// 在demo分支中，禁用证书删除功能
	c.JSON(http.StatusForbidden, models.APIResponse{
		Success: false,
		Message: "证书删除功能已禁用",
	})
	return

	// 以下是原始代码，在demo分支中不会执行
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "无效的证书ID",
		})
		return
	}

	// 检查记录是否存在
	var exists int
	err = database.DB.QueryRow("SELECT 1 FROM certificates WHERE id = ?", idInt).Scan(&exists)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "证书不存在",
		})
		return
	}

	_, err = database.DB.Exec("DELETE FROM certificates WHERE id = ?", idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "删除证书失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "删除证书成功",
	})
}
