<template>
  <div class="progress" style="height: 3px;">
    <div class="progress-bar" role="progressbar" :style="progressValueWidth" :aria-valuenow="progressValue" aria-valuemin="0" aria-valuemax="100"></div>
  </div>
</template>

<script setup lang='ts'>
  import { computed } from 'vue';
  import { useOverviewProgressbarStore } from '@/stores/overviewProgressbarState';

  const overviewProgressbarStore = useOverviewProgressbarStore();

  const progressValue = computed(() => overviewProgressbarStore.progressValue);

  overviewProgressbarStore.$subscribe((mutation, state) => {
    if (state.progressValue === 100) {
      setTimeout(() => overviewProgressbarStore.$reset(), 1000);
    }
  });

  const progressValueWidth = computed(() => `width: ${overviewProgressbarStore.progressValue}%;`);
</script>
