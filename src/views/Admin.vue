<template>
  <div class="admin-container">
    <header class="admin-header">
      <div class="logo">
        <h1>简历管理系统</h1>
      </div>
      <div class="user-actions">
        <button class="logout-btn" @click="handleLogout">
          <i class="fas fa-sign-out-alt"></i> 退出登录
        </button>
      </div>
    </header>
    
    <div class="admin-content">
      <aside class="admin-sidebar">
        <nav class="sidebar-nav">
          <div 
            v-for="(item, index) in menuItems" 
            :key="index"
            class="nav-item"
            :class="{ 'active': activeSection === item.id }"
            @click="setActiveSection(item.id)"
          >
            <i :class="item.icon"></i>
            <span>{{ item.name }}</span>
          </div>
        </nav>
      </aside>
      
      <main class="admin-main">
        <div v-if="activeSection === 'profile'" class="admin-section">
          <h2>个人信息管理</h2>
          <div class="section-content">
            <ProfileForm />
          </div>
        </div>
        
        <div v-else-if="activeSection === 'skills'" class="admin-section">
          <h2>技能管理</h2>
          <div class="section-content">
            <SkillsForm />
          </div>
        </div>
        
        <div v-else-if="activeSection === 'experiences'" class="admin-section">
          <h2>工作经历管理</h2>
          <div class="section-content">
            <ExperienceForm />
          </div>
        </div>
        
        <div v-else-if="activeSection === 'projects'" class="admin-section">
          <h2>项目经验管理</h2>
          <div class="section-content">
            <ProjectsForm />
          </div>
        </div>
        
        <div v-else-if="activeSection === 'certificates'" class="admin-section">
          <h2>证书管理</h2>
          <div class="section-content">
            <CertificatesForm />
          </div>
        </div>
        
        <div v-else-if="activeSection === 'settings'" class="admin-section">
          <h2>系统设置</h2>
          <div class="section-content">
            <SettingsForm />
          </div>
        </div>
        
        <div v-else-if="activeSection === 'visitor'" class="admin-section">
          <h2>访客密码管理</h2>
          <div class="section-content">
            <VisitorAccessForm />
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import ProfileForm from '../components/admin/ProfileForm.vue';
import SkillsForm from '../components/admin/SkillsForm.vue';
import ExperienceForm from '../components/admin/ExperienceForm.vue';
import ProjectsForm from '../components/admin/ProjectsForm.vue';
import CertificatesForm from '../components/admin/CertificatesForm.vue';
import SettingsForm from '../components/admin/SettingsForm.vue';
import VisitorAccessForm from '../components/admin/VisitorAccessForm.vue';

const router = useRouter();
const activeSection = ref('profile');

const menuItems = [
  { id: 'profile', name: '个人信息', icon: 'fas fa-user' },
  { id: 'skills', name: '技能管理', icon: 'fas fa-code' },
  { id: 'experiences', name: '工作经历', icon: 'fas fa-briefcase' },
  { id: 'projects', name: '项目经验', icon: 'fas fa-project-diagram' },
  { id: 'certificates', name: '证书管理', icon: 'fas fa-certificate' },
  { id: 'visitor', name: '访客密码', icon: 'fas fa-key' },
  { id: 'settings', name: '系统设置', icon: 'fas fa-cog' }
];

const setActiveSection = (sectionId) => {
  activeSection.value = sectionId;
};

const handleLogout = () => {
  localStorage.removeItem('token');
  router.push('/login');
};
</script>

<style scoped>
.admin-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.admin-header {
  background-color: var(--primary-color);
  color: white;
  padding: 15px 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo h1 {
  color: white;
  font-size: 1.5rem;
  margin: 0;
}

.logout-btn {
  background: none;
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  padding: 8px 15px;
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.logout-btn:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.admin-content {
  display: flex;
  flex: 1;
}

.admin-sidebar {
  width: 250px;
  background-color: #f1f5f9;
  padding: 30px 0;
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
}

.nav-item {
  padding: 15px 30px;
  display: flex;
  align-items: center;
  gap: 15px;
  cursor: pointer;
  transition: all 0.3s;
  border-left: 4px solid transparent;
}

.nav-item:hover {
  background-color: #e2e8f0;
}

.nav-item.active {
  background-color: #e2e8f0;
  border-left-color: var(--accent-color);
  color: var(--primary-color);
  font-weight: 600;
}

.nav-item i {
  font-size: 1.2rem;
  width: 20px;
  text-align: center;
}

.admin-main {
  flex: 1;
  padding: 30px;
  background-color: white;
}

.admin-section h2 {
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid #e2e8f0;
}

.section-content {
  background-color: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
}

@media (max-width: 768px) {
  .admin-content {
    flex-direction: column;
  }
  
  .admin-sidebar {
    width: 100%;
    padding: 15px 0;
  }
  
  .sidebar-nav {
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: center;
  }
  
  .nav-item {
    padding: 10px 15px;
    border-left: none;
    border-bottom: 2px solid transparent;
  }
  
  .nav-item.active {
    border-left-color: transparent;
    border-bottom-color: var(--accent-color);
  }
}
</style> 