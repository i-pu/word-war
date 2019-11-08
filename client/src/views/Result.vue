<template>
  <div class="result">
    <Hero />
    <Navbar />
    <section class="section">
      <p>result: {{ score }}</p>
    </section>
  </div>
</template>

<script lang="ts">
import { Component, Watch, Vue } from 'vue-property-decorator'
import Hero from '@/components/Hero.vue'
import Navbar from '@/components/Navbar.vue'
import { WordWarPromiseClient } from '@/pb/word_war_grpc_web_pb'
import { ResultRequest, ResultResponse } from '@/pb/word_war_pb'

@Component({
  components: {
    Hero,
    Navbar
  }
})
export default class Result extends Vue {
  private wordWarPromiseClient: WordWarPromiseClient = new WordWarPromiseClient(
    process.env.VUE_APP_API_ENDPOINT
  )
  private result?: ResultResponse
  private score: string = ''

  async created() {
    const req = new ResultRequest()
    req.setUserid(this.$store.getters['user.uid'])
    const result = await this.wordWarPromiseClient
      .result(req)
      .catch(console.error)
    if (result !== undefined) {
      this.result = result
      this.score = this.result.getScore()
    }
    console.log('score', this.score)
  }
}
</script>
