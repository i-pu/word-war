import Vue from 'vue'
import Vuex from 'vuex'

import { RootState } from '@/store/root_state'
import { user } from '@/store/user'
import { sample } from '@/store/sample'

Vue.use(Vuex)

export default new Vuex.Store<RootState>({
  state: {
    version: '0.1.0'
  },
  modules: {
    user,
    sample
  }
})
