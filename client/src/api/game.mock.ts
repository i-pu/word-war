/**
 * Game Api
 */
import { RoomInfo, Word } from '@/model'

async function* playersGenerator(limit: number = 4, wait: number = 2) {
  let i = 0
  while (i < limit) {
    await new Promise(resolve => setTimeout(resolve, wait * 1000))
    yield { playerId: `id-${i++}` }
  }
}

async function* messagesGenerator(count: number = 10, wait: number = 1) {
  let i = 0
  let id = 'hoge'
  while (i < count) {
    // sleep 1s
    await new Promise(resolve => setTimeout(resolve, wait * 1000))
    yield { id, text: `message${i++}` }
  }
}

/**
 *  Matching
 */
export const match = async (
  userId: string,
  roomId: string,
  onData: (roomInfo: RoomInfo) => void,
  onEnd: () => void,
) => {
  // 2秒ごとに合計4人来る
  const players = []
  for await (let player of playersGenerator(2, 2)) {
    players.push(player)
    onData({ players, timer: 5, limit: 4, roomId: "some-roomid" })
  }

  onEnd()
}

export const start = async (
  userId: string,
  roomId: string,
  onData: (word: Word) => void,
  onEnd: () => void,
) => {
  for await (let message of messagesGenerator()) {
    onData(message)
  }

  onEnd()
}

export const say = async (roomId: string, userId: string, message: string): Promise<void> => {
  return
}

export const result = async (userId: string, roomId: string, ): Promise<number> => {
  // TODO: まともなデータを入れる
  return 10
}