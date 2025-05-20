package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// 初始化数据库
func SetupDatabase() error {
	// 确保数据库目录存在
	dbDir := "./data"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err = os.MkdirAll(dbDir, 0755)
		if err != nil {
			return err
		}
	}

	dbPath := filepath.Join(dbDir, "resume.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	DB = db

	// 创建表
	err = createTables()
	if err != nil {
		return err
	}

	// 检查是否需要初始化数据
	empty, err := isDatabaseEmpty()
	if err != nil {
		return err
	}

	if empty {
		err = initializeData()
		if err != nil {
			return err
		}
	}

	return nil
}

// 创建数据库表
func createTables() error {
	// 个人信息表
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS profile (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		title TEXT NOT NULL,
		avatar TEXT,
		email TEXT,
		phone TEXT,
		location TEXT,
		introduction TEXT,
		years_of_exp INTEGER,
		education TEXT,
		job_status TEXT,
		philosophy TEXT,
		last_updated TIMESTAMP,
		resume_file_url TEXT
	)`)
	if err != nil {
		return err
	}

	// 技能分类表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS skill_categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		icon TEXT
	)`)
	if err != nil {
		return err
	}

	// 技能表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS skills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id INTEGER,
		name TEXT NOT NULL,
		level INTEGER,
		description TEXT,
		tags TEXT,
		FOREIGN KEY (category_id) REFERENCES skill_categories(id)
	)`)
	if err != nil {
		return err
	}

	// 工作经历表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS experiences (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		period TEXT NOT NULL,
		title TEXT NOT NULL,
		company TEXT NOT NULL,
		location TEXT,
		color TEXT,
		icon TEXT,
		responsibilities TEXT,
		achievements TEXT,
		technologies TEXT,
		sort_order INTEGER
	)`)
	if err != nil {
		return err
	}

	// 项目经验表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		category TEXT,
		description TEXT,
		image TEXT,
		demo_link TEXT,
		repo_link TEXT,
		show_architecture BOOLEAN,
		metrics TEXT,
		key_points TEXT,
		tech_stack TEXT,
		sort_order INTEGER
	)`)
	if err != nil {
		return err
	}

	// 证书认证表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS certificates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		organization TEXT,
		date TEXT,
		description TEXT,
		icon TEXT,
		link TEXT,
		sort_order INTEGER
	)`)
	if err != nil {
		return err
	}

	// 用户表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT
	)`)
	if err != nil {
		return err
	}

	// 访客密码表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS visitor_access (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		access_key VARCHAR(255) NOT NULL,
		access_type VARCHAR(50) NOT NULL,
		value VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(access_type, value)
	);`)
	if err != nil {
		log.Printf("创建访客密码表失败: %v", err)
		return err
	}

	log.Printf("所有数据库表已创建")
	return nil
}

// 检查数据库是否为空
func isDatabaseEmpty() (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM profile").Scan(&count)
	if err != nil {
		// 如果表不存在或其他错误，认为是空的
		return true, nil
	}
	return count == 0, nil
}

