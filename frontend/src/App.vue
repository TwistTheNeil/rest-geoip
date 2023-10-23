<template>
  <div>
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useFetch_GetMaptilerToken } from './composables/useFetch';

useFetch_GetMaptilerToken();
</script>

