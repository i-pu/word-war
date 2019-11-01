<template>
  <div>
    <p>Siritori</p>
    <b-field>
      <b-input></b-input>
    </b-field>
    <div class="buttons">
      <b-button @click="send">Send</b-button>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";

@Component
export default class Siritori extends Vue {
  private wordWarClient?;
  private message: string;
  constructor() {
    super();
    this.message = "";
  }
  private send() {
    console.log(this.message);
  }
}
</script>
