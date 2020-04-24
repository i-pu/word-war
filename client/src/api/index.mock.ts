/**
 *  Word War API
 */
import { User, ServerStatus, defaultUser } from '@/model'

/**
 * Health Check Mock
 */
export const healthCheck = async (): Promise<ServerStatus> => {
  return {
    active: true,
    version: '1.7',
  }
}

// TODO: まともなデータを入れる
export const getUserdata = async (uid: string): Promise<User> => {
  return defaultUser()
}

/**
 *
 */
export const setInitialUserdata = async (uid: string): Promise<void> => {
  return
}

/**
 *  Sign in Mock
 */
export const signIn = async ({ email, password }: { email: string, password: string }): Promise<User> => {
  return await getUserdata('')
}

/**
 *  Sign up Mock
 */
export const signUp = async ({ email, password }: { email: string, password: string }): Promise<User> => {
  // TODO: まともなデータを入れる
  return defaultUser()
}

/**
 *  Sign out Mock
 */
export const signOut = async (): Promise<void> => {
  return
}
