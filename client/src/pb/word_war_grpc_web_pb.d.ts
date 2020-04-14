import * as grpcWeb from 'grpc-web';

import {
  GameRequest,
  GameResponse,
  HealthCheckRequest,
  HealthCheckResponse,
  MatchingRequest,
  MatchingResponse,
  ResultRequest,
  ResultResponse,
  SayRequest,
  SayResponse} from './word_war_pb';

export class WordWarClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  matching(
    request: MatchingRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<MatchingResponse>;

  say(
    request: SayRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: SayResponse) => void
  ): grpcWeb.ClientReadableStream<SayResponse>;

  game(
    request: GameRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<GameResponse>;

  result(
    request: ResultRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ResultResponse) => void
  ): grpcWeb.ClientReadableStream<ResultResponse>;

  healthCheck(
    request: HealthCheckRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: HealthCheckResponse) => void
  ): grpcWeb.ClientReadableStream<HealthCheckResponse>;

}

export class WordWarPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  matching(
    request: MatchingRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<MatchingResponse>;

  say(
    request: SayRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<SayResponse>;

  game(
    request: GameRequest,
    metadata?: grpcWeb.Metadata
  ): grpcWeb.ClientReadableStream<GameResponse>;

  result(
    request: ResultRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ResultResponse>;

  healthCheck(
    request: HealthCheckRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<HealthCheckResponse>;

}

