<template>
  <div>
    <h1>Explore existing Feeds:</h1>
   <FeedItem v-for="item in data" :key="item.name" :name="item.name" :description="item.description" />
  
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted  } from 'vue';
import type { Ref } from 'vue';
import FeedItem from './FeedItem.vue'


interface Feed {
  name: string;
  type: string;
  feedURL: string;
  hostURL: string;
  author: string;
  description: string;
  imageURL: string;
  accountID: string;
}

export default defineComponent({
    name: 'Explore',
    setup() {
        // Define a reactive property to store the fetched data
        const data: Ref<Feed[]> = ref([]);
        // Fetch data from the API when the component is mounted
        onMounted(async () => {
            try {
                // Make a fetch request using TypeScript
                const response = await fetch("/api/feed");
                if (!response.ok) {
                    throw new Error('Failed to fetch data');
                }
                // Parse the JSON response
                const jsonData = await response.json();
                // Update the data property with the fetched data
                data.value = jsonData.items;
            }
            catch (error) {
                console.error('Error fetching data:', error);
            }
        });
        return { data };
    },
    components: { FeedItem }
});
</script>

<style scoped>
/* Add your component styles here */
</style>
