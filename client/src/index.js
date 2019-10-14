import firebase from 'firebase'
import credential from '../credential'

// firebase の初期化
firebase.initializeApp(credential)

// Parcel を使っているので直接 import できる
const { Elm } = require('./main.elm')

// Elm を初期化して #app に配置
const app = Elm.Main.init({
  node: document.getElementById('app'),
  // Elm flags
  flags: "Hoge"
})

export default {
  ports: app.ports,
  firebase
}
