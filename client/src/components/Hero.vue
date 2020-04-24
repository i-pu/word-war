<template>
  <div>
    <template v-if="$store.getters.serverHealth">
      <div class="hero is-primary">
        <div class="hero-body">
          <div class="container">
            <h2 class="title has-text-left">
              Word War {{ $store.getters.version }}
            </h2>
          </div>
        </div>
        <div class="hero-footer">
          <div class="container">
            <h2 class="subtitle">userId:{{ $store.getters.userId }}</h2>
          </div>
        </div>
      </div>
    </template>
    <template v-else>
      <div class="hero is-danger">
        <div class="hero-body">
          <div class="container">
            <h2 class="title has-text-left">
              Word War Not Available
            </h2>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { defineComponent, } from '@vue/composition-api'

export default defineComponent({
})
</script>
