<template>
  <div>
    <div class="container">
      <p>Siritori</p>
      <b-field>
        <b-input v-model="state.message" @keyup.native.enter="send"></b-input>
      </b-field>
      <p>room: {{ roomId }}</p>
      <p>{{ currentWord }}</p>
      <ul>
        <!-- FIXME: word.id を keyにするのダメそう -->
        <!-- 発言ごとに固有のランダムなIDを与えたい -->
        <li v-for="(word, i) in words" :key="i">
          {{ word.id }}: {{ word.text }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { Word, Scene } from '@/model'
import { defineComponent, reactive, computed, watch, onMounted } from '@vue/composition-api'
import { useGameStore } from '@/store/game'

export default defineComponent({
  setup(props: {}, { root }) {
    const state = reactive({
      message: ''
    })

    const { words, roomId, scene, start, say } = useGameStore()

    const currentWord = computed(() => {
      return words.value.length === 0
      ? 'り'
      : words.value[words.value.length - 1].text
    })

    watch([scene], () => {
      console.log("scene watched")
      console.log(`シーンが ${Scene[scene.value]} になったよ`)
      if (scene.value === Scene.End) {
        root.$router.push('/result')
      }
    })

    onMounted(async () => {
      await start()
      console.log('Game started')
    })

    const send = async() => {
      say(state.message)
      state.message = ''
    }

    return {
      state,
      currentWord,
      words,
      roomId,
      scene,
      send,
    }
  }
})
</script>
