import { ref } from 'vue';
import type { Ref } from 'vue';
import { defineStore } from 'pinia';

import { getMaptilerToken } from '@/composables/useFetch';

export const useConfigStore = defineStore('config', () => {
  const maptilerToken: Ref<string> = ref('');

  async function $reset() {
    // TODO: error?
    const { data, error } = await getMaptilerToken();
    maptilerToken.value = data.value;
  };

  return {
    maptilerToken,
    $reset,
  };
});
