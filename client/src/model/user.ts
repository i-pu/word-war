export interface User {
  userId: string
  history: Array<{ date: Date, rating: number }>
  name: string
  rating: number
}

export const defaultUser = (): User => {
  return {
    userId: 'mock-userid',
    history: [{ date: new Date(), rating: 1500 }],
    name: 'mockuser',
    rating: 1500
  }
}