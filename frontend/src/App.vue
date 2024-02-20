<template>
  <!-- https://vuejs.org/guide/essentials/component-basics.html#dynamic-components -->
  <component v-if="ready" :is="route.meta.layout || 'div'">
    <RouterView />
  </component>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from 'vue';
import { RouterView } from 'vue-router';
import { useRoute } from 'vue-router';

import { useConfigStore } from '@/stores/config';
import { useMaxmindDataStore } from '@/stores/maxmindDataStore';

const route = useRoute();
const maxmindDataStore = useMaxmindDataStore();
const configStore = useConfigStore();
const ready = ref(false);

onBeforeMount(async () => {
  await configStore.$reset();

  if (!route.query.address || route.query.address === "") {
    await maxmindDataStore.$reset();
  } else {
    await maxmindDataStore.fetchForIP(route.query.address as string);
  }
  ready.value = true;
});
</script>

