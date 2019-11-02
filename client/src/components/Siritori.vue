<template>
  <div>
    <p>Siritori</p>
    <b-field>
      <b-input v-model="message"></b-input>
    </b-field>
    <div class="buttons">
      <b-button @click="send">Send</b-button>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator"
import { WordWarPromiseClient } from "@/pb/word_war_grpc_web_pb"

@Component
export default class Siritori extends Vue {
  private wordWarPromiseClient: WordWarPromiseClient;
  private message: string;
  constructor() {
    super()
    // TODO: 環境変数で切り替えるようにする
    this.wordWarPromiseClient = new WordWarPromiseClient(
      "http://localhost:50051",
      null,
      null
    )
    this.message = ""
  }
  private send() {
    console.log(this.message)
  }
}
</script>
