<template>
  <!-- https://vuejs.org/guide/essentials/component-basics.html#dynamic-components -->
  <component :is="route.meta.layout || 'div'">
    <RouterView />
  </component>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router';
import { useFetch_GetMaptilerToken } from './composables/useFetch';
import { useRoute } from 'vue-router';

const route = useRoute();
console.log('do i get here');
console.log(route.meta.layout);

useFetch_GetMaptilerToken();
</script>

