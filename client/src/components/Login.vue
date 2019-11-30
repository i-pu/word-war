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

@Component
export default class Login extends Vue {
  private email: string = ''
  private password: string = ''

  // TODO: signInとsignUpをstoreに移動
  private async signIn() {
    await this.$store.dispatch('signIn', {
      email: this.email,
      password: this.password
    })
    this.$router.push('/home')
  }

  private async signUp() {
    await this.$store.dispatch('signUp', {
      email: this.email,
      password: this.password
    })
    this.$router.push('/home')
  }
}
</script>
