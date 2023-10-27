<template>
  <div>
    <OverviewProgressbar />
    <NavBar />

    <div class="container justify-content-center align-items-center">
      <slot />
    </div>
  </div>
</template>

<script setup lang='ts'>
  import NavBar from '@/components/NavBar.vue';
  import OverviewProgressbar from '@/components/OverviewProgressbar.vue';
</script>
