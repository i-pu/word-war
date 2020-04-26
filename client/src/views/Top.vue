<template>
  <div>
    <img @click="toLogin" class="img-responsive" src="https://imgc.appbank.net/c/wp-content/uploads/2018/12/24-1.jpg"/>
  </div>
</template>

<script lang="ts">
import { defineComponent, SetupContext, onMounted, reactive } from '@vue/composition-api'
export default defineComponent({
  setup(props:{}, {root}: SetupContext) {
    const toLogin = () => {
      root.$router.push('/login')
    }

    return {
      toLogin
    }
  }
})
</script>

<style scoped>
.img-responsive {
  display: block;
  width: auto;
  min-height: 100vh;
}
</style>