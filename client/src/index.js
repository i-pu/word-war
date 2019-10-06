import pb from './sample_grpc_web_pb'
import firebase from 'firebase'
import auth from 'firebase/auth'
import credential from '../credential'

firebase.initializeApp(credential)

const { Elm } = require('./main.elm')

const endpoint = 'http://localhost:50051'
const client = new pb.GreeterPromiseClient(endpoint)
const app = Elm.Main.init({
  node: document.getElementById('app'),
  // Elmにわたすもの
  flags: "Hoge"
})

// Elm -> JS
app.ports.toJS.subscribe(async (data) => {
  console.log(data)
  // gRPC Unary RPCs
  const req = new pb.HelloRequest()
  req.setName(data.name)
  req.setMessage(data.message)
  const { array } = await client.sayHello(req)
    .catch(console.error)

  // JS -> Elm
  app.ports.toElm.send({ name: array[0], message: array[1] })
})

app.ports.signinWithFirebase.subscribe(async ({ email, password }) => {
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