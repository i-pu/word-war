<template>
  <div class="game">
    <Hero />
    <Navbar />
    <Siritori />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Siritori from '@/components/Siritori.vue'
import Hero from '@/components/Hero.vue'

@Component({
  components: {
    Siritori,
    Hero
  }
})
export default class Game extends Vue {}
</script>
