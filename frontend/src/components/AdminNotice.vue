<template>
  <div
    v-if="adminNotice"
    class="alert alert-info d-flex align-items-center"
    role="alert"
  >
    {{ adminNotice }}
  </div>
</template>

<script setup lang="ts">
  import { storeToRefs } from 'pinia';
  import { useConfigStore } from '@/stores/config';

  const configStore = useConfigStore();
  const { adminNotice } = storeToRefs(configStore);
</script>
