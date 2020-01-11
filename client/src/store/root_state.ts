export interface User {
  userId: string
}

export interface RootState {
  version: string
  user: User
  serverHealth: boolean
}
