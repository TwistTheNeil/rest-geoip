<template>
  <div style="margin-bottom: 1rem;">
    <IconField>
      <InputIcon class="pi pi-search" />
      <InputText
        id="searchbar"
        v-model="queryString"
        placeholder="Query details about IP Address..."
        aria-label="Query details about an IP Address"
        @keyup.enter="$emit('queryRequest', queryString)"
        fluid
      />
    </IconField>
  </div>
</template>

<script setup lang='ts'>
  import { ref, watch } from 'vue';
  import type { Ref } from 'vue';
  import InputText from 'primevue/inputtext';
  import IconField from 'primevue/iconfield';
  import InputIcon from 'primevue/inputicon';

  import { useMaxmindDataStore } from '@/stores/maxmindDataStore';
  import { storeToRefs } from 'pinia';

  defineEmits([
    'queryRequest',
  ]);

  const maxmindStore = useMaxmindDataStore();
  const { data } = storeToRefs(maxmindStore);

  const queryString: Ref<string> = ref('');
  watch(data, () => {
    if (queryString.value !== data.value!.IP) {
      queryString.value = '';
    }
  });
</script>
