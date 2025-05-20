<template>
  <header class="navbar" :class="{ 'navbar-scrolled': scrolled }">
    <div class="container navbar-container">
      <div class="logo">张三</div>
      <div class="menu-toggle" @click="toggleMenu">
        <i class="fas fa-bars"></i>
      </div>
      <nav :class="{ 'active': menuActive }">
        <ul>
          <li v-for="item in navItems" :key="item.id">
            <a :href="'#' + item.id" @click="closeMenu">{{ item.name }}</a>
          </li>
          <li class="github-link">
            <a href="https://github.com/lcy0828/cv-portfolio" target="_blank" title="GitHub仓库">
              <i class="fab fa-github"></i>
            </a>
          </li>
        </ul>
      </nav>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

const scrolled = ref(false);
const menuActive = ref(false);

const navItems = [
  { id: 'about', name: '关于我' },
  { id: 'skills', name: '专业技能' },
  { id: 'experience', name: '工作经历' },
  { id: 'projects', name: '项目经验' },
  { id: 'certificates', name: '证书认证' },
  { id: 'contact', name: '联系方式' }
];

const handleScroll = () => {
  scrolled.value = window.scrollY > 50;
};

const toggleMenu = () => {
  menuActive.value = !menuActive.value;
};

const closeMenu = () => {
  menuActive.value = false;
};

onMounted(() => {
  window.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});
</script>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 1000;
  padding: 20px 0;
  transition: all 0.5s;
  background-color: transparent;
}

.navbar-scrolled {
  padding: 12px 0;
  background-color: rgba(30, 58, 138, 0.95);
  box-shadow: 0 2px 15px rgba(0, 0, 0, 0.1);
}

.navbar-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  font-size: 1.8rem;
  font-weight: 700;
  color: var(--light-color);
  text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.3);
}

.navbar-scrolled .logo {
  color: var(--light-color);
}

nav ul {
  display: flex;
  list-style: none;
  align-items: center;
}

nav ul li {
  margin-left: 30px;
}

nav ul li a {
  color: var(--light-color);
  font-weight: 500;
  font-size: 1rem;
  padding: 8px 0;
  position: relative;
}

nav ul li a:after {
  content: '';
  position: absolute;
  width: 0;
  height: 2px;
  background: var(--accent-color);
  bottom: 0;
  left: 0;
  transition: width 0.3s;
}

nav ul li a:hover:after {
  width: 100%;
}

.github-link a {
  font-size: 1.2rem;
}

.github-link a:hover {
  color: var(--accent-color);
}

.github-link a:after {
  display: none;
}

.menu-toggle {
  display: none;
  color: var(--light-color);
  font-size: 1.5rem;
  cursor: pointer;
}

@media (max-width: 768px) {
  .menu-toggle {
    display: block;
  }

  nav {
    position: fixed;
    top: 70px;
    right: -100%;
    width: 250px;
    height: calc(100vh - 70px);
    background-color: var(--primary-color);
    transition: 0.5s;
    padding: 20px;
    z-index: 999;
  }

  nav.active {
    right: 0;
  }

  nav ul {
    display: block;
    margin-top: 20px;
  }

  nav ul li {
    margin: 0 0 20px 0;
  }

  nav ul li a {
    display: block;
    font-size: 1.1rem;
  }

  .github-link {
    margin-top: 20px;
  }

  .github-link a {
    display: inline-block;
    font-size: 1.5rem;
  }
}
</style> 