// 初始化默认数据
func initializeData() error {
	// 初始化管理员用户
	_, err := DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)",
		"admin", "$2a$10$PYsUkJ2LUeAsbMtD3WkqNeETS4LjLxvg02ER2vKkBp0VIafe0nQTO", "admin") // 默认密码：admin123
	if err != nil {
		return err
	}

	// 初始化默认访客密码
	var visitorCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM visitor_access").Scan(&visitorCount)
	if err != nil {
		log.Printf("查询访客密码数量失败: %v", err)

		// 尝试创建表
		_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS visitor_access (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			access_key VARCHAR(255) NOT NULL,
			access_type VARCHAR(50) NOT NULL,
			value VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(access_type, value)
		);`)
		if err != nil {
			log.Printf("尝试再次创建访客密码表失败: %v", err)
			return err
		}

		visitorCount = 0
		log.Printf("重新创建了访客密码表")
	}

	log.Printf("当前访客密码记录数: %d", visitorCount)

	if visitorCount == 0 {
		_, err = DB.Exec(`
		INSERT INTO visitor_access (access_type, value, access_key) 
		VALUES ('password', 'default_password', 'visitor_key');`)
		if err != nil {
			log.Printf("创建默认访客密码失败: %v", err)
			return err
		}
		log.Printf("成功创建默认访客密码")
	}

	// 初始化个人信息
	_, err = DB.Exec(`
		INSERT INTO profile 
		(name, title, avatar, email, phone, location, introduction, years_of_exp, education, job_status, philosophy, last_updated)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"张三", "云原生运维工程师", "https://via.placeholder.com/150x150", "zhangsan@example.com", "138-1234-5678", "北京",
		"我是一名拥有5年经验的云原生运维工程师，专注于Kubernetes集群管理、自动化部署和系统架构设计。我热衷于利用最新技术解决复杂的运维挑战，并且善于构建高可用、可扩展的基础设施。",
		5,
		"计算机科学学士",
		"可入职",
		"基础设施即代码，自动化一切可自动化的事物",
		time.Now())
	if err != nil {
		return err
	}

	// 初始化技能分类
	categories := []map[string]interface{}{
		{
			"name":        "容器编排与云原生",
			"description": "在容器化和云原生技术领域的专业技能",
			"icon":        "fas fa-dharmachakra",
			"skills": []map[string]interface{}{
				{
					"name":        "Kubernetes",
					"level":       95,
					"description": "精通K8s集群管理、部署、扩展和故障排查",
					"tags":        []string{"集群管理", "HELM", "Operator", "微服务"},
				},
				{
					"name":        "Docker",
					"level":       90,
					"description": "熟练使用Docker进行容器化应用开发和部署",
					"tags":        []string{"镜像优化", "Docker Compose", "安全实践"},
				},
				{
					"name":        "Service Mesh",
					"level":       80,
					"description": "熟悉Istio等服务网格技术的部署和配置",
					"tags":        []string{"Istio", "流量管理", "服务发现"},
				},
				{
					"name":        "Prometheus & Grafana",
					"level":       85,
					"description": "熟练配置监控系统和构建可视化仪表板",
					"tags":        []string{"监控告警", "指标收集", "可视化"},
				},
			},
		},
		{
			"name":        "系统运维",
			"description": "在Linux/UNIX系统管理和自动化运维方面的能力",
			"icon":        "fas fa-server",
			"skills": []map[string]interface{}{
				{
					"name":        "Linux系统管理",
					"level":       90,
					"description": "深入理解Linux系统架构和系统调优",
					"tags":        []string{"故障诊断", "性能优化", "安全加固"},
				},
				{
					"name":        "Shell脚本",
					"level":       85,
					"description": "编写高效的Shell脚本实现自动化运维任务",
					"tags":        []string{"Bash", "自动化脚本", "任务调度"},
				},
				{
					"name":        "网络配置",
					"level":       80,
					"description": "理解网络架构和配置，实现高可用网络环境",
					"tags":        []string{"TCP/IP", "防火墙", "负载均衡"},
				},
				{
					"name":        "数据库运维",
					"level":       75,
					"description": "管理和维护MySQL、PostgreSQL等数据库系统",
					"tags":        []string{"备份恢复", "性能调优", "高可用配置"},
				},
			},
		},
	}

	for _, category := range categories {
		// 插入分类
		res, err := DB.Exec(
			"INSERT INTO skill_categories (name, description, icon) VALUES (?, ?, ?)",
			category["name"], category["description"], category["icon"])
		if err != nil {
			return err
		}

		categoryID, err := res.LastInsertId()
		if err != nil {
			return err
		}

		// 插入该分类下的技能
		for _, skill := range category["skills"].([]map[string]interface{}) {
			tags, err := json.Marshal(skill["tags"])
			if err != nil {
				return err
			}

			_, err = DB.Exec(
				"INSERT INTO skills (category_id, name, level, description, tags) VALUES (?, ?, ?, ?, ?)",
				categoryID, skill["name"], skill["level"], skill["description"], string(tags))
			if err != nil {
				return err
			}
		}
	}

	// 初始化工作经历
	experiences := []map[string]interface{}{
		{
			"period":   "2021 - 至今",
			"title":    "高级云原生运维工程师",
			"company":  "ABC云科技有限公司",
			"location": "北京",
			"color":    "#06B6D4",
			"icon":     "fas fa-cloud",
			"responsibilities": []string{
				"负责公司核心业务的Kubernetes集群规划、部署和日常运维",
				"设计并实现多环境（开发、测试、生产）的CI/CD自动化部署流水线",
				"针对微服务架构实现服务网格（Istio）和监控告警体系",
				"带领团队完成从传统架构向云原生架构的迁移工作",
			},
			"achievements": []string{
				"将应用部署时间从数小时缩短至15分钟，提升部署效率90%",
				"通过优化K8s集群资源配置，降低云资源成本约30%",
				"设计的监控告警系统提前发现95%的潜在问题，显著降低系统宕机风险",
			},
			"technologies": []string{"Kubernetes", "Docker", "Istio", "Jenkins", "Prometheus", "Terraform", "Helm"},
			"sort_order":   1,
		},
		{
			"period":   "2019 - 2021",
			"title":    "DevOps工程师",
			"company":  "智联云服务有限公司",
			"location": "上海",
			"color":    "#1E3A8A",
			"icon":     "fas fa-sync-alt",
			"responsibilities": []string{
				"负责构建和维护公司的CI/CD流水线和自动化测试框架",
				"实施基础设施即代码(IaC)实践，使用Terraform管理云资源",
				"优化容器化应用的构建、部署和运行流程",
				"与开发团队协作，提供技术支持和性能优化建议",
			},
			"achievements": []string{
				"实现了全自动化的部署流程，减少了90%的人工干预",
				"通过容器化改造，将资源利用率提升40%",
				"编写的自动化脚本库被公司多个团队采用，成为内部最佳实践",
			},
			"technologies": []string{"GitLab CI", "Docker", "Ansible", "AWS", "Terraform", "Python", "Shell"},
			"sort_order":   2,
		},
	}

	for _, exp := range experiences {
		responsibilities, err := json.Marshal(exp["responsibilities"])
		if err != nil {
			return err
		}

		achievements, err := json.Marshal(exp["achievements"])
		if err != nil {
			return err
		}

		technologies, err := json.Marshal(exp["technologies"])
		if err != nil {
			return err
		}

		_, err = DB.Exec(
			`INSERT INTO experiences 
			(period, title, company, location, color, icon, responsibilities, achievements, technologies, sort_order) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			exp["period"], exp["title"], exp["company"], exp["location"], exp["color"], exp["icon"],
			string(responsibilities), string(achievements), string(technologies), exp["sort_order"])
		if err != nil {
			return err
		}
	}

	// 初始化项目经验
	projects := []map[string]interface{}{
		{
			"title": "企业级微服务容器化迁移", "category": "云原生迁移", "description": "将传统单体应用拆分为微服务架构，并迁移至Kubernetes平台，提升系统弹性和可扩展性。", "image": "https://via.placeholder.com/800x600",
			"demo_link":         "#",
			"repo_link":         "#",
			"show_architecture": true,
			"metrics": []map[string]string{
				{"value": "99.99%", "label": "系统可用性"},
				{"value": "80%", "label": "资源利用率提升"},
				{"value": "60%", "label": "运维成本降低"},
			},
			"key_points": []string{
				"设计了微服务架构和容器化策略",
				"使用Helm构建了可复用的K8s部署模板",
				"实现了完整的CI/CD流水线自动部署",
				"配置了服务网格实现微服务治理",
			},
			"tech_stack": []string{"Kubernetes", "Docker", "Istio", "Helm", "Jenkins", "Prometheus", "Grafana"},
			"sort_order": 1,
		},
		{
			"title": "多云环境自动化运维平台", "category": "自动化部署", "description": "构建了一套跨多云平台的自动化运维系统，实现资源统一管理和应用自动部署。", "image": "https://via.placeholder.com/800x600",
			"demo_link":         "#",
			"repo_link":         "#",
			"show_architecture": false,
			"metrics": []map[string]string{
				{"value": "90%", "label": "部署效率提升"},
				{"value": "30%", "label": "多云成本优化"},
				{"value": "100+", "label": "自动化任务"},
			},
			"key_points": []string{
				"基于Terraform实现了多云资源统一管理",
				"开发了自动化部署工具与现有系统集成",
				"构建了中央化配置管理系统",
				"实现了多环境部署策略和回滚机制",
			},
			"tech_stack": []string{"Terraform", "Ansible", "Python", "GitLab CI", "AWS", "阿里云", "Shell"},
			"sort_order": 2,
		},
	}

	for _, project := range projects {
		metrics, err := json.Marshal(project["metrics"])
		if err != nil {
			return err
		}

		keyPoints, err := json.Marshal(project["key_points"])
		if err != nil {
			return err
		}

		techStack, err := json.Marshal(project["tech_stack"])
		if err != nil {
			return err
		}

		_, err = DB.Exec(
			`INSERT INTO projects 
			(title, category, description, image, demo_link, repo_link, show_architecture, metrics, key_points, tech_stack, sort_order) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			project["title"], project["category"], project["description"], project["image"],
			project["demo_link"], project["repo_link"], project["show_architecture"],
			string(metrics), string(keyPoints), string(techStack), project["sort_order"])
		if err != nil {
			return err
		}
	}

	// 初始化证书
	certificates := []map[string]string{
		{
			"name":         "Certified Kubernetes Administrator (CKA)",
			"organization": "Cloud Native Computing Foundation",
			"date":         "2022年3月",
			"description":  "Kubernetes管理员认证，证明了在部署和管理生产级Kubernetes集群方面的专业知识。",
			"icon":         "fas fa-dharmachakra",
			"link":         "#",
		},
		{
			"name":         "AWS Certified Solutions Architect",
			"organization": "Amazon Web Services",
			"date":         "2021年8月",
			"description":  "AWS解决方案架构师认证，验证了在AWS平台上设计和部署可扩展、高可用、容错系统的能力。",
			"icon":         "fab fa-aws",
			"link":         "#",
		},
	}

	for i, cert := range certificates {
		_, err = DB.Exec(
			`INSERT INTO certificates 
			(name, organization, date, description, icon, link, sort_order) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			cert["name"], cert["organization"], cert["date"], cert["description"],
			cert["icon"], cert["link"], i+1)
		if err != nil {
			return err
		}
	}

	log.Println("数据库初始化完成")
	return nil
}
