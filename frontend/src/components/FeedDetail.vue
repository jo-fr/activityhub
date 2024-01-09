<template>
    <div>
      <h1>{{data?.feed?.name}}</h1>
      <div clas="links">
        <div>
          {{ data?.feed?.account.preferredUsername }}
        </div>
        <div>
         <a target=”_blank” :href="data?.feed?.feedURL">Feed URL</a>
        </div>
        <div>
         <a target=”_blank” :href="data?.feed?.hostURL">Host Website</a>
        </div>
      </div>
      <p>{{data?.feed?.description}}</p>

        <div>
          <h2>Feed</h2>
            <div class="item" v-for="status in data?.status" :key="status.createdAt" v-html="status.Content"></div>
        </div>
    
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, ref, onMounted, defineProps  } from 'vue';
  import { stringifyQuery, useRoute } from 'vue-router'; 
  import type { Ref } from 'vue';
  
  interface Feed {
    name: string;
    type: string;
    feedURL: string;
    hostURL: string;
    author: string;
    description: string;
    imageURL: string;
    accountID: string;
    account: {
      preferredUsername: string;
  };
  }

  interface Status {
    createdAt: string;
    Content: string
  }

  
  export default defineComponent({
    name: 'FeedDetail',


    setup() {
      const data: Ref<{ feed: Feed | null; status: Status[] | null }> = ref({
      feed: null,
      status: null,
    });
        const route = useRoute();


        // Fetch data from the API when the component is mounted
        onMounted(async () => {
            try {
                // Make a fetch request using TypeScript
                const response = await fetch(`/api/feed/${route.params.id}`);
                if (!response.ok) {
                    throw new Error('Failed to fetch data');
                }
                // Parse the JSON response
                const jsonData = await response.json();
                // Update the data property with the fetched data
                data.value.feed = jsonData;
            }
            catch (error) {
                console.error('Error fetching data:', error);
            }
            try {
                // Make a fetch request using TypeScript
                const response = await fetch(`/api/feed/${route.params.id}/status?limit=10`);
                if (!response.ok) {
                    throw new Error('Failed to fetch data');
                }
                // Parse the JSON response
                const jsonData = await response.json();
                // Update the data property with the fetched data
                console.log(jsonData)
                data.value.status = jsonData.items;
            }
            catch (error) {
                console.error('Error fetching data:', error);
            }
        });
        return { data };
    },
  });
  </script>

  
  <style scoped>
  .item {
    border: 1px solid var( --color-border);
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
    display: inline;;
  }

  strong {
    color: forestgreen;
    font-weight: 800;
  }
  </style>
  