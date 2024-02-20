import { ref } from 'vue';
import type { Ref } from 'vue';
import { defineStore } from 'pinia';

import { getConfig } from '@/composables/useFetch';

export const useConfigStore = defineStore('config', () => {
  const maptilerToken: Ref<string> = ref('');
  const adminNotice: Ref<string | null> = ref(null);

  async function $reset() {
    // TODO: error?
    const { data, error } = await getConfig();
    maptilerToken.value = data!.MaptilerToken;
    adminNotice.value = data!.AdminNotice;
  };

  return {
    maptilerToken,
    adminNotice,
    $reset,
  };
});
