<template>
  <Message
    v-if="adminNotice"
    severity="info"
    :closable="false"
    style="margin-bottom: 1rem; margin-top: 1.5rem;"
  >
    {{ adminNotice }}
  </Message>
</template>

<script setup lang="ts">
  import { storeToRefs } from 'pinia';
  import Message from 'primevue/message';
  import { useConfigStore } from '@/stores/config';

  const configStore = useConfigStore();
  const { adminNotice } = storeToRefs(configStore);
</script>
