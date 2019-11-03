import Vue from 'vue'
import Buefy from 'buefy'
import VueHead from 'vue-head'

import 'buefy/dist/buefy.css'
import '@mdi/font/css/materialdesignicons.css'

import App from './App.vue'
import router from './router'
import store from './store/index'

Vue.config.productionTip = false

Vue.use(Buefy)
Vue.use(VueHead)

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
