<template>
  <div class="top">
    <Hero />
    <Login />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Login from '@/components/Login.vue'
import Hero from '@/components/Hero.vue'

@Component({
  components: { Hero, Login }
})
export default class Top extends Vue {
  mounted() {
    this.$store.dispatch('healthCheck')
  }
}
</script>
