import firebase from 'firebase/app'
import 'firebase/auth'
import credential from '../credential'

// firebase の初期化
firebase.initializeApp(credential)
import pb from './word_war_grpc_web_pb'
// Parcel を使っているので直接 import できる
const { Elm } = require('./main.elm')

// Elm を初期化して #app に配置
const app = Elm.Main.init({
  node: document.getElementById('app'),
  // Elm flags
  flags: "Hoge",
})

const STUBBED = true

const store = {}

// gRPC API のエンドポイント
const endpoint = 'http://localhost:8080'
// gRPC のクライアント
const client = new pb.WordWarPromiseClient(endpoint)

// ====================
//        Top
// ====================

// ログインボタンが押されて Elm から呼ばれる
app.ports.signinWithFirebase.subscribe(async ({ email, password }) => {
  // firebase auth でユーザー認証
  const auth = await firebase.auth()
    .signInWithEmailAndPassword(email, password)
    .catch(console.log)

  if (!auth) {
    return
  }

  console.log(`logged in as ${auth.user.uid}`)

  app.ports.signinCallback.send({ uid: auth.user.uid })
})

// サインアップボタンが押されたときの処理
app.ports.signupWithFirebase.subscribe(async ({ email, password }) => {
  console.log({ email, password })
  const auth = await firebase.auth()
    .createUserWithEmailAndPassword(email, password)
    .catch(console.error)

  if (!auth) {
    return
  }

  console.log(`signed up and logged in as ${auth.user.uid}`)

  app.ports.signinCallback.send({ uid: auth.user.uid })
})

// ====================
//        Home
// ====================


// ====================
//        Game
// ====================
app.ports.startGame.subscribe(async userId => {
  console.log('start game')
  if (!STUBBED) {
    // Server Streaming RPCs
    const req = new pb.GameRequest()
    req.setUserid(userId)
    const stream = client.game(req)

    stream.on('data', res => {
      const message = res.getMessage()
      const userId = req.getUserid()
      // JS -> Elm
      app.ports.onMessage.send({ userId, message })
    })

    stream.on('end', () => {
      app.ports.onFinish.send(null)
    })
  } else {
    return new Promise(async () => {
      // stub
      store.messages = []
      for (let i = 0; i < 5; i++) {
        app.ports.onMessage.send({ userId: 'u012', message: 'message' })
        await new Promise(resolve => setTimeout(resolve, 2000))
      }

      app.ports.onFinish.send(null)
    })
  }
})

app.ports.say.subscribe(async ({ userId, message }) => {
  console.log({ userId, message })

  if (!STUBBED) {
    // gRPC Unary RPCs
    const req = new pb.SayRequest()
    req.setUserid(userId)
    req.setMessage(message)

    const res = await client.say(req)
      .catch(console.error)

    // Elm の onMessage を呼ぶ
    app.ports.onMessage.send({ "userId": res.getUserid(), "message": res.getMessage() })
  } else {
    store.messages.push({ userId, message })
    app.ports.onMessage.send({ userId, message })
  }
})

// ====================
//       Result
// ====================

app.ports.requestResult.subscribe(async userId => {
  if (!STUBBED) {
    const req = new pb.ResultRequest()
    req.setUserid(userId)

    const res = await client.result(req)
      .catch(console.error)

    app.ports.onResult.send({ "userId": res.getUserid(), "score": res.getScore() })
  } else {
    const score = store.messages.length
    console.log({ userId, score })
    app.ports.onResult.send({ userId, score })
  }
})

export default app
