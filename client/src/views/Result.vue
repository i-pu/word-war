<template>
  <div class="result">
    <Hero />
    <Navbar />
    <section class="section">
      <p>result: {{ score }}</p>
    </section>

    <b-button @click="to('/home')" >ホームに戻る</b-button>
  </div>
</template>

<script lang="ts">
import Hero from '@/components/Hero.vue'
import Navbar from '@/components/Navbar.vue'
import { useGameStore } from '@/store/game'
import { defineComponent, reactive, SetupContext, onMounted } from '@vue/composition-api'

export default defineComponent({
  components: {
    Hero,
    Navbar
  },
  setup(props: {}, { root }: SetupContext) {
    const { result, score } = useGameStore()
    onMounted(async () => {
      await result()
      console.log(`score: ${score}`)
    })

    const to = (path: string) => {
      root.$router.push(path)
    }

    return {
      score, to
    }
  }
})
</script>
