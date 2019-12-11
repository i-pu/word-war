import { Module } from 'vuex'
import { RootState } from '@/store/root_state'
import { WordWarPromiseClient } from '@/pb/word_war_grpc_web_pb'
import {
  SayRequest,
  MatchingRequest,
  GameRequest,
  GameResponse,
  ResultRequest,
  ResultResponse
} from '@/pb/word_war_pb'

interface GameState {
  isPlaying: boolean
  roomId: string
  score: number
  client: WordWarPromiseClient
  words: GameResponse[]
}

export const game: Module<GameState, RootState> = {
  namespaced: true,
  state: {
    isPlaying: false,
    roomId: '',
    score: 0,
    client: new WordWarPromiseClient(process.env.VUE_APP_API_ENDPOINT),
    words: []
  },
  getters: {
    getWords: state => state.words,
    score: state => state.score,
    roomId: state => state.roomId
  },
  mutations: {
    start(commit) {
      commit.isPlaying = true
    },
    end(commit) {
      commit.isPlaying = false
    },
    score(commit, score: number) {
      commit.score = score
    },
    room(commit, roomId: string) {
      commit.roomId = roomId
    },
    push(commit, word: GameResponse) {
      commit.words.push(word)
    }
  },
  actions: {
    async matchAndStart({ commit, state, rootGetters }) {
      console.log(rootGetters.userId)
      const matchingReq: MatchingRequest = new MatchingRequest()
      matchingReq.setUserid(rootGetters.userId)
      const matchingRes = await state.client
        .matching(matchingReq)
        .catch(console.error)

      if (!matchingRes) {
        return
      }

      console.log(matchingRes.getRoomid())
      commit('room', matchingRes.getRoomid())

      const gameReq: GameRequest = new GameRequest()
      gameReq.setRoomid(matchingRes.getRoomid())
      gameReq.setUserid(rootGetters.userId)
      const stream = state.client.game(gameReq)

      stream.on('data', res => {
        console.log(`${res.getRoomid()} ${res.getUserid()} ${res.getMessage()}`)
        commit('push', res)
      })

      stream.on('status', status => {
        console.log('status', status)
        if (status.code === 0) {
          commit('end')
        }
      })

      stream.on('error', res => {
        console.log('error', res)
      })

      commit('start')
    },
    async say(
      { state, commit, rootGetters },
      { message }: { message: string }
    ) {
      // TODO: uid
      const req: SayRequest = new SayRequest()
      req.setRoomid(state.roomId)
      req.setUserid(rootGetters.userId)
      req.setMessage(message)

      console.log(req)

      await state.client.say(req).catch(console.error)
    },
    async result({ state, commit, rootGetters }) {
      const req = new ResultRequest()
      req.setUserid(rootGetters.userId)
      req.setRoomid(state.roomId)
      const result = await state.client.result(req).catch(console.error)
      if (result) {
        console.log(result)
        commit('score', result.getScore())
      }
    }
  }
}