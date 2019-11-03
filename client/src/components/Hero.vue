<template>
  <div>
    <div class="hero is-primary">
      <div class="container">
        <h2 class="subtitle">uid: {{ $store.state.user.uid }}</h2>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import 'bulma/css/bulma.css'
import { Component, Vue } from 'vue-property-decorator'

@Component
export default class Hero extends Vue {}
</script>
