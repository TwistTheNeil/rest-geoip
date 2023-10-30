import { ref } from 'vue';
import type { Ref } from 'vue';

import { getMaptilerToken } from '@/composables/useFetch';

const maptilerToken: Ref<string> = ref("");

export const useMaptilerToken = async () : Promise<string> => {
  if (maptilerToken.value !== "") {
    return maptilerToken.value;
  }

  // TODO: error?
  const { data, error } = await getMaptilerToken();
  maptilerToken.value = data.value;

  return maptilerToken.value;
};

