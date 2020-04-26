<template>
  <div class="home">
    <Hero />
    <!-- <Navbar /> -->

    <div class="tile is-ancester">
      <div class="tile is-parent is-vertical">
        <div class="tile is-child notification is-danger is-paddingless rounded">
          <router-link to="/waiting" tag="b-button" class="v-30 is-primary is-fullwidth rounded">
              <p>ランダムマッチ</p>
          </router-link>
        </div>
        <div class="tile is-child notification is-paddingless rounded">
          <router-link to="/waiting" tag="b-button" class="v-20 is-info is-fullwidth rounded" disabled>
              <p>プライベートマッチ</p>
          </router-link>
        </div>
      </div>
    </div>

    <section class="section">

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
  </div>
</template>

<script lang="ts">
import { useRootStore } from '@/store'
import Hero from '@/components/Hero.vue'
// import Navbar from '@/components/Navbar.vue'
import { defineComponent, reactive, SetupContext } from '@vue/composition-api'

export default defineComponent({
  components: { Hero },
})
</script>

<style scoped>
.rounded {
  border-radius: 2em;
}

.v-30 {
  height: 30vh;
}

.v-20 {
  height: 20vh;
}
</style>