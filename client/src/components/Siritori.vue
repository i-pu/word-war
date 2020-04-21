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
          {{ word.id }}: {{ word.text }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import { Word, Scene } from '@/store/game'

@Component({})
export default class Siritori extends Vue {
  private message: string = ''

  private get words(): Word[] {
    return this.$store.getters['game/getWords'] as Word[]
  }

  private get roomId(): string {
    return this.$store.getters['game/roomId']
  }

  private get scene(): Scene {
    return this.$store.getters['game/scene']
  }

  private get currentWord(): string {
    return this.words.length === 0
      ? 'り'
      : this.words[this.words.length - 1].text
  }

  @Watch('scene')
  private onSceneChanged(scene: Scene) {
    console.log(`シーンが ${Scene[scene]} になったよ`)
    if (scene === Scene.End) {
      this.$router.push('/result')
    }
  }

  async mounted() {
    // start
    await this.$store.dispatch('game/start')
    console.log('Game started')
  }

  private async send() {
    this.$store.dispatch('game/say', { message: this.message })
    this.message = ''
  }
}
</script>
