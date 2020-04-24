<template>
  <div class="waiting">
    <Hero />

    <b-loading is-full-page :active="state.loading" />

    <h1>待機中 scene: {{ scene }}. {{ players }}..</h1>

    <p v-for="player in players" :key="player.playerId">
      {{ player }}
    </p>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, SetupContext, onMounted, watch } from '@vue/composition-api'
import { useGameStore } from '@/store/game'
import Hero from '@/components/Hero.vue'
import Navbar from '@/components/Navbar.vue'
import { Scene, User } from '@/model'

export default defineComponent({
  components: {
    Hero,
    Navbar
  },
  setup(props: {}, {root}: SetupContext) {
    const { scene, players, match } = useGameStore()
    const state = reactive({
      loading: true
    })

    watch([scene], () => {
      console.log(`sceneChanged ${Scene[scene.value]}`)
      if (scene.value === Scene.Gaming) {
        root.$buefy.toast.open({
          duration: 3000,
          message: 'はじまるよ',
          type: 'is-success'
        })
        state.loading = false
        root.$router.push('/game')
      }
    })

    onMounted(async () => {
      await match()
    })

    return { scene, players, state }
  }
})
</script>
