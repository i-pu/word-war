import { createStore } from 'pinia'
// import * as api from '@/api'
import * as api from '@/api/index.mock'
import { defaultUser } from '@/model'

export const useRootStore = createStore({
  id: 'root',
  state: () => ({
    status: {
      active: false,
      version: '',
    },
    user: defaultUser()
  }),
  getters: {
    userId: state => {
      return state.user.userId
    },
    status: state => {
      return state.status
    },
  },
  actions: {
    async healthCheck() {
      try {
        this.state.status = await api.healthCheck()
        console.log(this.state.status)
      } catch (e) {
        console.error(e)
      }
    },
    async signIn({ email, password }: { email: string, password: string }) {
      const user = await api.signIn({ email, password })

      console.log(`${this.state.user.userId} -> ${user.userId}`)
      this.state.user = user
    },
    async signUp({ email, password }: { email: string, password: string }) {
      // create user
      const user = await api.signUp({ email, password })
      this.state.user = user

      console.log(`signUp: ${email}, ${password}`)
    },
    async signOut() {
      this.state.user = defaultUser()
      await api.signOut()
    },
  }
})
