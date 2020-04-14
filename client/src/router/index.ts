import Vue from 'vue'
import VueRouter from 'vue-router'
import firebase from 'firebase'
import store from './../store'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'top',
    component: () => import('@/views/Top.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/home',
    name: 'home',
    component: () => import('@/views/Home.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/waiting',
    name: 'waiting',
    component: () => import('@/views/Waiting.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/game',
    name: 'game',
    component: () => import('@/views/Game.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/result',
    name: 'result',
    component: () => import('@/views/Result.vue'),
    meta: { requiresAuth: true }
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  // un-authrized
  if (requiresAuth) {
    if (firebase.auth().currentUser) {
      next()
    } else {
      console.log('ログインしてください')
      next('/')
    }
  }

  next()

  // when un-authrized
  firebase.auth().onAuthStateChanged(user => {
    if (!user) {
      next('/')
    }
  })
})

export default router
