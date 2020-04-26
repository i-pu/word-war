export interface User {
  userId: string
  avatarUrl: string
  history: Array<{ date: Date, rating: number }>
  name: string
  rating: number
}

const dummyRatings = (): Array<{ date: Date, rating: number }> => {
  const ratings: Array<{ date: Date, rating: number }> = []
  for (let d = new Date(2020, 4 - 1, 10); d <= new Date(); d.setDate(d.getDate() + 1)) {
    if (Math.random() < 0.5) {
      ratings.push({ date: new Date(d), rating: (15.0 + parseFloat((Math.random() * 10.0 - 5.0).toFixed(2)))})
    }
  }
  return ratings
}

export const defaultUser = (): User => {
  return {
    userId: 'mock-userid',
    avatarUrl: 'https://img.atcoder.jp/icons/267f5de4d8768543b1570f07e47b5316.jpg',
    history: dummyRatings(),
    name: 'John Smith',
    rating: 16.75,
  }
}