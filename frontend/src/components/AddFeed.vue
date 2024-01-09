<template>
    <div>
    <input type="text"  v-model="feedURL" placeholder="add new Feed URL"/>
    <button @click="submitData">Submit</button>
    <p v-if="errorMessage">{{errorMessage}}</p>
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
  
          const response = await fetch("/api/feed", {
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
  