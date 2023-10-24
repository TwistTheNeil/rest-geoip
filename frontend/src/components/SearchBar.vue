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
  import { ref } from 'vue';
  import type { Ref } from 'vue';

  import SearchIcon from '@/components/icons/SearchIcon.vue';

  const props = defineProps({
    currentIPAddress: {
      type: String,
      required: true,
    },
  });

  const emit = defineEmits([
    'queryRequest',
  ]);

  const queryString: Ref<string> = ref(props.currentIPAddress);
</script>
