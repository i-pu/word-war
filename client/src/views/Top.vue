<template>
  <div class="top">
    <p>WIP</p>
    <p>firebaseの認証をする</p>
    <Login />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator"
import Login from "@/components/Login.vue"

@Component({
  components: { Login }
})
export default class Top extends Vue {}
</script>
