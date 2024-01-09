<template>
    <div>
    <input class="addFeed" type="text"  v-model="feedURL" v-on:keyup.enter="submitData" placeholder="add new Feed URL"/>
    <p class="error" v-if="errorMessage">{{errorMessage}}</p>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent, ref, onMounted, defineProps  } from 'vue';
  import { stringifyQuery, useRoute } from 'vue-router'; 
  import type { Ref } from 'vue';
  
  interface Feed {
    id: string;
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
          const response = await fetch("/api/feeds", {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              feedURL: this.feedURL
            })
          });
          if (!response.ok) {
            if (response.status == 400) {
                const jsonData = await response.json();
     
                    this.errorMessage = jsonData.errors[0]?.message
            }
              throw new Error('error');
            }
    
    
          // Parse the JSON response
          const jsonData = await response.json();
          const feed: Feed = jsonData;

          this.$router.push({ name: 'feedDetail', params: { id: feed.id } });
        
      }
    }
  });
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
  