<template>
  <div class="waiting">
    <Hero />

    <b-loading is-full-page :active="loading" />

    <h1>待機中 scene: {{ scene }}. {{ users }}..</h1>

    <p v-for="user in users" :key="user.playerId">
      {{ user }}
    </p>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import Hero from '@/components/Hero.vue'
import Navbar from '@/components/Navbar.vue'
import { Scene, User } from '@/store/game'

@Component({
  components: {
    Hero,
    Navbar
  }
})
export default class Waiting extends Vue {
  private loading: boolean = true

  // computed
  private get scene(): Scene {
    return this.$store.getters['game/scene']
  }

  private get users(): User[] {
    return this.$store.getters['game/users']
  }

  @Watch('scene')
  sceneChanged() {
    console.log(`sceneChanged ${this.scene}`)
    if (this.scene === Scene.Gaming) {
      this.$buefy.toast.open({
        duration: 3000,
        message: 'はじまるよ',
        type: 'is-success'
      })
      this.loading = false
      this.$router.push('/game')
    }
  }

  mounted() {
    this.$store.dispatch('game/match')
  }
}
</script>
