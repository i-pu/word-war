import { ports, firebase } from './index'
import pb from './word_war_grpc_web_pb'

// gRPC API のエンドポイント
const endpoint = 'http://localhost:8080'
// gRPC のクライアント
const client = new pb.WordWarPromiseClient()

// ====================
//        Top
// ====================

// ログインボタンが押されて Elm から呼ばれる
ports.signinWithFirebase.subscribe(async ({ email, password }) => {
  // firebase auth でユーザー認証
  const auth = await firebase.auth()
    .signInWithEmailAndPassword(email, password)
    .catch(console.log)

  if (!auth) {
    return
  }

  console.log(`logged in as ${auth.user.uid}`)

  ports.signinCallback.send({ uid: auth.user.uid })
})

// サインアップボタンが押されたときの処理
ports.signupWithFirebase.subscribe(async ({ email, password }) => {
  console.log({ email, password })
  const auth = await firebase.auth()
    .createUserWithEmailAndPassword(email, password)
    .catch(console.error)

  if (!auth) {
    return
  }

  console.log(`signed up and logged in as ${auth.user.uid}`)

  ports.signinCallback.send({ uid: auth.user.uid })
})

// ====================
//        Home
// ====================


// ====================
//        Game
// ====================
ports.startGame.subscribe(async ({ userId }) => {
  // Server Streaming RPCs
  const req = new pb.GameRequest()
  req.setUserId(userId)
  const stream = client.game(req)

  stream.on('data', res => {
    const message = res.getMessage()
    const userId = req.getUserId()
    // JS -> Elm
    ports.onMessage.send({ userId, message })
  })

  stream.on('status', status => {
    console.log(status)
  })
})

ports.say.subscribe(async ({ userId, message }) => {
  console.log({ userId, message })
  // gRPC Unary RPCs
  const req = new pb.SayRequest()
  req.setUserId(userId)
  req.setMessage(message)

  const [ userId, message ] = await client.Say(req)
    .then(res => res.array)
    .catch(console.error)

  // Elm の onMessage を呼ぶ
  ports.toElm.send({ userId, message })
})

// ====================
//       Result
// ====================

ports.requestResult.subscribe(({ userId }) => {
  const req = new pb.ResultRequest()
  req.setUserId(userId)

  const [ userId, score ] = await client.Result(req)
    .then(res => res.array)
    .catch(console.error)

  ports.resultCallback.send({ userId, score })
})
