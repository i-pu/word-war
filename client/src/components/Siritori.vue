<template>
  <div>
    <div class="container">
      <p>Siritori</p>
      <b-field>
        <b-input v-model="message" @keyup.native.enter="send"></b-input>
      </b-field>
      <ul>
        <!-- FIXME: word.id を keyにするのダメそう -->
        <li v-for="(word, i) in words" :key="i">
          {{ word.getUserid() }}: {{ word.getMessage() }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
// TODO: room_idを指定してゲームが始められるように
import { Component, Vue } from 'vue-property-decorator'

@Component
export default class Siritori extends Vue {
  private message: string = ''

  private get words() {
    return this.$store.getters['game/getWords']
  }

  mounted() {
    this.$store.dispatch('game/matchAndStart')

    this.$store.watch(
      (state, getter) => {
        return state.game.isPlaying
      },
      (isPlaying, old) => {
        console.log(`${old} => ${isPlaying}`)
        if (!isPlaying) {
          this.$router.push('/result')
        }
      }
    )
  }

  private async send() {
    this.$store.dispatch('game/say', { message: this.message })
    this.message = ''
  }
}
</script>
