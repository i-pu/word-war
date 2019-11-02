import * as grpcWeb from 'grpc-web';

import {
  GameRequest,
  GameResponse,
  ResultRequest,
  ResultResponse,
  SayRequest,
  SayResponse} from './word_war_pb';

export class WordWarClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

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

}

export class WordWarPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

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

}

