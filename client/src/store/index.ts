import Vue from "vue";
import Vuex, { Module } from "vuex";

Vue.use(Vuex);

interface IUser {
  uid: string;
}

const user: Module<IUser, RootState> = {
  namespaced: true,
  state: {
    uid: ""
  },
  mutations: {
    setUid(state, { uid }: { uid: string }) {
      state.uid = uid;
    }
  },
  getters: {}
};

interface ISample {
  count: number;
}

interface RootState {
  version: string;
}

const sample: Module<ISample, RootState> = {
  namespaced: true,
  state: {
    count: 0
  },
  mutations: {
    increment(commit) {
      commit.count++;
    }
  }
};

export default new Vuex.Store<RootState>({
  state: {
    version: "0.1.0"
  },
  modules: {
    user,
    sample
  }
});
