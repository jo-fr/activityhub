<template>
  <div>
    <AddFeed />
    <h1>Explore existing Feeds:</h1>
    <FeedItem
      class="clickable"
      v-for="item in data"
      :key="item.name"
      :name="item.name"
      :description="item.description"
      @click="goToDetail(item.account.username)"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue'
import type { Ref } from 'vue'
import FeedItem from './FeedItem.vue'
import FeedDetail from './FeedDetail.vue'
import AddFeed from './AddFeed.vue'
import type { Feed } from '../models/models'
import { fetchFeeds } from '../api//api'

export default defineComponent({
  name: 'Explore',
  setup() {
    // Define a reactive property to store the fetched data
    const data: Ref<Feed[]> = ref([])
    // Fetch data from the API when the component is mounted
    onMounted(async () => {
      const json = await fetchFeeds()
      data.value = json.items
    })
    return { data }
  },

  methods: {
    goToDetail(username: string) {
      // Navigate to the DetailView with a prop (e.g., id)
      this.$router.push({ name: 'feedDetail', params: { username: username } })
    }
  },
  components: { FeedItem, FeedDetail, AddFeed }
})
</script>

<style scoped>
.clickable {
  cursor: pointer;
}
</style>
