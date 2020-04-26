<template>
  <div>
    <LoginForm />
  </div>
</template>

<script lang="ts">
import LoginForm from '@/components/LoginForm.vue'
import { defineComponent } from '@vue/composition-api'

export default defineComponent({
  components: { LoginForm },
})
</script>

<style scoped>
.img-responsive {
  display: block;
  width: auto;
  min-height: 100vh;
}
</style>