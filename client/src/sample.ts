import { WordWarPromiseClient } from '@/pb/word_war_grpc_web_pb'
import { SayRequest, GameRequest, ResultRequest } from '@/pb/word_war_pb'

const wordWarPromiseClient = new WordWarPromiseClient(
  'http://localhost:50051',
  null,
  null
)

const req: SayRequest = new SayRequest()
req.setUserid('')
req.setMessage('')

wordWarPromiseClient.say(req)

const reqgame: GameRequest = new GameRequest()
reqgame.setUserid('')
const stream = wordWarPromiseClient.game(reqgame)
stream.on('data', res => {
  console.log(res)
  res.getUserid()
  res.getMessage()
})

const reqresult: ResultRequest = new ResultRequest()
reqresult.setUserid('')
reqresult.setUserid('')
