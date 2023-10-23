import { ref } from 'vue';
import type { Ref } from 'vue';

import { useFetch_GetMaptilerToken } from '@/composables/useFetch';

const maptilerToken: Ref<string> = ref("");

export const useMaptilerToken = async () : Promise<string> => {
  if (maptilerToken.value !== "") {
    return maptilerToken.value;
  }

  // TODO: error?
  const { data, error } = await useFetch_GetMaptilerToken();
  maptilerToken.value = data.value;

  return maptilerToken.value;
};

