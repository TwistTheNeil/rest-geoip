<template>
  <div>
    <OverviewProgressbar />
    <NavBar />

    <div class="container justify-content-center align-items-center">
      <SearchBar @query-request="fetchIPAddressDetails" />

      <slot />
    </div>
  </div>
</template>

<script setup lang='ts'>
  import { isIP } from 'is-ip';

  import NavBar from '@/components/NavBar.vue';
  import OverviewProgressbar from '@/components/OverviewProgressbar.vue';
  import SearchBar from '@/components/SearchBar.vue';
  import { useOverviewProgressbarStore } from '@/stores/overviewProgressbarState';
  import { useMaxmindDataStore } from '@/stores/maxmindDataStore';

  const overviewProgressbarStore = useOverviewProgressbarStore();
  const maxmindDataStore = useMaxmindDataStore();

  const fetchIPAddressDetails = async (ipAddress: string) => {
    if (!isIP(ipAddress)) {
      return;
    }

    overviewProgressbarStore.update(30);

    await maxmindDataStore.fetchForIP(ipAddress);

    overviewProgressbarStore.update(100);
  };
</script>
