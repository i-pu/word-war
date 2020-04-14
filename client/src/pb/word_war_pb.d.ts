import * as jspb from "google-protobuf"

export class SayRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getMessage(): string;
  setMessage(value: string): void;

  getRoomid(): string;
  setRoomid(value: string): void;

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
    roomid: string,
  }
}

export class SayResponse extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getMessage(): string;
  setMessage(value: string): void;

  getRoomid(): string;
  setRoomid(value: string): void;

  getValid(): boolean;
  setValid(value: boolean): void;

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
    roomid: string,
    valid: boolean,
  }
}

export class MatchingRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchingRequest.AsObject;
  static toObject(includeInstance: boolean, msg: MatchingRequest): MatchingRequest.AsObject;
  static serializeBinaryToWriter(message: MatchingRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchingRequest;
  static deserializeBinaryFromReader(message: MatchingRequest, reader: jspb.BinaryReader): MatchingRequest;
}

export namespace MatchingRequest {
  export type AsObject = {
    userid: string,
  }
}

export class MatchingResponse extends jspb.Message {
  getRoomid(): string;
  setRoomid(value: string): void;

  getUserList(): Array<User>;
  setUserList(value: Array<User>): void;
  clearUserList(): void;
  addUser(value?: User, index?: number): User;

  getRoomuserlimit(): number;
  setRoomuserlimit(value: number): void;

  getTimerseconds(): number;
  setTimerseconds(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MatchingResponse.AsObject;
  static toObject(includeInstance: boolean, msg: MatchingResponse): MatchingResponse.AsObject;
  static serializeBinaryToWriter(message: MatchingResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MatchingResponse;
  static deserializeBinaryFromReader(message: MatchingResponse, reader: jspb.BinaryReader): MatchingResponse;
}

export namespace MatchingResponse {
  export type AsObject = {
    roomid: string,
    userList: Array<User.AsObject>,
    roomuserlimit: number,
    timerseconds: number,
  }
}

export class User extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    userid: string,
  }
}

export class GameRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getRoomid(): string;
  setRoomid(value: string): void;

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
    roomid: string,
  }
}

export class GameResponse extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getMessage(): string;
  setMessage(value: string): void;

  getRoomid(): string;
  setRoomid(value: string): void;

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
    roomid: string,
  }
}

export class ResultRequest extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getRoomid(): string;
  setRoomid(value: string): void;

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
    roomid: string,
  }
}

export class ResultResponse extends jspb.Message {
  getUserid(): string;
  setUserid(value: string): void;

  getScore(): string;
  setScore(value: string): void;

  getRoomid(): string;
  setRoomid(value: string): void;

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
    roomid: string,
  }
}

export class HealthCheckRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HealthCheckRequest.AsObject;
  static toObject(includeInstance: boolean, msg: HealthCheckRequest): HealthCheckRequest.AsObject;
  static serializeBinaryToWriter(message: HealthCheckRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HealthCheckRequest;
  static deserializeBinaryFromReader(message: HealthCheckRequest, reader: jspb.BinaryReader): HealthCheckRequest;
}

export namespace HealthCheckRequest {
  export type AsObject = {
  }
}

export class HealthCheckResponse extends jspb.Message {
  getActive(): boolean;
  setActive(value: boolean): void;

  getServerversion(): string;
  setServerversion(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HealthCheckResponse.AsObject;
  static toObject(includeInstance: boolean, msg: HealthCheckResponse): HealthCheckResponse.AsObject;
  static serializeBinaryToWriter(message: HealthCheckResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HealthCheckResponse;
  static deserializeBinaryFromReader(message: HealthCheckResponse, reader: jspb.BinaryReader): HealthCheckResponse;
}

export namespace HealthCheckResponse {
  export type AsObject = {
    active: boolean,
    serverversion: string,
  }
}

