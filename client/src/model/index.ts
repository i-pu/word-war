/**
 *  Types of UI, Model
 */
export { User, defaultUser } from './user'
export interface ServerStatus {
  active: boolean
  version: string
}

export enum Scene {
  None,
  // 待機中
  Matching,
  // ゲーム中
  Gaming,
  // 終了
  End
}

export interface Word {
  id: string
  text: string
}

export interface RoomInfo {
  players: Player[]
  timer: number
  limit: number
  roomId: string
}

export interface Player {
  playerId: string
}