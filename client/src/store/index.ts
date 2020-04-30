import { createStore } from 'pinia'
import * as api from '@/api'
// import * as api from '@/api/index.mock'
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
  },
  actions: {
    async healthCheck() {
      const status = await api.healthCheck()
        .catch(console.error)
      if (!status) {
        console.warn('statusがなんもはいってないよやばいよhealthcheckとおってないかもよ')
        return
      }
      this.state.status = status
      console.log(this.state.status)
    },
    async onAuthChanged(auth: firebase.User) {
      if (!auth) {
        return
      }

      const firstLogin = !(await api.existUserData(auth.uid))

      // はじめて google でログイン
      // firebase に ユーザーデータがあるか確認
      if (firstLogin) {
        await api.registerUser(auth)
      }

      // fetch user
      const user = await api.getUserdata(auth.uid)
      // set user
      this.state.user = user
      console.log(`${this.state.user ? this.state.user.userId : 'null'} -> ${user.userId}`)
    },
    async signIn({ email, password }: { email: string, password: string }) {
      await api.signIn({ email, password })
    },
    async signUp({ email, password }: { email: string, password: string }) {
      // create user
      await api.signUp({ email, password })
      console.log(`signUp: ${email}, ${password}`)
    },
    async google() {
      await api.signUpWithGoogle()
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
