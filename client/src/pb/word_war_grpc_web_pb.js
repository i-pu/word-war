/**
 * @fileoverview gRPC-Web generated client stub for word_war
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.word_war = require('./word_war_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.word_war.WordWarClient =
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
proto.word_war.WordWarPromiseClient =
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
 *   !proto.word_war.MatchingRequest,
 *   !proto.word_war.MatchingResponse>}
 */
const methodDescriptor_WordWar_Matching = new grpc.web.MethodDescriptor(
  '/word_war.WordWar/Matching',
  grpc.web.MethodType.UNARY,
  proto.word_war.MatchingRequest,
  proto.word_war.MatchingResponse,
  /** @param {!proto.word_war.MatchingRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.word_war.MatchingResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.word_war.MatchingRequest,
 *   !proto.word_war.MatchingResponse>}
 */
const methodInfo_WordWar_Matching = new grpc.web.AbstractClientBase.MethodInfo(
  proto.word_war.MatchingResponse,
  /** @param {!proto.word_war.MatchingRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.word_war.MatchingResponse.deserializeBinary
);


/**
 * @param {!proto.word_war.MatchingRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.word_war.MatchingResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.word_war.MatchingResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.word_war.WordWarClient.prototype.matching =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/word_war.WordWar/Matching',
      request,
      metadata || {},
      methodDescriptor_WordWar_Matching,
      callback);
};


/**
 * @param {!proto.word_war.MatchingRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.word_war.MatchingResponse>}
 *     A native promise that resolves to the response
 */
proto.word_war.WordWarPromiseClient.prototype.matching =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/word_war.WordWar/Matching',
      request,
      metadata || {},
      methodDescriptor_WordWar_Matching);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.word_war.SayRequest,
 *   !proto.word_war.SayResponse>}
 */
const methodDescriptor_WordWar_Say = new grpc.web.MethodDescriptor(
  '/word_war.WordWar/Say',
  grpc.web.MethodType.UNARY,
  proto.word_war.SayRequest,
  proto.word_war.SayResponse,
  /** @param {!proto.word_war.SayRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.word_war.SayResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.word_war.SayRequest,
 *   !proto.word_war.SayResponse>}
 */
const methodInfo_WordWar_Say = new grpc.web.AbstractClientBase.MethodInfo(
  proto.word_war.SayResponse,
  /** @param {!proto.word_war.SayRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.word_war.SayResponse.deserializeBinary
);


/**
 * @param {!proto.word_war.SayRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.word_war.SayResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.word_war.SayResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.word_war.WordWarClient.prototype.say =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/word_war.WordWar/Say',
      request,
      metadata || {},
      methodDescriptor_WordWar_Say,
      callback);
};


/**
 * @param {!proto.word_war.SayRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.word_war.SayResponse>}
 *     A native promise that resolves to the response
 */
proto.word_war.WordWarPromiseClient.prototype.say =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/word_war.WordWar/Say',
      request,
      metadata || {},
      methodDescriptor_WordWar_Say);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.word_war.GameRequest,
 *   !proto.word_war.GameResponse>}
 */
const methodDescriptor_WordWar_Game = new grpc.web.MethodDescriptor(
  '/word_war.WordWar/Game',
  grpc.web.MethodType.SERVER_STREAMING,
  proto.word_war.GameRequest,
  proto.word_war.GameResponse,
  /** @param {!proto.word_war.GameRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.word_war.GameResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.word_war.GameRequest,
 *   !proto.word_war.GameResponse>}
 */
const methodInfo_WordWar_Game = new grpc.web.AbstractClientBase.MethodInfo(
  proto.word_war.GameResponse,
  /** @param {!proto.word_war.GameRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.word_war.GameResponse.deserializeBinary
);


/**
 * @param {!proto.word_war.GameRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.word_war.GameResponse>}
 *     The XHR Node Readable Stream
 */
proto.word_war.WordWarClient.prototype.game =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/word_war.WordWar/Game',
      request,
      metadata || {},
      methodDescriptor_WordWar_Game);
};


/**
 * @param {!proto.word_war.GameRequest} request The request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.word_war.GameResponse>}
 *     The XHR Node Readable Stream
 */
proto.word_war.WordWarPromiseClient.prototype.game =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/word_war.WordWar/Game',
      request,
      metadata || {},
      methodDescriptor_WordWar_Game);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.word_war.ResultRequest,
 *   !proto.word_war.ResultResponse>}
 */
const methodDescriptor_WordWar_Result = new grpc.web.MethodDescriptor(
  '/word_war.WordWar/Result',
  grpc.web.MethodType.UNARY,
  proto.word_war.ResultRequest,
  proto.word_war.ResultResponse,
  /** @param {!proto.word_war.ResultRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.word_war.ResultResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.word_war.ResultRequest,
 *   !proto.word_war.ResultResponse>}
 */
const methodInfo_WordWar_Result = new grpc.web.AbstractClientBase.MethodInfo(
  proto.word_war.ResultResponse,
  /** @param {!proto.word_war.ResultRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.word_war.ResultResponse.deserializeBinary
);


/**
 * @param {!proto.word_war.ResultRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.word_war.ResultResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.word_war.ResultResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.word_war.WordWarClient.prototype.result =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/word_war.WordWar/Result',
      request,
      metadata || {},
      methodDescriptor_WordWar_Result,
      callback);
};


/**
 * @param {!proto.word_war.ResultRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.word_war.ResultResponse>}
 *     A native promise that resolves to the response
 */
proto.word_war.WordWarPromiseClient.prototype.result =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/word_war.WordWar/Result',
      request,
      metadata || {},
      methodDescriptor_WordWar_Result);
};


module.exports = proto.word_war;

