/**
 *  Word War API
 */
import { WordWarPromiseClient } from '@/pb/word_war_grpc_web_pb'
import { HealthCheckRequest, HealthCheckResponse } from '@/pb/word_war_pb'
import firebase from 'firebase'
import { User, ServerStatus, defaultUser } from '@/model'

/**
 *  Health Check
 */
export const healthCheck = async (): Promise<ServerStatus> => {
  const client = new WordWarPromiseClient(process.env.VUE_APP_API_ENDPOINT)
  const req = new HealthCheckRequest()
  const res = await client.healthCheck(req)
  return {
    active: res.getActive(),
    version: res.getServerversion(),
  }
}

/**
 * ユーザデータの初期化をする
 */
// Firebase Auth のデータでユーザー作成
export const registerUser = async (fbUser: firebase.User): Promise<void> => {
  const avatarUrl = fbUser.photoURL || ''
  const name = fbUser.displayName || 'ななし'
  console.log('ゲームデータを作成します')

  const user: User = {
    userId: fbUser.uid,
    name,
    avatarUrl,
    history: [],
    rating: 1500
  }

  return firebase
    .firestore()
    .collection('users')
    .doc(user.userId)
    .set(user)
    .catch(console.error)
}

export const existUserData = async (uid: string): Promise<boolean> => {
  return (await firebase.firestore().collection('users').doc(uid).get()).exists
}

export const getUserdata = async (uid: string): Promise<User> => {
  const ss = await firebase
    .firestore()
    .collection('users')
    .doc(uid)
    .get()

  // TODO: validation
  const { history, name, rating, avatarUrl } = ss.data()!

  console.log('ユーザーデータを習得しました')

  return {
    userId: uid,
    avatarUrl,
    history,
    name,
    rating
  }
}

/**
 *  Sign in
 */
export const signIn = async ({ email, password }: { email: string, password: string }): Promise<void> => {
  console.log(`signIn: ${email}, ${password}`)

  const result = await firebase
    .auth()
    .signInWithEmailAndPassword(email, password)
    .catch(console.error)

  if (!result || !result.user) {
    throw new Error(`can't authorized: ${email}, ${password}`)
  }

  console.log('Firebase Auth にサインインしました')
}


/**
 *  Sign iu
 */
export const signUp = async ({ email, password }: { email: string, password: string }): Promise<void> => {
  // create user
  const result = await firebase
    .auth()
    .createUserWithEmailAndPassword(email, password)
    .catch(console.error)

  if (!result || !result.user) {
    throw new Error(`can't authorized: ${email}, ${password}`)
  }

  console.log('Firebase Auth アカウントを作成しました')
}

/**
 *  Sign up with Twitter Mock
 */
export const signUpWithTwitter = async (): Promise<User> => {
  throw 'Not Implemented'
}

/**
 *  Sign up with GitHub Mock
 */
export const signUpWithGitHub = async (): Promise<User> => {
  throw 'Not Implemented'
}

export const signUpWithGoogle = async () => {
  var provider = new firebase.auth.GoogleAuthProvider()
  firebase.auth().useDeviceLanguage();
  firebase.auth().signInWithRedirect(provider).catch(console.error)
  return
}


/**
 *  Sign out
 */
export const signOut = async (): Promise<void> => {
  // create user
  await firebase
    .auth()
    .signOut()
    .catch(console.error)
}
