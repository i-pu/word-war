import Vue from 'vue'
import VueRouter, { Route, NavigationGuard } from 'vue-router'
import { useRootStore } from '@/store'
import firebase from 'firebase'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'top',
    component: () => import('@/views/Top.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/Login.vue'),
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
  },
  {
    path: '/user/:id',
    name: 'user',
    component: () => import('@/views/User.vue'),
    meta: { requiresAuth: true }
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

// const onStateChanged = (): Promise<firebase.User | null> => {
//   return new Promise((resolve, reject) => {
//     try {
//       const unsubscribe = firebase.auth().onAuthStateChanged(async user => user ? resolve(user) : resolve(null))
//       unsubscribe()
//     } catch (e) {
//       reject(e)
//     }
//   })
// }

router.beforeEach(async (to, from, next) => {
  console.log('enter@beforeEach')
  await new Promise(res => setTimeout(res, 500))

  const { status, healthCheck, user } = useRootStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  console.log(`したいです。from: ${from.path} -> to: ${to.path}`)

  await healthCheck()
  if (!status.value.active && to.path !== '/') {
    console.log('やばいですアクティブじゃないです')
    next({ path: '/' })
    return
  }

  // un-authrized
  if (requiresAuth) {
    if (!user.value) {
      console.log('ログインしてください')
      next({ path: '/' })
    }
  }

  next()
})

export default router
