import { ref } from 'vue';
import type { Ref } from 'vue'

import type { MaxmindBackendResponse } from '@/types/MaxmindBackend';

// TODO: why param? remove
export async function useFetch_GetMaxmindData(url: string): Promise<{ data: Ref<MaxmindBackendResponse| null>, error: Ref<string | null> }> {
  const data: Ref<MaxmindBackendResponse | null> = ref(null);
  const error: Ref<string | null> = ref(null);

  try {
    const fetchPromise = await fetch(url);
    data.value = await fetchPromise.json();
  } catch (err: any) {
    error.value = err.message;
  }

  return { data, error };
};

export async function useFetch_GetMaptilerToken(): Promise<{ data: Ref<string>, error: Ref<string | null> }> {
  const data: Ref<string> = ref("");
  const error: Ref<string | null> = ref(null);

  try {
    const fetchPromise = await fetch("/api/maptiler/token");
    data.value = (await fetchPromise.json()).MaptilerToken;
  } catch (err: any) {
    error.value = err.message;
  }

  return { data, error };
};

