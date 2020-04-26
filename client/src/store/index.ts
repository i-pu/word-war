import { createStore } from 'pinia'
// import * as api from '@/api'
import * as api from '@/api/index.mock'
import { defaultUser, User } from '@/model'

export enum SignInMethod {
  Email, Google, Twitter, GitHub
}

export const useRootStore = createStore({
  id: 'root',
  state: () => ({
    status: {
      active: false,
      version: '',
    },
    user: null as User | null,
    auth: null as firebase.User | null,
  }),
  getters: {
    user: state => {
      return state.user
    },
    userId: state => {
      return state.user ? state.user!.userId : ''
    },
    status: state => {
      return state.status
    },
    auth: state => state.auth
  },
  actions: {
    async healthCheck() {
      const status = await api.healthCheck()
        .catch(console.error)
      if (!status) {
        console.warn('やばいよ')
        return
      }
      this.state.status = status
      console.log(this.state.status)
    },
    async updateAuth(auth: firebase.User) {
      this.state.user = defaultUser()
      this.state.auth = auth
      if (auth.displayName) {
        this.state.user.name = auth.displayName
      }
      this.state.user.avatarUrl = auth.photoURL || ''
    },
    async signIn({ email, password }: { email: string, password: string }) {
      const user = await api.signIn({ email, password })
      console.log(`${this.state.user ? this.state.user.userId : 'null'} -> ${user.userId}`)
      this.state.user = user
    },
    async signUp({ email, password }: { email: string, password: string }) {
      // create user
      const user = await api.signUp({ email, password })
      this.state.user = user

      console.log(`signUp: ${email}, ${password}`)
    },
    async google() {
      const user = await api.signUpWithGoogle()
      this.state.user = user

      console.log(`signUp with google`)
    },
    async signUpWithTwitter() {
      const user = await api.signUpWithTwitter()
      this.state.user = user

      console.log(`signUp with twitter`)
    },
    async signUpWithGitHub() {
      // const user = await api.signUpWithGitHub()
      // TODO: リダイレクトなら何も買いってこない
      // ポップアップならUser帰ってくる
      await api.signUpWithGitHub()
      // this.state.user = user

      console.log(`signUp with github`)
    },
    async signOut() {
      this.state.user = null
      await api.signOut()
    },
  }
})
