<template>
  <div class="result">
    <Hero />
    <Navbar />
    <section class="section">
      <p>result: {{ $store.getters['game/score'] }}</p>
    </section>
  </div>
</template>

<script lang="ts">
import { Component, Watch, Vue } from 'vue-property-decorator'
import Hero from '@/components/Hero.vue'
import Navbar from '@/components/Navbar.vue'

@Component({
  components: {
    Hero,
    Navbar
  }
})
export default class Result extends Vue {
  async mounted() {
    await this.$store.dispatch('game/result')
    console.log('score', this.$store.getters['game/score'])
  }
}
</script>
