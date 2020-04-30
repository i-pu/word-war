<template>
  <div>
    <LoginForm />
  </div>
</template>

<script lang="ts">
import firebase from 'firebase'
import LoginForm from '@/components/LoginForm.vue'
import { defineComponent, onMounted } from '@vue/composition-api'
import { existUserData, registerUser } from '@/api/index'
import { useRootStore } from '@/store'

export default defineComponent({
  components: { LoginForm },
  setup(props: {}, { root }) {
    const { user, onAuthChanged } = useRootStore()

    onMounted(() => {
      firebase.auth().onAuthStateChanged(async auth => {
        console.log(`auth state changed user.value: ${user.value}, auth: ${auth}`)
        if (auth) {
          onAuthChanged(auth)
          console.log("/homeに自動遷移する")
          root.$router.push('/home')
        } else {
          console.log("自動遷移しない")
        }
      })
    })
  }
})
</script>

<style scoped>
.img-responsive {
  display: block;
  width: auto;
  min-height: 100vh;
}
</style>