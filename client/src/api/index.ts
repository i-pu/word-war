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
 *
 */
export const setInitialUserdata = async (uid: string): Promise<void> => {
  return firebase
    .firestore()
    .collection('users')
    .doc(uid)
    .set(defaultUser())
}

export const getUserdata = async (uid: string): Promise<User> => {
  const ss = await firebase
    .firestore()
    .collection('users')
    .doc(uid)
    .get()

  // TODO: validation
  const { history, name, rating } = ss.data()!

  return {
    userId: uid,
    history,
    name,
    rating
  }
}

/**
 *  Sign in
 */
export const signIn = async ({ email, password }: { email: string, password: string }): Promise<User> => {
  console.log(`signIn: ${email}, ${password}`)

  const result = await firebase
    .auth()
    .signInWithEmailAndPassword(email, password)
    .catch(console.error)

  if (!result || !result.user) {
    throw new Error(`can't authorized: ${email}, ${password}`)
  }

  return await getUserdata(result.user.uid)
}


/**
 *  Sign iu
 */
export const signUp = async ({ email, password }: { email: string, password: string }): Promise<User> => {
  // create user
  const result = await firebase
    .auth()
    .createUserWithEmailAndPassword(email, password)
    .catch(console.error)

  if (!result || !result.user) {
    throw new Error(`can't authorized: ${email}, ${password}`)
  }

  setInitialUserdata(result.user.uid)

  return getUserdata(result.user.uid)
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
