import firebase from 'firebase/app'
import 'firebase/auth'
import 'firebase/firestore'
import credential from '../credential'

// firebase の初期化
firebase.initializeApp(credential)
import pb from './word_war_grpc_web_pb'
// Parcel を使っているので直接 import できる
const { Elm } = require('./Main.elm')

// Elm を初期化して #app に配置
const app = Elm.Main.init({
  node: document.getElementById('app'),
  // Elm flags
  flags: "Hoge",
})

const store = {}

const isDevelop = process.env.NODE_ENV === 'development'
const isStub = true // isDevelop

// gRPC API のエンドポイント
const endpoint = isDevelop
  ? process.env.ELM_APP_API_ENDPOINT_DEV
  : process.env.ELM_APP_API_ENDPOINT

if (isStub) {
  console.log(`stub mode`)
} else {
  console.log(`endpoint is ${endpoint}`)
}

// gRPC のクライアント
const client = new pb.WordWarPromiseClient(endpoint)

// ====================
//        Top
// ====================

// ログインボタンが押されて Elm から呼ばれる
app.ports.signinWithFirebase.subscribe(async ({ email, password }) => {
  try {
    // firebase auth でユーザー認証
    const auth = await firebase.auth()
      .signInWithEmailAndPassword(email, password)

    const userId = auth.user.uid

    // create user data
    const snapshot = await firebase.firestore()
      .collection('users')
      .doc(userId)
      .get()

    const userData = snapshot.data()

    console.log(`logged in as ${userId}`)
    console.log(userData)

    app.ports.signinCallback.send({ uid: userId, rating: userData.rating })

  } catch (e) {
    console.error(e)
    return
  }
})

// サインアップボタンが押されたときの処理
app.ports.signupWithFirebase.subscribe(async ({ name, email, password }) => {
  console.log({ name, email, password })

  try {
    // create user
    const auth = await firebase.auth()
      .createUserWithEmailAndPassword(email, password)

    const userId = auth.user.uid

    console.log(`signed up and logged in as ${userId}`)

    // create user data
    await firebase.firestore()
      .collection('users')
      .doc(userId)
      .set({
        name,
        rating: 1500,
        history: [1500]
      })

    app.ports.signinCallback.send({ uid: userId, rating: 1500 })

  } catch (e) {
    console.error(e)
    return
  }
})

// ====================
//        Home
// ====================

// ====================
//        Game
// ====================
app.ports.startGame.subscribe(async userId => {
  console.log('start game')


  if (!isStub) {
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

    stream.on('status', status => {
      // be sent when finish streaming
      if (status.code === 0) {
        app.ports.onFinish.send(null)
      }
    })
  } else {
    return new Promise(async () => {
      // stub
      store.messages = []
      for (let i = 0; i < 5; i++) {
        app.ports.onMessage.send({ userId: 'bcEudVn6dlYAdAPqWSrDA7dWAv82', message: 'りんご' })
        await new Promise(resolve => setTimeout(resolve, 2000))
      }
    })
  }
})

app.ports.say.subscribe(async ({ userId, message }) => {
  console.log({ userId, message })

  if (!isStub) {
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

    if (store.messages.length >= 3) {
      app.ports.onFinish.send(null)
    }
  }
})

// ====================
//       Result
// ====================

app.ports.requestResult.subscribe(async userId => {
  if (!isStub) {
    const req = new pb.ResultRequest()
    req.setUserid(userId)

    const res = await client.result(req)
      .catch(console.error)

    app.ports.onResult.send({ "userId": res.getUserid(), "score": res.getScore() })
  } else {
    const score = store.messages ? store.messages.length : 0
    console.log({ userId, score })
    app.ports.onResult.send({ userId, score })
  }
})

export default app
