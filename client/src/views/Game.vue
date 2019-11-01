<template>
  <div class="game">
    <p>WIP</p>
    <Siritori></Siritori>
  </div>
</template>

<script lang="ts">
import Siritori from "@/components/Siritori.vue";
import { Component, Vue } from "vue-property-decorator";

@Component({
  components: {
    Siritori
  }
})
export default class Game extends Vue {}
</script>
