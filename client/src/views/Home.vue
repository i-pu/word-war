<template>
  <div class="home">
    <Hero />
    <Navbar />
    <section class="section">
      <b-button type="is-link" tag="router-link" to="/game">
        へやをつくる
      </b-button>

      <b-button @click="joinRoom">
        へやにさんかする
      </b-button>
      <b-field>
        <b-input placeholder="input room ID" v-model="roomIdInput"></b-input>
      </b-field>
      <b-button type="is-dark" @click="onClick">vuex</b-button>
      <p>{{ count }}</p>
    </section>
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

  get count(): number {
    return this.$store.state.sample.count
  }

  private onClick() {
    this.$store.commit('sample/increment')
    console.log(this.$store.state.sample.count)
  }

  private joinRoom() {
    console.log(`roomId: ${this.roomIdInput}`)
    this.$router.push(`/game?roomid=${this.roomIdInput}`)
  }
}
</script>
