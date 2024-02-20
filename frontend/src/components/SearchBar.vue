<template>
  <div class="input-group mb-3">
    <span class="input-group-text">
      <SearchIcon />
    </span>
    <input
      type="text"
      class="form-control"
      id="searchbar"
      placeholder="Query details about IP Address..."
      aria-label="Query details about an IP Address"
      v-model="queryString"
      @keyup.enter="$emit('queryRequest', queryString)"
    />
  </div>
</template>

<script setup lang='ts'>
  import { ref, watch } from 'vue';
  import type { Ref } from 'vue';

  import SearchIcon from '@/components/icons/SearchIcon.vue';
  import { useMaxmindDataStore } from '@/stores/maxmindDataStore';
  import { storeToRefs } from 'pinia';

  defineEmits([
    'queryRequest',
  ]);

  const maxmindStore = useMaxmindDataStore();
  const { data } = storeToRefs(maxmindStore);

  const queryString: Ref<string> = ref('');
  watch(data, () => {
    if (queryString.value !== data.value.IP) {
      queryString.value = '';
    }
  });
</script>
