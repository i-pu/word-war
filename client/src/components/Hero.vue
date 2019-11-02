<template>
  <div>
    <div class="hero">
      <h2>uid: {{ $store.state.user.uid }}</h2>
    </div>
  </div>
</template>

<script lang="ts">
import "bulma/css/bulma.css"
import { Component, Vue } from "vue-property-decorator"

@Component
export default class Hero extends Vue {}
</script>
