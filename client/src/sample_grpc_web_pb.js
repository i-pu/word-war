/**
 * @fileoverview gRPC-Web generated client stub for sample
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.sample = require('./sample_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.sample.GreeterClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.sample.GreeterPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.sample.HelloRequest,
 *   !proto.sample.HelloResponse>}
 */
const methodDescriptor_Greeter_SayHello = new grpc.web.MethodDescriptor(
  '/sample.Greeter/SayHello',
  grpc.web.MethodType.UNARY,
  proto.sample.HelloRequest,
  proto.sample.HelloResponse,
  /** @param {!proto.sample.HelloRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.sample.HelloResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.sample.HelloRequest,
 *   !proto.sample.HelloResponse>}
 */
const methodInfo_Greeter_SayHello = new grpc.web.AbstractClientBase.MethodInfo(
  proto.sample.HelloResponse,
  /** @param {!proto.sample.HelloRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.sample.HelloResponse.deserializeBinary
);


/**
 * @param {!proto.sample.HelloRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.sample.HelloResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.sample.HelloResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.sample.GreeterClient.prototype.sayHello =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/sample.Greeter/SayHello',
      request,
      metadata || {},
      methodDescriptor_Greeter_SayHello,
      callback);
};


/**
 * @param {!proto.sample.HelloRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.sample.HelloResponse>}
 *     A native promise that resolves to the response
 */
proto.sample.GreeterPromiseClient.prototype.sayHello =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/sample.Greeter/SayHello',
      request,
      metadata || {},
      methodDescriptor_Greeter_SayHello);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.sample.HelloRequest,
 *   !proto.sample.HelloResponse>}
 */
const methodDescriptor_Greeter_SayHelloManyTimes = new grpc.web.MethodDescriptor(
  '/sample.Greeter/SayHelloManyTimes',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.sample.HelloRequest,
  proto.sample.HelloResponse,
  /** @param {!proto.sample.HelloRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.sample.HelloResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.sample.HelloRequest,
 *   !proto.sample.HelloResponse>}
 */
const methodInfo_Greeter_SayHelloManyTimes = new grpc.web.AbstractClientBase.MethodInfo(
  proto.sample.HelloResponse,
  /** @param {!proto.sample.HelloRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.sample.HelloResponse.deserializeBinary
);


/**
 * @param {!proto.sample.HelloRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.sample.HelloResponse>}
 *     The XHR Node Readable Stream
 */
proto.sample.GreeterClient.prototype.sayHelloManyTimes =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/sample.Greeter/SayHelloManyTimes',
      request,
      metadata || {},
      methodDescriptor_Greeter_SayHelloManyTimes);
};


/**
 * @param {!proto.sample.HelloRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.sample.HelloResponse>}
 *     The XHR Node Readable Stream
 */
proto.sample.GreeterPromiseClient.prototype.sayHelloManyTimes =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/sample.Greeter/SayHelloManyTimes',
      request,
      metadata || {},
      methodDescriptor_Greeter_SayHelloManyTimes);
};


module.exports = proto.sample;

