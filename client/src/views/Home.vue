<template>
  <div class="home">
    <Hero />
    <b-button type="is-link" tag="router-link" to="/game">Game Start</b-button>

    <b-button type="is-dark" @click="onClick">vuex</b-button>
    <p>{{ count }}</p>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Hero from '@/components/Hero.vue'

@Component({
  components: {
    Hero
  }
})
export default class Home extends Vue {
  get count(): number {
    return this.$store.state.sample.count
  }
  private onClick() {
    this.$store.commit('sample/increment')
    console.log(this.$store.state.sample.count)
  }
}
</script>
