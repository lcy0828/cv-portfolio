import axios from 'axios';

// 创建axios实例
const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

// 请求拦截器 - 可以在这里添加认证信息等
api.interceptors.request.use(
  config => {
    // 如果存在token，可以在这里添加到请求头
    const token = localStorage.getItem('token');
    const visitorToken = localStorage.getItem('visitorToken');
    
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    } else if (visitorToken) {
      // 使用访客令牌
      config.headers.Authorization = `Bearer ${visitorToken}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截器 - 处理常见错误
api.interceptors.response.use(
  response => {
    return response.data;
  },
  error => {
    if (error.response && error.response.status === 401) {
      // 未授权，可能需要重新登录或访客验证
      if (localStorage.getItem('token')) {
        // 管理员登录失效
        localStorage.removeItem('token');
        // 可以添加重定向到登录页面
      } else if (localStorage.getItem('visitorToken')) {
        // 访客令牌失效
        localStorage.removeItem('visitorToken');
        // 重定向到访客验证页面
        window.location.href = '/verify';
      }
    }
    return Promise.reject(error);
  }
);

// API服务对象
const apiService = {
  // 访客验证
  verifyVisitor(data) {
    return api.post('/verify', data);
  },
  
  // 访客密码管理
  getVisitorAccess() {
    return api.get('/admin/visitor/access');
  },
  addVisitorAccess(data) {
    return api.post('/admin/visitor/access', data);
  },
  deleteVisitorAccess(id) {
    return api.delete(`/admin/visitor/access/${id}`);
  },
  
  // 个人信息相关
  getProfile() {
    return api.get('/profile');
  },
  updateProfile(data) {
    return api.put('/admin/profile', data);
  },
  
  // 技能相关
  getSkills() {
    return api.get('/skills');
  },
  getSkill(id) {
    return api.get(`/skills/${id}`);
  },
  createSkill(data) {
    return api.post('/admin/skills', data);
  },
  updateSkill(id, data) {
    return api.put(`/admin/skills/${id}`, data);
  },
  deleteSkill(id) {
    return api.delete(`/admin/skills/${id}`);
  },
  
  // 工作经历相关
  getExperiences() {
    return api.get('/experiences');
  },
  getExperience(id) {
    return api.get(`/experiences/${id}`);
  },
  createExperience(data) {
    return api.post('/admin/experiences', data);
  },
  updateExperience(id, data) {
    return api.put(`/admin/experiences/${id}`, data);
  },
  deleteExperience(id) {
    return api.delete(`/admin/experiences/${id}`);
  },
  
  // 项目经验相关
  getProjects() {
    return api.get('/projects');
  },
  getProject(id) {
    return api.get(`/projects/${id}`);
  },
  createProject(data) {
    return api.post('/admin/projects', data);
  },
  updateProject(id, data) {
    return api.put(`/admin/projects/${id}`, data);
  },
  deleteProject(id) {
    return api.delete(`/admin/projects/${id}`);
  },
  
  // 证书相关
  getCertificates() {
    return api.get('/certificates');
  },
  getCertificate(id) {
    return api.get(`/certificates/${id}`);
  },
  createCertificate(data) {
    return api.post('/admin/certificates', data);
  },
  updateCertificate(id, data) {
    return api.put(`/admin/certificates/${id}`, data);
  },
  deleteCertificate(id) {
    return api.delete(`/admin/certificates/${id}`);
  },
  
  // 系统设置
  changePassword(data) {
    return api.put('/admin/settings/password', data);
  },
  
  // 登录
  login(credentials) {
    return api.post('/login', credentials);
  }
};

export default apiService; 