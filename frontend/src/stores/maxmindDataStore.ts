import { ref } from 'vue';
import type { Ref } from 'vue';
import { defineStore } from 'pinia';
import type { MaxmindBackendResponse } from '@/types/MaxmindBackend';
import { useFetch_GetMaxmindData } from '@/composables/useFetch';

export const useMaxmindDataStore = defineStore('maxindData', () => {
  const data: Ref<MaxmindBackendResponse | null> = ref(null);
  const error: Ref<string | null> = ref(null);

  async function $reset() {
    const { data: d, error: e } = await useFetch_GetMaxmindData('/api/geoip');
    data.value = d.value;
    error.value = e.value;
  };

  async function fetchForIP(ipAddress: string) {
    if (data.value!.IP === ipAddress) {
      return;
    }

    const { data: d, error: e } = await useFetch_GetMaxmindData(`/api/geoip/${ipAddress}`);
    data.value = d.value;
    error.value = e.value;
  };

  return {
    data,
    error,
    $reset,
    fetchForIP,
  };
});
