import { createStore } from 'pinia'
import { useRootStore } from '@/store'
import { WordWarPromiseClient } from '@/pb/word_war_grpc_web_pb'
import {
  SayRequest,
  MatchingRequest,
  GameRequest,
  GameResponse,
  ResultRequest,
  ResultResponse,
  MatchingResponse
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

export const useGameStore = createStore({
  id: 'game',
  state: () => ({
    scene: Scene.None,
    roomId: '',
    timer: 0,
    users: [] as User[],
    limit: 0,
    score: 0,
    client: new WordWarPromiseClient(process.env.VUE_APP_API_ENDPOINT),
    words: [] as Word[],
  }),
  getters: {
    words: state => state.words,
    score: state => state.score,
    roomId: state => state.roomId,
    users: state => state.users,
    scene: state => state.scene
  },
  actions: {
    prepareRoom() {
      this.state.words = []
      this.state.score = 0
      this.state.scene = Scene.Gaming
    },
    score(score: number) {
      this.state.score = score
    },
    push(word: Word) {
      this.state.words.push(word)
    },
    setScene(scene: Scene) {
      console.log(`${Scene[this.state.scene]} -> ${Scene[scene]}`)
      this.state.scene = scene
    },
    setTimer(timer: number) {
      this.state.timer = timer
    },
    setRoomInfo(info: RoomInfo) {
      console.log('setRoomInfo')
      this.state.users = info.users
      this.state.timer = info.timer
      this.state.limit = info.limit
      this.state.roomId = info.roomId
    },
    async match() {
      console.log('Matching ...')
      const { userId } = useRootStore()
      console.log(userId.value)
      const matchingReq: MatchingRequest = new MatchingRequest()
      matchingReq.setUserid(userId.value)
      const stream = this.state.client.matching(matchingReq)
      this.setScene(Scene.Matching)

      stream.on('data', (matchingRes: MatchingResponse) => {
        const [roomId, pbUsers, limit, timer] = [
          matchingRes.getRoomid(),
          matchingRes.getUserList(),
          matchingRes.getRoomuserlimit(),
          matchingRes.getTimerseconds()
        ]
        const users: User[] = pbUsers.map(u => ({ playerId: u.getUserid() }))
        console.log(`${roomId} ${users} ${limit} ${timer}`)

        this.setRoomInfo({ roomId, users, timer, limit })
      })

      stream.on('status', status => {
        console.log('status', status)
        if (status.code === 0) {
          // ゲーム開始
          this.prepareRoom()
        } else {
          throw `やばいね, ${status}`
        }
      })
    },

    async start() {
      const { userId } = useRootStore()

      if (!this.state.roomId) {
        throw 'RoomId is empty'
      }

      setInterval(() => {
        this.setTimer(this.state.timer - 1)
      }, 1000)

      console.log(`roomId in store: ${this.state.roomId}`)

      const gameReq: GameRequest = new GameRequest()
      gameReq.setRoomid(this.state.roomId)
      gameReq.setUserid(userId.value)
      const stream = this.state.client.game(gameReq)

      // on message
      stream.on('data', (res: SayRequest) => {
        const [roomId, userId, message] = [res.getRoomid(), res.getUserid(), res.getMessage()]
        console.log(`game response data ${roomId} ${userId} ${message}`)
        this.push({ id: userId, text: message })
      })

      // fire on stream end
      stream.on('status', status => {
        console.log('status', status)
        if (status.code === 0) {
          console.log("無事終わりました")
          this.setScene(Scene.End)
        } else {
          console.error(`無事に終わりませんでした ${status}`)
        }
      })

      stream.on('error', error => {
        console.log('Error')
        console.error(error)
      })
    },
    async say(message: string) {
      const req: SayRequest = new SayRequest()
      const { userId } = useRootStore()
      req.setRoomid(this.state.roomId)
      req.setUserid(userId.value)
      req.setMessage(message)

      console.log(`Said ${req.getMessage()}`)

      try {
        const res = await this.state.client.say(req)
        console.log(`Response: ${res.getMessage()}`)
      } catch (e) {
        console.error(e)
      }
    },
    async result() {
      const { userId } = useRootStore()
      const req = new ResultRequest()
      req.setUserid(userId.value)
      req.setRoomid(this.state.roomId)
      const result = await this.state.client.result(req).catch(console.error)
      if (result) {
        console.log(result)
        this.score(+result.getScore())
      }
    }
  }
})
