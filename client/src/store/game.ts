import { createStore } from 'pinia'
import { useRootStore } from '@/store'
import * as api from '@/api/game'
// import * as api from '@/api/game.mock'
import { Scene, Player, Word, RoomInfo } from '@/model'

export const useGameStore = createStore({
  id: 'game',
  state: () => ({
    scene: Scene.None,
    roomId: '',
    timer: 0,
    players: [] as Player[],
    limit: 0,
    score: 0,
    words: [] as Word[],
    // TODO: bad
    cancel: (() => {}) as () => void,
  }),
  getters: {
    words: state => state.words,
    score: state => state.score,
    roomId: state => state.roomId,
    players: state => state.players,
    scene: state => state.scene
  },
  actions: {
    prepareRoom() {
      this.state.words = []
      this.state.score = 0
      this.state.scene = Scene.Gaming
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
      this.state.players = info.players
      this.state.timer = info.timer
      this.state.limit = info.limit
      this.state.roomId = info.roomId
    },
    async match() {
      const { userId } = useRootStore()
      const matchOnData = (roomInfo: RoomInfo) => {
        console.log(roomInfo)
        this.setRoomInfo(roomInfo)
      }
      const matchOnEnd = () => {
        console.log(`[room ${this.state.roomId}] Match end`)
        // ゲーム開始
        this.prepareRoom()
      }

      this.setScene(Scene.Matching)

      this.state.cancel = await api.match(userId.value, matchOnData, matchOnEnd)
    },
    async cancelMatching() {
      if (typeof this.state.cancel === 'function') {
        // @ts-ignore
        this.state.cancel()
        this.state.scene = Scene.None
      }
    },
    async start() {
      const { userId } = useRootStore()
      const onData = (word: Word) => {
        this.state.words.push(word)
      }
      const onEnd = () => {
        this.setScene(Scene.End)
      }

      // timer
      setInterval(() => {
        this.setTimer(this.state.timer - 1)
      }, 1000)

      console.log(`roomId in store: ${this.state.roomId}`)

      api.start(userId.value, this.roomId.value, onData, onEnd)
    },
    async say(message: string) {
      const { userId } = useRootStore()
      await api.say(this.roomId.value, userId.value, message)
        .catch(console.error)
    },

    async result() {
      const { userId } = useRootStore()
      this.state.score = await api.result(userId.value, this.state.roomId)
    }
  }
})
