/**
 *  Word War API
 */
import { User, ServerStatus, defaultUser } from '@/model'
import firebase from 'firebase'

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

// /**
//  *  Sign up with Google Mock
//  */
// export const signUpWithGoogle = async (): Promise<User> => {
//   // TODO: まともなデータを入れる
//   return defaultUser()
// }

/**
 *  Sign up with Twitter Mock
 */
export const signUpWithTwitter = async (): Promise<User> => {
  // TODO: まともなデータを入れる
  return defaultUser()
}

/**
 *  Sign up with GitHub Mock
 */
export const signUpWithGitHub = async (): Promise<User> => {
  // TODO: まともなデータを入れる
  return defaultUser()
}

/**
 *  Sign up with Google Mock
 */
export const signUpWithGoogle = async (): Promise<User> => {
  var provider = new firebase.auth.GoogleAuthProvider()
  firebase.auth().useDeviceLanguage();
  // TODO: signInWithRedirectはむずかしい
  // firebase.auth().signInWithRedirect(provider).catch(console.error)
  const result = await firebase.auth().signInWithPopup(provider).catch(console.error)
  if (!result || !result.user) {
    throw 'WTF'
  }
  console.log(`result: ${result}`)
  await setInitialUserdata(result.user.uid)
  return await getUserdata(result.user.uid)
}

/**
 *  Sign out Mock
 */
export const signOut = async (): Promise<void> => {
  await firebase
    .auth()
    .signOut()
    .catch(console.error)
}
