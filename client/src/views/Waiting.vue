<template>
  <div class="waiting">
    <b-loading is-full-page :active="state.loading" />
    <div class="content">
      <h1>RoomId <strong>{{ roomId }}</strong></h1>
      <h3>部屋の設定
        <ul>
          <li>ゴニョゴニョ</li>
          <li>ゴニョゴニョ</li>
          <li>ゴニョゴニョ</li>
        </ul>
      </h3>
      <h1>{{ Scene[scene] }}</h1>
    </div>

    <div class="tile is-ancester">
      <div class="tile is-parent is-vertical">
        <div v-for="player in players" :key="player.playerId" class="tile is-child notification is-danger rounded">
          <article class="media">
            <figure class="media-left">
              <p class="image is-64x64">
                <img class="is-rounded" src="https://versions.bulma.io/0.7.0/images/placeholders/128x128.png">
              </p>
            </figure>
            <div class="media-content">
              <div class="content">
                <p><strong>John Smith</strong></p>
                <p> {{ player.playerId }} </p>
              </div>
            </div>
          </article>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive, SetupContext, onMounted, watch, onUnmounted } from '@vue/composition-api'
import { useGameStore } from '@/store/game'
import { Scene, User } from '@/model'

export default defineComponent({
  setup(props: {}, {root}: SetupContext) {
    const { scene, players, match, roomId, cancelMatching } = useGameStore()
    const state = reactive({
      loading: true
    })

    watch([scene], () => {
      console.log(`sceneChanged ${Scene[scene.value]}`)
      if (scene.value === Scene.Gaming) {
        root.$buefy.toast.open({
          duration: 3000,
          message: 'はじまるよ',
          type: 'is-success'
        })
        state.loading = false
        root.$router.push('/game')
      }
      // else if (scene.value === Scene.None) {
      //   alert('ホームに戻ります')
      //   state.loading = false
      //   root.$router.push('/home')
      // }
    })

    onMounted(async () => {
      await match()
    })

    onUnmounted(async () => {
      console.log('マッチングが中断されます')
      await cancelMatching()
    })

    return { roomId, scene, players, state, Scene }
  }
})
</script>
<style scoped>
.rounded {
  border-top-left-radius: 2em;
  border-bottom-right-radius: 2em;
}
</style>