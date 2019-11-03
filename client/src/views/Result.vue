<template>
  <div class="result">
    <Hero></Hero>
    <p>result: {{ score }}</p>
  </div>
</template>

<script lang="ts">
import { Component, Watch, Vue } from 'vue-property-decorator'
import Hero from '@/components/Hero.vue'
import { WordWarPromiseClient } from '@/pb/word_war_grpc_web_pb'
import { ResultRequest, ResultResponse } from '@/pb/word_war_pb'

@Component({
  components: {
    Hero
  }
})
export default class Result extends Vue {
  private wordWarPromiseClient: WordWarPromiseClient
  private result?: ResultResponse
  private score: string

  constructor() {
    super()
    // TODO: 環境変数で切り替えるようにする
    this.wordWarPromiseClient = new WordWarPromiseClient(
      'http://localhost:8080',
      null,
      null
    )
    this.score = ''
  }

  created() {
    const req = new ResultRequest()
    req.setUserid(this.$store.state.user.uid)
    this.wordWarPromiseClient
      .result(req)
      .then(result => {
        this.result = result
        console.log('result', this.result)
        this.score = result.getScore()
        console.log('score', this.score)
      })
      .catch(err => {
        console.log(err)
      })
  }
}
</script>
