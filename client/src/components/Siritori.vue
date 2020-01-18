<template>
  <div>
    <div class="container">
      <p>Siritori</p>
      <b-field>
        <b-input v-model="message" @keyup.native.enter="send"></b-input>
      </b-field>
      <p>room: {{ roomId }}</p>
      <p>{{ currentWord }}</p>
      <ul>
        <!-- FIXME: word.id を keyにするのダメそう -->
        <!-- 発言ごとに固有のランダムなIDを与えたい -->
        <li v-for="(word, i) in words" :key="i">
          {{ word.getUserid() }}: {{ word.getMessage() }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'

@Component
export default class Siritori extends Vue {
  private message: string = ''

  private get words() {
    return this.$store.getters['game/getWords']
  }

  private get roomId() {
    return this.$store.getters['game/roomId']
  }

  private get currentWord(): string {
    return this.words.length === 0
      ? 'り'
      : this.words[this.words.length - 1].getMessage()
  }

  async mounted() {
    try {
      console.log(this.$route.params)
      console.log(this.$route.query)
      if (this.$route.query.roomid) {
        const roomId = this.$route.query.roomid
        await this.$store.dispatch('game/start', { roomId })
        console.log(`joined ${roomId}`)
      } else {
        const roomId = await this.$store.dispatch('game/match')
        console.log(`match ${roomId}`)

        if (!roomId) {
          console.error('invalid roomid')
          return
        }

        await this.$store.dispatch('game/start', { roomId })
        console.log(`created and joined ${roomId}`)
      }
    } catch (e) {
      console.error(e)
    }

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
