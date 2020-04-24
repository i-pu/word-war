import Vue from 'vue'
import Buefy from 'buefy'
import VueHead from 'vue-head'

import 'buefy/dist/buefy.css'
import '@mdi/font/css/materialdesignicons.css'
import 'bulma/css/bulma.css'
import '@/config/firebase'

import App from './App.vue'
import router from './router'
import VueCompositionApi from "@vue/composition-api"
Vue.config.productionTip = false

Vue.use(Buefy)
Vue.use(VueHead)
Vue.use(VueCompositionApi)

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
