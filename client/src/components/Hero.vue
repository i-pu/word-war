<template>
  <div>
    <template v-if="status.active && user">
      <div class="hero is-info is-small">
        <div class="hero-body">
          <div class="media">
            <div class="media-left mt-auto">
              <router-link :to="`/user/${user.userId}`" tag="figure" class="image is-48x48">
                <img class="is-rounded" :src="user.avatarUrl">
              </router-link>
            </div>
            <div class="media-content">
              <div class="mt-auto">
                  <p class="up has-text-white is-size-4 no-wrap">{{user.name}} id: {{ user.userId }}</p>
              </div>
              <div class="columns is-mobile">
                <div class="column is-1 mt-auto">
                  <p class="has-text-white-ter">Rating</p>
                  <p class="up has-text-white is-size-4">{{ user.rating }}</p>
                </div>
                <div class="column is-3 mt-auto">
                  <p class="has-text-white-ter">Points</p>
                  <p class="up has-text-white is-size-6">0</p>
                </div>
              </div>
            </div>
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
    const { status, user } = useRootStore()
    return {
      status,
      user,
    }
  }
})
</script>

<style scoped>
.mt-auto {
  margin-right: auto;
  margin-left: auto;
  margin-top: auto;
}

.up {
  padding-top: auto;
  margin-top: -0.5em;
}

.no-wrap {
  white-space: nowrap;
}

.is-vertical-center {
  display: flex;
  align-items: center;
}
</style>