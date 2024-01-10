<template>
  <div>
    <h1>{{ data?.feed?.name }}</h1>
    <div clas="links">
      <div>
        {{ data?.feed?.account.uri }}
      </div>
      <div>
        <a target="”_blank”" :href="data?.feed?.feedURL">Feed URL</a>
      </div>
      <div>
        <a target="”_blank”" :href="data?.feed?.hostURL">Host Website</a>
      </div>
    </div>
    <p>{{ data?.feed?.description }}</p>

    <div>
      <h2>Feed</h2>
      <div
        class="item"
        v-for="status in data?.status"
        :key="status.createdAt"
        v-html="status.content"
      ></div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, defineProps } from 'vue'
import { useRoute } from 'vue-router'
import type { Ref } from 'vue'
import type { Feed, Status } from '../models/models'
import { fetchUserFeed, fetchFeedStatus } from '../api/api'

export default defineComponent({
  name: 'FeedDetail',

  setup() {
    const data: Ref<{ feed: Feed | null; status: Status[] | null }> = ref({
      feed: null,
      status: null
    })
    const route = useRoute()

    // Fetch data from the API when the component is mounted
    onMounted(async () => {
      try {
        const input = route.params.username
        const username = Array.isArray(input) ? input[0] : input
        data.value.feed = await fetchUserFeed(username)
      } catch (error) {
        console.error('Error fetching user feed data:', error)
      }

      try {
        data.value.status = await fetchFeedStatus(data.value.feed?.id!)
      } catch (error) {
        console.error('Error fetching feed status data:', error)
      }
    })
    return { data }
  }
})
</script>

<style scoped>
.item {
  border: 1px solid var(--color-border);
  border-radius: 4px;
  padding: 10px;
  margin: 10px 0;
}

.links {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.links > div {
  display: inline;
}

strong {
  color: forestgreen;
  font-weight: 800;
}
</style>
