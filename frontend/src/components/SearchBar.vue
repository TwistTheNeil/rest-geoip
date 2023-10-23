<template>
  <div class="input-group mb-3">
    <span class="input-group-text">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-search" viewBox="0 0 16 16">
        <path d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001c.03.04.062.078.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1.007 1.007 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0z"/>
      </svg>
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
