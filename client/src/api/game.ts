/**
 * Game Api
 */
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
import { Player, RoomInfo, Word } from '@/model'

export const client = new WordWarPromiseClient(process.env.VUE_APP_API_ENDPOINT)

async function* playersGenerator(limit: number = 2, wait: number = 2) {
  let i = 0
  while (i < limit) {
    await new Promise(resolve => setTimeout(resolve, wait * 1000))
    yield { playerId: `id-${i}` }
  }
}

async function* messagesGenerator(count: number = 10, wait: number = 1) {
  let i = 0
  let id = 'hoge'
  while (i < count) {
    // sleep 1s
    await new Promise(resolve => setTimeout(resolve, wait * 1000))
    yield { id, text: `message${i}` }
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
  if (!roomId) {
    throw 'RoomId is empty'
  }

  console.log(`Matching ${roomId} ...`)

  console.log(userId)

  const matchingReq: MatchingRequest = new MatchingRequest()
  matchingReq.setUserid(userId)

  const stream = client.matching(matchingReq)

  stream.on('data', (matchingRes: MatchingResponse) => {
    // convert
    const [roomId, pbUsers, limit, timer] = [
      matchingRes.getRoomid(),
      matchingRes.getUserList(),
      matchingRes.getRoomuserlimit(),
      matchingRes.getTimerseconds()
    ]
    const players: Player[] = pbUsers.map(u => ({ playerId: u.getUserid() }))
    console.log(`${roomId} ${players} ${limit} ${timer}`)

    // callback
    onData({ roomId, players, timer, limit })
  })

  stream.on('status', status => {
    console.log('status', status)
    if (status.code === 0) {
      // ゲーム開始
      onEnd()
    } else {
      throw `やばいね, ${status}`
    }
  })
}

/**
 * Start
 */
export const start = async (
  userId: string,
  roomId: string,
  onData: (word: Word) => void,
  onEnd: () => void,
) => {
  console.log(`roomId in store: ${roomId}`)

  const gameReq: GameRequest = new GameRequest()
  gameReq.setRoomid(roomId)
  gameReq.setUserid(userId)
  const stream = client.game(gameReq)

  // on message
  stream.on('data', (res: SayRequest) => {
    const [roomId, userId, message] = [res.getRoomid(), res.getUserid(), res.getMessage()]
    console.log(`game response data ${roomId} ${userId} ${message}`)
    onData({ id: userId, text: message })
  })

  // fire on stream end
  stream.on('status', status => {
    console.log('status', status)
    if (status.code === 0) {
      console.log("無事終わりました")
      onEnd()
    } else {
      throw `無事に終わりませんでした ${status}`
    }
  })

  stream.on('error', error => {
    throw error
  })
}

export const say = async (roomId: string, userId: string, message: string): Promise<void> => {
  const req: SayRequest = new SayRequest()
  req.setRoomid(roomId)
  req.setUserid(userId)
  req.setMessage(message)
  console.log(`Said ${req.getMessage()}`)

  try {
    const res = await client.say(req)
    console.log(`Response: ${res.getMessage()}`)
  } catch (e) {
    console.error(e)
  }
}

export const result = async (userId: string, roomId: string, ): Promise<number> => {
  const req = new ResultRequest()
  req.setUserid(userId)
  req.setRoomid(roomId)
  return client.result(req).then((res) => {
    console.log(res)
    return +res.getScore()
  })
}
