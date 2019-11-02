<template>
  <div class="result">
    <Hero></Hero>
    <p>WIP</p>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator"
import Hero from "@/components/Hero.vue"

@Component({
  components: {
    Hero
  }
})
export default class Result extends Vue {}
</script>
