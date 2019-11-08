<template>
  <div class="container">
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
      <b-button type="is-dark" @click="signUp">Sign up</b-button>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import firebase from '@/config/firebase'

@Component
export default class Login extends Vue {
  private email: string = ''
  private password: string = ''

  private async signIn() {
    const result = await firebase
      .auth()
      .signInWithEmailAndPassword(this.email, this.password)
      .catch(console.error)

    if (!result || !result.user) {
      throw new Error(`can't authorized: ${this.email}, ${this.password}`)
    }

    this.$store.commit('user/setUid', { uid: result.user.uid })
    this.$router.push('/home')

    console.log(`signIn: ${this.email}, ${this.password}`)
  }

  private async signUp() {
    // create user
    const result = await firebase
      .auth()
      .createUserWithEmailAndPassword(this.email, this.password)
      .catch(console.error)

    if (!result || !result.user) {
      throw new Error(`can't authorized: ${this.email}, ${this.password}`)
    }

    this.$store.commit('user/setUid', { uid: result.user.uid })
    this.$router.push('/home')

    console.log(`signUp: ${this.email}, ${this.password}`)
  }
}
</script>
