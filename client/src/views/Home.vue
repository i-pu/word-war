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

    <b-button class="is-danger" @click="signOut">ログアウト</b-button>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Hero from '@/components/Hero.vue'
import Navbar from '@/components/Navbar.vue'

@Component({
  components: {
    Hero,
    Navbar
  }
})
export default class Home extends Vue {
  private roomIdInput: string = ''

  // private joinRoom() {
  //   console.log(`roomId: ${this.roomIdInput}`)
  //   this.$router.push(`/waiting?roomid=${this.roomIdInput}`)
  // }

  private signOut() {
    this.$store.dispatch('signOut')
  }
}
</script>
