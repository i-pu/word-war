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

export interface User {
  playerId: string
}

export enum Scene {
  None,
  // 待機中
  Matching,
  // ゲーム中
  Gaming,
  // 終了
  End
}

export interface Word {
  id: string
  text: string
}

interface GameState {
  scene: Scene
  roomId: string
  users: User[]
  limit: number
  timer: number
  score: number
  client: WordWarPromiseClient
  words: Word[]
}

interface RoomInfo {
  users: User[]
  timer: number
  limit: number
  roomId: string
}

export const game: Module<GameState, RootState> = {
  namespaced: true,
  state: {
    scene: Scene.None,
    roomId: '',
    timer: 0,
    users: [],
    limit: 0,
    score: 0,
    client: new WordWarPromiseClient(process.env.VUE_APP_API_ENDPOINT),
    words: []
  },
  getters: {
    getWords: state => state.words,
    score: state => state.score,
    roomId: state => state.roomId,
    users: state => state.users,
    scene: state => state.scene
  },
  mutations: {
    prepareRoom(commit) {
      commit.words = []
      commit.score = 0
      commit.scene = Scene.Gaming
    },
    score(commit, score: number) {
      commit.score = score
    },
    push(commit, word: Word) {
      commit.words.push(word)
    },
    setScene(commit, scene: Scene) {
      console.log(`${Scene[commit.scene]} -> ${Scene[scene]}`)
      commit.scene = scene
    },
    setTimer(commit, timer: number) {
      commit.timer = timer
    },
    setRoomInfo(commit, info: RoomInfo) {
      console.log('setRoomInfo')
      commit.users = info.users
      commit.timer = info.timer
      commit.limit = info.limit
      commit.roomId = info.roomId
    }
  },
  actions: {
    async match({ commit, state, rootGetters }) {
      console.log('Matching ...')
      console.log(rootGetters.userId)
      const matchingReq: MatchingRequest = new MatchingRequest()
      matchingReq.setUserid(rootGetters.userId)
      const stream = state.client.matching(matchingReq)
      commit('setScene', Scene.Matching)

      stream.on('data', matchingRes => {
        const [roomId, pbUsers, limit, timer] = [
          matchingRes.getRoomid(),
          matchingRes.getUserList(),
          matchingRes.getRoomuserlimit(),
          matchingRes.getTimerseconds()
        ]
        const users: User[] = pbUsers.map(u => ({ playerId: u.getUserid() }))
        console.log(`${roomId} ${users} ${limit} ${timer}`)

        commit('setRoomInfo', { roomId, users, timer, limit })
      })

      stream.on('status', status => {
        console.log('status', status)
        if (status.code === 0) {
          // ゲーム開始
          commit('prepareRoom')
        } else {
          throw `やばいね, ${status}`
        }
      })

      stream.on('error', res => {
        throw res
      })
    },

    async start({ commit, state, rootGetters }) {
      if (!state.roomId) {
        throw 'RoomId is empty'
      }

      setInterval(() => {
        commit('setTimer', state.timer - 1)
      }, 1000)

      console.log(`roomId in store: ${state.roomId}`)

      const gameReq: GameRequest = new GameRequest()
      gameReq.setRoomid(state.roomId)
      gameReq.setUserid(rootGetters.userId)
      const stream = state.client.game(gameReq)

      // on message
      stream.on('data', res => {
        const [roomId, userId, message] = [res.getRoomid(), res.getUserid(), res.getMessage()]
        console.log(`game response data ${roomId} ${userId} ${message}`)
        commit('push', { id: userId, text: message })
      })

      // fire on stream end
      stream.on('status', status => {
        console.log('status', status)
        if (status.code === 0) {
          console.log("無事終わりました")
          commit('setScene', Scene.End)
        } else {
          console.error(`無事に終わりませんでした ${status}`)
        }
      })

      stream.on('error', error => {
        console.log('Error')
        console.error(error)
      })
    },
    async say(
      { state, commit, rootGetters },
      { message }: { message: string }
    ) {
      const req: SayRequest = new SayRequest()
      req.setRoomid(state.roomId)
      req.setUserid(rootGetters.userId)
      req.setMessage(message)

      console.log(`Said ${req.getMessage()}`)

      try {
        const res = await state.client.say(req)
        console.log(`Response: ${res.getMessage()}`)
      } catch (e) {
        console.error(e)
      }
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
