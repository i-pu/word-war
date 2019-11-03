import * as jspb from "google-protobuf"

export class SayRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SayRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SayRequest): SayRequest.AsObject;
  static serializeBinaryToWriter(message: SayRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SayRequest;
  static deserializeBinaryFromReader(message: SayRequest, reader: jspb.BinaryReader): SayRequest;
}

export namespace SayRequest {
  export type AsObject = {
    userid: string,
    message: string,
  }
}

export class SayResponse extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SayResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SayResponse): SayResponse.AsObject;
  static serializeBinaryToWriter(message: SayResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SayResponse;
  static deserializeBinaryFromReader(message: SayResponse, reader: jspb.BinaryReader): SayResponse;
}

export namespace SayResponse {
  export type AsObject = {
    userid: string,
    message: string,
  }
}

export class GameRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GameRequest): GameRequest.AsObject;
  static serializeBinaryToWriter(message: GameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GameRequest;
  static deserializeBinaryFromReader(message: GameRequest, reader: jspb.BinaryReader): GameRequest;
}

export namespace GameRequest {
  export type AsObject = {
    userid: string,
  }
}

export class GameResponse extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GameResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GameResponse): GameResponse.AsObject;
  static serializeBinaryToWriter(message: GameResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GameResponse;
  static deserializeBinaryFromReader(message: GameResponse, reader: jspb.BinaryReader): GameResponse;
}

export namespace GameResponse {
  export type AsObject = {
    userid: string,
    message: string,
  }
}

export class ResultRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResultRequest): ResultRequest.AsObject;
  static serializeBinaryToWriter(message: ResultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResultRequest;
  static deserializeBinaryFromReader(message: ResultRequest, reader: jspb.BinaryReader): ResultRequest;
}

export namespace ResultRequest {
  export type AsObject = {
    userid: string,
  }
}

export class ResultResponse extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getScore(): string;
  setScore(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResultResponse): ResultResponse.AsObject;
  static serializeBinaryToWriter(message: ResultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResultResponse;
  static deserializeBinaryFromReader(message: ResultResponse, reader: jspb.BinaryReader): ResultResponse;
}

export namespace ResultResponse {
  export type AsObject = {
    userid: string,
    score: string,
  }
}

