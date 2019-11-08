<template>
  <div>
    <div class="hero is-primary">
      <div class="hero-body">
        <div class="container">
          <h2 class="title has-text-left">Word War</h2>
        </div>
      </div>
      <div class="hero-footer">
        <div class="container">
          <h2 class="subtitle">uid: {{ $store.getters['user/uid'] }}</h2>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'

@Component
export default class Hero extends Vue {}
</script>
