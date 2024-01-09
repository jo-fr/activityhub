<template>
    <div>
      <h1>{{data.name}}</h1>
      <p>{{data.description}}</p>
    
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, ref, onMounted, defineProps  } from 'vue';
  import { useRoute } from 'vue-router'; 
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
}

  
  export default defineComponent({
    name: 'FeedDetail',


    setup() {
        const data: Ref<Feed> = ref({name: "", type: "", feedURL: "", hostURL: "", author: "", description: "", imageURL: "", accountID: ""});
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
                data.value = jsonData;
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
  /* Add your component styles here */
  </style>
  