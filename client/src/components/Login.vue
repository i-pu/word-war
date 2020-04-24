<template>
  <div class="container">
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
      <b-button @click="signIn">Sign in</b-button>
      <b-button type="is-dark" @click="signUp">Sign up</b-button>
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

    const signIn = async () => {
      await store.signIn(form)
      root.$router.push('/home')
    }

    const signUp = async () => {
      await store.signUp(form)
      root.$router.push('/home')
    }
    return {
      form,
      signIn,
      signUp,
    }
  }
})
</script>
