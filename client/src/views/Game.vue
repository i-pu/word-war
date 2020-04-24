<template>
  <div class="game">
    <Hero />
    <Navbar />
    <Siritori />c
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@vue/composition-api'
import Siritori from '@/components/Siritori.vue'
import Hero from '@/components/Hero.vue'

export default defineComponent({
  components: {
    Siritori,
    Hero
  }
})
</script>
