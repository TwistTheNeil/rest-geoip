<template>
  <ProgressBar :value="progressValue" class="fast-progress" :showValue="false" />
</template>

<script setup lang='ts'>
  import { computed } from 'vue';
  import ProgressBar from 'primevue/progressbar';
  import { useOverviewProgressbarStore } from '@/stores/overviewProgressbarState';

  const overviewProgressbarStore = useOverviewProgressbarStore();

  const progressValue = computed(() => overviewProgressbarStore.progressValue);

  overviewProgressbarStore.$subscribe((mutation, state) => {
    if (state.progressValue === 100) {
      setTimeout(() => overviewProgressbarStore.$reset(), 1000);
    }
  });
</script>

<style scoped>
.fast-progress {
  height: 3px;
  border-radius: 0;
}

.fast-progress :deep(.p-progressbar-value) {
  transition: width 0.1s ease-in-out;
}
</style>
