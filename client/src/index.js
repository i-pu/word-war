import pb from './sample_grpc_web_pb'
const { Elm } = require('./main.elm')

const endpoint = 'http://localhost:50051'
const client = new pb.GreeterPromiseClient(endpoint)
const app = Elm.Main.init({
  node: document.getElementById('app'),
  // Elmにわたすもの
  flags: "Hoge"
})

// Elm -> JS
app.ports.toJS.subscribe(async data => {
  console.log(data)
  // gRPC Unary RPCs
  const req = new pb.HelloRequest()
  req.setName('Unary RPC')
  const { array } = await client.sayHello(req)
    .catch(console.error)

  // JS -> Elm
  app.ports.toElm.send(array[0])
})

// Server Streaming RPCs
const req = new pb.HelloRequest()
req.setName('Server Stream')
const stream = client.sayHelloManyTimes(req)

stream.on('data', res => {
  const data = res.getMessage()
  // JS -> Elm
  app.ports.toElm.send(data)
})

stream.on('status', status => {
  console.log(status)
})
