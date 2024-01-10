<template>
  <div>
    <input
      class="addFeed"
      type="text"
      v-model="feedURL"
      v-on:keyup.enter="submitData"
      placeholder="add new Feed URL"
    />
    <p class="error" v-if="errorMessage">{{ errorMessage }}</p>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { addFeed } from '../api/api'

export default defineComponent({
  name: 'AddFeed',
  data() {
    return {
      feedURL: '',
      errorMessage: ''
    }
  },

  methods: {
    async submitData() {
      try {
        const feed = await addFeed(this.feedURL)

        if (feed) {
          this.$router.push({ name: 'feedDetail', params: { username: feed.account.username } })
        }
      } catch (error) {
        console.log(error)
      }
    }
  }
})
</script>

<style scoped>
.addFeed {
  width: 100%;
  padding: 12px 20px;
  margin: 8px 0;
  background-color: var(--main-color);
  color: white;
  border: 2px solid var(--color-border);
  border-radius: 4px;
  box-sizing: border-box;
  font-size: 1.5rem;
}

.error {
  color: red;
}
</style>
