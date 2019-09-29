import sample from './sample_grpc_web_pb'

const greeter = new sample.GreeterPromiseClient('http://localhost:50051')
const req = new sample.HelloRequest()
req.setName('World')
greeter.sayHello(req)
    .then((res) => {
            console.log(res)
            console.log(res.array)
        })
    .catch((err) => {
            console.error(err)
        })

const stream = greeter.sayHelloManyTimes(req)

stream.on('data', function(response) {
  console.log(response.getMessage());
});
stream.on('status', function(status) {
  console.log(status.code);
  console.log(status.details);
  console.log(status.metadata);
});
stream.on('end', function(end) {
  // stream end signal
});