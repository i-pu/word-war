<template>
  <div class="top">
    <Hero />
    <Login />
  </div>
</template>

<script lang="ts">
import Login from '@/components/Login.vue'
import Hero from '@/components/Hero.vue'
import { defineComponent, SetupContext, onMounted } from '@vue/composition-api'
import { useRootStore } from '@/store'

export default defineComponent({
  components: {
    Hero,
    Login
  },
  setup(props:{}, {root}: SetupContext) {
    const {healthCheck} = useRootStore()
    onMounted(() => {
      healthCheck()
    })
    return {}
  }
})
</script>
