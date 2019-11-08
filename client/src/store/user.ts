import { Module } from 'vuex'
import { RootState } from '@/store/root_state'

interface IUser {
  uid: string
}

export const user: Module<IUser, RootState> = {
  namespaced: true,
  state: {
    uid: ''
  },
  mutations: {
    setUid(state, { uid }: { uid: string }) {
      state.uid = uid
    }
  },
  getters: {
    uid: state => {
      return state.uid
    }
  }
}
