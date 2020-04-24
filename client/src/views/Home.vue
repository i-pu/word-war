<template>
  <div class="home">
    <Hero />
    <Navbar />
    <section class="section">
      <b-button type="is-link" tag="router-link" to="/waiting">
        ランダムマッチ
      </b-button>

      <!-- <b-button disabled type="is-link" tag="router-link" to="/game">
        へやをつくる
      </b-button>

      <b-button disabled @click="joinRoom">
        へやにさんかする
      </b-button> -->
      <!-- <b-field>
        <b-input v-model="roomIdInput" placeholder="input room ID"></b-input>
      </b-field>
      <b-button type="is-dark" @click="onClick">vuex</b-button>
      <p>{{ count }}</p> -->
    </section>

    <b-button class="is-danger">ログアウト</b-button>
  </div>
</template>

<script lang="ts">
import { useRootStore } from '@/store'
import Hero from '@/components/Hero.vue'
import Navbar from '@/components/Navbar.vue'
import { defineComponent, reactive, SetupContext } from '@vue/composition-api'

export default defineComponent({
  components: {
    Hero,
    Navbar
  },
})
</script>
