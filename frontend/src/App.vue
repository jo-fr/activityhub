<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header>

    <h1><RouterLink to="/">ActivityHub</RouterLink></h1> 

  </header>
  <RouterView />
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}


h1 {
  font-size: 2rem;
  font-weight: 700;
  margin: 0;
  text-align: center;
}


</style>
