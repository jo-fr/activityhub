<template>
  <div class="item">
    <h2>{{ $props.name }}</h2>
    <p>{{ $props.description }}</p>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted  } from 'vue';
export default defineComponent({
  name: 'FeedItem',
  props: {
    name: String,
   description: String,
  },
 
})
</script>


<style scoped>
.item {
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 15px;
  background-color:var(--color-background);

}

h2 {
  margin-bottom: 10px;
}

p {
  margin: 0;
}


</style>
