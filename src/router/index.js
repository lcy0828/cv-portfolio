import { createRouter, createWebHistory } from 'vue-router';
import App from '../App.vue';
import Login from '../views/Login.vue';
import VisitorVerification from '../components/VisitorVerification.vue';

const routes = [
  {
    path: '/',
    name: 'Home',
    component: App,
    meta: { requiresVisitor: true }
  },
  {
    path: '/verify',
    name: 'VisitorVerification',
    component: VisitorVerification
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('../views/Admin.vue'),
    meta: { requiresAuth: true }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' };
    } else if (savedPosition) {
      return savedPosition;
    } else {
      return { top: 0 };
    }
  }
});

// 路由守卫，检查权限
router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // 检查是否已登录（管理员）
    const token = localStorage.getItem('token');
    if (!token) {
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      });
    } else {
      next();
    }
  } else if (to.matched.some(record => record.meta.requiresVisitor)) {
    // 检查是否已验证访客身份
    const visitorToken = localStorage.getItem('visitorToken');
    if (!visitorToken) {
      next({
        path: '/verify',
        query: { redirect: to.fullPath }
      });
    } else {
      next();
    }
  } else {
    next();
  }
});

export default router; 