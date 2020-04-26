<template>
  <div v-if="user">
    <!--
    <button class="button is-danger back-rounded">
      <span class="icon">
        <i class="mdi mdi-48px mdi-menu-left"></i>
      </span>
      <span>戻る</span>
    </button>
    -->
    <section class="section is-mediam">
      <article class="media">
        <figure class="media-left">
          <p class="image is-64x64">
            <img class="is-rounded" :src="user.avatarUrl">
          </p>
        </figure>
        <div class="media-content">
          <div class="content">
            <p class="is-marginless" :class="ratingToColor(16.0)">
              <i class="mdi mdi-30px mdi-crown"></i>
              {{ user.name }}
            </p>
            <p> ID: {{ user.userId }} </p>
          </div>
        </div>
      </article>
    </section>

    <div class="tile is-ancester">
      <div class="tile is-parent is-vertical">
        <div class="tile is-child notification is-grey-light rounded">
          <div class="content">
            <h1>戦績</h1>
          </div>
          <user-rating-graph :ratings="user.history" />
        </div>
      </div>
    </div>

    <b-button @click="logOut" rounded class="is-danger">ログアウト</b-button>
  </div>
</template>

<script lang="ts">
import { defineComponent, SetupContext, onMounted, reactive, ref, Ref } from '@vue/composition-api'
import { useRootStore } from '@/store'
import UserRatingGraph from '@/components/UserRatingGraph.vue'

const ratingToColor = (rating: number): object => {
  return {
    'is-size-3': true,
    'has-text-grey': rating <= 5.0,
    'has-text-success': rating <= 10.0,
    'has-text-warning': rating <= 15.0,
    'has-text-danger': 15.0 < rating,
  }
}

export default defineComponent({
  components: { UserRatingGraph },
  setup(props:{}, {root}: SetupContext) {
    const userId = ref(root.$route.params['id'])
    const { user, signOut } = useRootStore()

    const logOut = async () => {
      await signOut()
      console.log('ログアウトしました')
      root.$router.push('/')
    }

    return {
      userId,
      user,
      ratingToColor,
      logOut,
    }
  }
})
</script>

<style scoped>
.rounded {
  border-radius: 2em;
}
.back-rounded {
  border-top-left-radius: 0px;
  border-top-right-radius: 0px;
  border-bottom-left-radius: 0px;
  border-bottom-right-radius: 20px;
}
</style>