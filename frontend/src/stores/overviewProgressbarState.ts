import { ref } from 'vue';
import type { Ref } from 'vue';
import { defineStore } from 'pinia';

export const useOverviewProgressbarStore = defineStore('overviewProgressbar', () => {
  const progressValue: Ref<number> = ref(0); 

  function update(value: number) {
    progressValue.value = value;
  };

  function $reset() {
    progressValue.value = 0;
  };

  return {
    progressValue,
    update,
    $reset,
  };
});
