<template>
  <div class="form-floating mt-4 mb-4">
    <input
      type="text"
      class="form-control border-0 shadow-none"
      id="searchbar"
      placeholder="Query details about IP Address..."
      v-model="queryString"
      @keyup.enter="$emit('queryRequest', queryString)"
    />
    <label for="searchBar">Query details about IP Address...</label>
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
