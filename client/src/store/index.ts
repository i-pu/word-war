import Vue from 'vue'
import Vuex from 'vuex'

import firebase from '@/config/firebase'
import { RootState, User } from '@/store/root_state'
import { sample } from '@/store/sample'
import { game } from '@/store/game'
import { WordWarPromiseClient } from '@/pb/word_war_grpc_web_pb'
import { HealthCheckRequest, HealthCheckResponse } from '@/pb/word_war_pb'

Vue.use(Vuex)
// TODO: linterの細かい調整

export default new Vuex.Store<RootState>({
  state: {
    // server info
    version: '',
    serverHealth: false,
    user: {
      userId: ''
    }
  },
  mutations: {
    setUserId(state, { userId }: { userId: string }) {
      state.user.userId = userId
    },
    serverHealth(state, res: HealthCheckResponse) {
      state.version = res.getServerversion()
      state.serverHealth = res.getActive()
    }
  },
  getters: {
    userId: state => {
      return state.user.userId
    },
    version: state => {
      return state.version
    },
    serverHealth: state => {
      return state.serverHealth
    }
  },
  actions: {
    async healthCheck({ commit }) {
      const client = new WordWarPromiseClient(process.env.VUE_APP_API_ENDPOINT)
      const req = new HealthCheckRequest()
      const res = await client.healthCheck(req).catch(console.error)

      console.log(res)

      // :thinking_face:
      commit('serverHealth', res)
    },

    async signIn({ commit }, { email, password }) {
      const result = await firebase
        .auth()
        .signInWithEmailAndPassword(email, password)
        .catch(console.error)

      if (!result || !result.user) {
        throw new Error(`can't authorized: ${email}, ${password}`)
      }

      commit('setUserId', { userId: result.user.uid })

      console.log(`signIn: ${email}, ${password}`)
    },

    async signUp({ commit }, { email, password }) {
      // create user
      const result = await firebase
        .auth()
        .createUserWithEmailAndPassword(email, password)
        .catch(console.error)

      if (!result || !result.user) {
        throw new Error(`can't authorized: ${email}, ${password}`)
      }

      commit('user/setUserId', { userId: result.user.uid })

      console.log(`signUp: ${email}, ${password}`)
    }
  },
  modules: {
    sample,
    game
  }
})
