<template>
  <div>
    <div class="container">
      <p>Siritori</p>
      <b-field>
        <b-input v-model="message"></b-input>
      </b-field>
      <div class="buttons">
        <b-button @click="send">Send</b-button>
      </div>
      <ul>
        <!-- FIXME: word.id を keyにするのダメそう -->
        <li v-for="(word, i) in siritoriWords" :key="i">
          {{ word.uid }}: {{ word.uid }}: {{ word.message }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { WordWarPromiseClient } from '@/pb/word_war_grpc_web_pb'
import { SayRequest, GameRequest } from '@/pb/word_war_pb'

@Component
export default class Siritori extends Vue {
  // TODO: 環境変数で切り替えるようにする
  private wordWarPromiseClient: WordWarPromiseClient = new WordWarPromiseClient(
    'http://localhost:8080'
  )
  private message: string = ''
  private siritoriWords: Array<{ uid: string; message: string }> = []

  created() {
    const req: GameRequest = new GameRequest()
    req.setUserid(this.$store.state.user.uid)
    const stream = this.wordWarPromiseClient.game(req)

    stream.on('data', res => {
      console.log(res)
      this.siritoriWords.push({
        uid: res.getUserid(),
        message: res.getMessage()
      })
    })

    stream.on('status', status => {
      console.log('status', status)
      if (status.code === 0) {
        this.$router.push('/result')
      }
    })

    stream.on('error', res => {
      console.log('error', res)
    })
  }

  private send() {
    const req: SayRequest = new SayRequest()
    req.setUserid(this.$store.state.user.uid)
    req.setMessage(this.message)
    console.log(req)
    this.wordWarPromiseClient
      .say(req)
      .then(result => {
        console.log(result)
      })
      .catch(err => {
        console.log(err)
      })
  }
}
</script>
