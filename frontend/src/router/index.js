import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue'),
    meta: { guest: true }
  },
  {
    path: '/',
    component: () => import('../views/Layout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue')
      },
      {
        path: 'ports',
        name: 'Ports',
        component: () => import('../views/Ports.vue')
      },
      {
        path: 'devices',
        name: 'Devices',
        component: () => import('../views/Devices.vue')
      },
      {
        path: 'orders',
        name: 'Orders',
        component: () => import('../views/Orders.vue')
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('../views/Profile.vue')
      }
    ]
  },
  {
    path: '/admin',
    component: () => import('../views/admin/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: () => import('../views/admin/Dashboard.vue')
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('../views/admin/Users.vue')
      },
      {
        path: 'ports',
        name: 'AdminPorts',
        component: () => import('../views/admin/Ports.vue')
      },
      {
        path: 'devices',
        name: 'AdminDevices',
        component: () => import('../views/admin/Devices.vue')
      },
      {
        path: 'orders',
        name: 'AdminOrders',
        component: () => import('../views/admin/Orders.vue')
      },
      {
        path: 'config',
        name: 'AdminConfig',
        component: () => import('../views/admin/Config.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
  } else if (to.meta.guest && userStore.isLoggedIn) {
    next('/')
  } else if (to.meta.requiresAdmin && userStore.user?.role !== 'admin') {
    next('/')
  } else {
    next()
  }
})

export default router
