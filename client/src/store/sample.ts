import { Module } from 'vuex'
import { RootState } from '@/store/root_state'

interface ISample {
  count: number
}

export const sample: Module<ISample, RootState> = {
  namespaced: true,
  state: {
    count: 0
  },
  mutations: {
    increment(commit) {
      commit.count++
    }
  }
}
