import { createStore } from 'pinia'
import firebase from '@/config/firebase'
import { WordWarPromiseClient } from '@/pb/word_war_grpc_web_pb'
import { HealthCheckRequest, HealthCheckResponse } from '@/pb/word_war_pb'

export interface User {
  userId: string
}

export interface RootState {
  version: string
  user: User
  serverHealth: boolean
}


const setInitialUserdata = (uid: string) => {
  return firebase
    .firestore()
    .collection('users')
    .doc(uid)
    .set({
      history: [{ date: new Date(), rating: 1500 }],
      name: 'AAA',
      rating: 1500
    })
}

export const useRootStore = createStore({
  id: 'root',
  state: () => ({
    // server info
    version: '',
    serverHealth: false,
    user: {
      userId: ''
    }
  }),
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
    setUserId({ userId }: { userId: string }) {
      console.log(`${this.state.user.userId} -> ${userId}`)
      this.state.user.userId = userId
    },
    serverHealth(res: HealthCheckResponse) {
      this.state.version = res.getServerversion()
      this.state.serverHealth = res.getActive()
    },
    async healthCheck() {
      try {
        const client = new WordWarPromiseClient(process.env.VUE_APP_API_ENDPOINT)
        const req = new HealthCheckRequest()
        const res = await client.healthCheck(req)

        console.log(res)
        this.serverHealth(res)
      } catch (e) {
        console.error(e)
      }
    },
    async signIn({ email, password }: { email: string, password: string }) {
      const result = await firebase
        .auth()
        .signInWithEmailAndPassword(email, password)
        .catch(console.error)

      if (!result || !result.user) {
        throw new Error(`can't authorized: ${email}, ${password}`)
      }

      this.setUserId( { userId: result.user.uid })

      console.log(`signIn: ${email}, ${password}`)
    },
    async signUp({ email, password }: { email: string, password: string }) {
      // create user
      const result = await firebase
        .auth()
        .createUserWithEmailAndPassword(email, password)
        .catch(console.error)

      if (!result || !result.user) {
        throw new Error(`can't authorized: ${email}, ${password}`)
      }

      await setInitialUserdata(result.user.uid).catch(console.error)

      console.log('initialized userdata in firestore')

      this.setUserId({ userId: result.user.uid })

      console.log(`signUp: ${email}, ${password}`)
    },
    async signOut() {
      this.setUserId({ userId: '' })
      await firebase
        .auth()
        .signOut()
        .catch(console.error)
    },
  }
})
