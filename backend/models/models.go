package models

import "time"

// Profile 个人信息模型
type Profile struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Title         string    `json:"title"`
	Avatar        string    `json:"avatar"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Location      string    `json:"location"`
	Introduction  string    `json:"introduction"`
	YearsOfExp    int       `json:"years_of_exp"`
	Education     string    `json:"education"`
	JobStatus     string    `json:"job_status"`
	Philosophy    string    `json:"philosophy"`
	LastUpdated   time.Time `json:"last_updated"`
	ResumeFileURL string    `json:"resume_file_url"`
}

// SkillCategory 技能分类模型
type SkillCategory struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	Skills      []Skill `json:"skills,omitempty"`
}

// Skill 技能模型
type Skill struct {
	ID          int      `json:"id"`
	CategoryID  int      `json:"category_id"`
	Name        string   `json:"name"`
	Level       int      `json:"level"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// Experience 工作经历模型
type Experience struct {
	ID               int      `json:"id"`
	Period           string   `json:"period"`
	Title            string   `json:"title"`
	Company          string   `json:"company"`
	Location         string   `json:"location"`
	Color            string   `json:"color"`
	Icon             string   `json:"icon"`
	Responsibilities []string `json:"responsibilities"`
	Achievements     []string `json:"achievements"`
	Technologies     []string `json:"technologies"`
	SortOrder        int      `json:"sort_order"`
}

// Project 项目经验模型
type Project struct {
	ID               int      `json:"id"`
	Title            string   `json:"title"`
	Category         string   `json:"category"`
	Description      string   `json:"description"`
	Image            string   `json:"image"`
	DemoLink         string   `json:"demo_link"`
	RepoLink         string   `json:"repo_link"`
	ShowArchitecture bool     `json:"show_architecture"`
	Metrics          []Metric `json:"metrics"`
	KeyPoints        []string `json:"key_points"`
	TechStack        []string `json:"tech_stack"`
	SortOrder        int      `json:"sort_order"`
}

// Metric 项目成果指标模型
type Metric struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// Certificate 证书认证模型
type Certificate struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Date         string `json:"date"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	Link         string `json:"link"`
	SortOrder    int    `json:"sort_order"`
}

// User 用户模型(管理员)
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // 不输出到JSON
	Role     string `json:"role"`
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录响应
type LoginResponse struct {
	Token string `json:"token"`
}

// APIResponse 通用API响应
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// VerificationRequest 访客验证请求
type VerificationRequest struct {
	VerificationType string `json:"verification_type"` // "name", "email", "phone", "password"
	Value            string `json:"value"`
}

// VerificationResponse 验证响应
type VerificationResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
}
