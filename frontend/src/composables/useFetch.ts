import { ref } from 'vue';
import type { Ref } from 'vue'

import type { MaxmindBackendResponse, ConfigBackendResponse } from '@/types/MaxmindBackend';

const BACKEND_HOST = import.meta.env.VITE_BACKEND_HOST;

export async function getMaxmindData(ipAddress?: string): Promise<{ data: Ref<MaxmindBackendResponse | null>, error: Ref<string | null> }> {
  const data: Ref<MaxmindBackendResponse | null> = ref(null);
  const error: Ref<string | null> = ref(null);

  try {
    let url = `${BACKEND_HOST}/api/geoip`;
    if (ipAddress) {
      url = `${url}/${ipAddress}`;
    }
    const fetchPromise = await fetch(url);
    data.value = await fetchPromise.json();
  } catch (err: any) {
    error.value = err.message;
  }

  return { data, error };
};

export async function getConfig(): Promise<{ data: ConfigBackendResponse | null, error: string | null }> {
  try {
    const fetchPromise = await fetch(`${BACKEND_HOST}/api/config`);
    const data: ConfigBackendResponse = await fetchPromise.json();
    return { data, error: null };
  } catch (err: any) {
    return { data: null, error: err.message };
  }
};

