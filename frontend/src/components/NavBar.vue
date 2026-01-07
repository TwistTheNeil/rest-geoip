<template>
  <Menubar :model="items" style="margin-bottom: 1rem; border: none;">
    <template #end>
      <Button label="My IP Address" size="small" @click="$emit('queryReset')" />
    </template>
  </Menubar>
</template>

<script setup lang="ts">
import Menubar from 'primevue/menubar';
import Button from 'primevue/button';
import { ref } from 'vue';

const items = ref([
  {
    label: 'About',
    url: 'https://github.com/TwistTheNeil/rest-geoip'
  }
]);
</script>
