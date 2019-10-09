import pb from './word_war_grpc_web_pb'
import firebase from 'firebase'
import credential from '../credential'

// firebase の初期化
firebase.initializeApp(credential)

// Parcel を使っているので直接 import できる
const { Elm } = require('./main.elm')

// gRPC API のエンドポイント
const endpoint = 'http://localhost:50051'
// gRPC のクライアント
const client = new pb.WordWarPromiseClient()
// Elm を初期化して #app に配置
const app = Elm.Main.init({
  node: document.getElementById('app'),
  // Elm flags
  flags: "Hoge"
})

// Elm から toJS が呼ばれたときの処理
// コールバックの頭に async をつけると非同期関数となり中で await が使える
// await には Promise を剥がす効果があり,
// Promise を返す関数を呼ぶ際につけると Promise の中身を取り出すことが出来る
// ```js
// const f = () => new Promise(resolve => resolve(1))
// const p = f() // :Promise<number>
// const v = await f() // 1 :number
// ```

// catch は Promise を返す関数にチェーンさせることができて
// Promise の処理が失敗した時に例外をキャッチ出来る

// ちなみに .catch(console.error) は省略型
//    .catch(function (err) { console.error(err) })
// => .catch((err) => { console.error(err) }) [arrow function 化]
// => .catch(err => console.error(err) ) [引数が1つのときは()を省略でき, 1行の処理で{}を省略出来る]
// => .catch(console.error) [callback引数の引数はcallbackの引数に渡される]

app.ports.toJS.subscribe(async (data) => {
  console.log(data)
  // gRPC Unary RPCs
  const req = new pb.SayRequest()
  req.setUserId(data.name)
  req.setMessage(data.message)
  // array には レスポンスが配列型で順番に入っている(なぜ配列)
  const { array } = await client.Say(req)
    .catch(console.error)

  // 直す
  const res = { name: array[0], message: array[1] }

  // Elm の toElm を呼ぶ
  app.ports.toElm.send(res)
})

// ログインボタンが押されて Elm から呼ばれる
// function の引数で {} を使うと引数の中のプロパティをとれる
// Ex:
// const v = { a: 1, b: 2 }
// const { a, b } = v // a == 1, b == 2
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

app.ports.requestResult.subscribe(async () => {
  const result = { score: 1, userId: 'hoge' }
  app.ports.resultCallback.send(result)
})

// Server Streaming RPCs
const req = new pb.HelloRequest()
req.setName('Server Stream')
req.setMessage('Hello')
const stream = client.sayHelloManyTimes(req)

stream.on('data', res => {
  const message = res.getMessage()
  const name = req.getName()
  // JS -> Elm
  app.ports.toElm.send({ message, name })
})

stream.on('status', status => {
  console.log(status)
})