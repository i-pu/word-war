import Vue from "vue"
import VueRouter from "vue-router"

Vue.use(VueRouter)
// TODO: router.beforeEachを使ってログイン済みか確認する
// see also https://qiita.com/sunadorinekop/items/f3486da415d3024c7ed4

const routes = [
  {
    path: "/",
    name: "top",
    component: () => import("@/views/Top.vue")
  },
  {
    path: "/about",
    name: "about",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/About.vue")
  },
  {
    path: "/home",
    name: "home",
    component: () => import("@/views/Home.vue")
  },
  {
    path: "/game",
    name: "game",
    component: () => import("@/views/Game.vue")
  },
  {
    path: "/result",
    name: "result",
    component: () => import("@/views/Result.vue")
  }
]

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
})

export default router