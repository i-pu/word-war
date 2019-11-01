import firebase from "firebase/app";
import "firebase/auth";
import "firebase/firestore";

const config = {
  apiKey: "AIzaSyAy9nNPPNOyzYkookvIHPiu68IlYazhA4c",
  authDomain: "word-war-9e392.firebaseapp.com",
  databaseURL: "https://word-war-9e392.firebaseio.com",
  projectId: "word-war-9e392",
  storageBucket: "",
  messagingSenderId: "973244315901",
  appId: "1:973244315901:web:3a1d5260f030c0fc8ae694",
  measurementId: "G-84Y9JV8623"
};

firebase.initializeApp(config);

export default firebase;
