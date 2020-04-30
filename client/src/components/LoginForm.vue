<template>
  <div class="container">
    <div class="container">
      <h2 class="title has-text-left">
        Word War {{ status.version }}
      </h2>
    </div>
    <b-field label="email">
      <b-input
        v-model="form.email"
        type="email"
        icon="email"
        placeholder="Input your email"
      >
      </b-input>
    </b-field>

    <b-field label="password">
      <b-input
        v-model="form.password"
        type="password"
        icon="lock"
        placeholder="Input your password"
      >
      </b-input>
    </b-field>
    <div class="buttons">
      <b-button @click="onClickSignIn">Sign in</b-button>
      <b-button type="is-dark" @click="onClickSignUp">Sign up</b-button>
    </div>

    <div class="buttons">
      <button @click="onClickSignInWithGoogle" class="button is-danger is-rounded">Sign up with Google</button>
      <!-- <button @click="signUpWithGitHub" class="button is-black is-rounded">Sign up with GitHub</button>
      <button @click="signUpWithTwitter" class="button is-info is-rounded">Sign up with Twitter</button> -->
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive } from '@vue/composition-api'
import { useRootStore } from '@/store/index'

export default defineComponent({
  setup(props: {}, { root }) {
    const store = useRootStore()
    const form = reactive({
      email: '',
      password: '',
    })

    const onClickSignIn = async () => {
      await store.signIn(form)
      // root.$router.push('/home')
    }

    const onClickSignInWithGoogle = async () => {
      await store.google()
    }

    // const signUpWithTwitter = async () => {
    //   await store.signUpWithTwitter()
    // }

    // const signUpWithGitHub = async () => {
    //   await store.signUpWithGitHub()
    // }

    const onClickSignUp = async () => {
      await store.signUp(form)
    }
    return {
      status: store.status,
      form,
      onClickSignIn,
      onClickSignUp,
      onClickSignInWithGoogle,
      // signUpWithTwitter,
      // signUpWithGitHub,
    }
  }
})
</script>
