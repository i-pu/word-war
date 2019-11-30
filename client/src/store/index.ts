import Vue from 'vue'
import Vuex from 'vuex'

import firebase from '@/config/firebase'
import { RootState, User } from '@/store/root_state'
import { sample } from '@/store/sample'
import { game } from '@/store/game'

Vue.use(Vuex)

export default new Vuex.Store<RootState>({
  state: {
    version: '0.0.1',
    user: {
      userId: ''
    }
  },
  mutations: {
    setUserId(state, { userId }: { userId: string }) {
      state.user.userId = userId
    }
  },
  getters: {
    userId: state => {
      return state.user.userId
    }
  },
  actions: {
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
