<template>
  <div>
    <b-field label="email">
      <b-input
        v-model="email"
        type="email"
        icon="email"
        placeholder="Input your email"
      >
      </b-input>
    </b-field>

    <b-field label="password">
      <b-input
        v-model="password"
        type="password"
        icon="lock"
        placeholder="Input your password"
      >
      </b-input>
    </b-field>
    <div class="buttons">
      <b-button @click="signIn">Sign in</b-button>
      <b-button @click="signUp" type="is-dark">Sign up</b-button>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import firebase from "@/config/firebase";

@Component
export default class Login extends Vue {
  private email: string;
  private password: string;
  private uid: string;

  constructor() {
    super();
    this.email = "";
    this.password = "";
    this.uid = "";
  }

  private async signIn() {
    await firebase
      .auth()
      .signInWithEmailAndPassword(this.email, this.password)
      .then((result: firebase.auth.UserCredential | null) => {
        if (result != null && result.user != null) {
          this.uid = result.user.uid;
          this.$router.push("/home");
        } else {
          throw new Error(`can't authorized: ${this.email}, ${this.password}`);
        }
      })
      .catch(error => {
        console.log(error);
      });
    console.log(`signIn: ${this.email}, ${this.password}`);
  }

  private async signUp() {
    // create user
    await firebase
      .auth()
      .createUserWithEmailAndPassword(this.email, this.password)
      .then((result: firebase.auth.UserCredential | null) => {
        if (result != null && result.user != null) {
          this.uid = result.user.uid;
          this.$router.push("/home");
        } else {
          throw new Error(`can't authorized: ${this.email}, ${this.password}`);
        }
      })
      .catch(error => {
        console.log(error);
      });
    console.log(`signUp: ${this.email}, ${this.password}`);
  }
}
</script>
