<template>
  <div>
    <template v-if="status.active">
      <div class="hero is-primary">
        <div class="hero-body">
          <div class="container">
            <h2 class="title has-text-left">
              Word War {{ status.version }}
            </h2>
          </div>
        </div>
        <div class="hero-footer">
          <div class="container">
            <h2 class="subtitle">userId:{{ userId }}</h2>
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
import { defineComponent } from '@vue/composition-api'
import { useRootStore } from '@/store'

export default defineComponent({
  setup() {
    const { status, userId } = useRootStore()
    return {
      status,
      userId
    }
  }
})
</script>
