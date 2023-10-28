<template>
  <!-- TODO: figure out why height: 100% doesn't work on mobile -->
  <div id='map' style='height: 500px'></div>
</template>

<script setup lang='ts'>
  import { onMounted, onUpdated, computed } from 'vue';
  import maplibregl from 'maplibre-gl';
  import { storeToRefs } from 'pinia';

  import { useMaptilerToken } from '@/composables/useMaptilerToken';
  import { useMaxmindDataStore } from '@/stores/maxmindDataStore';

  const maxmindDataStore = useMaxmindDataStore();
  const { data } = storeToRefs(maxmindDataStore);

  const center = computed(() => ({ lng: data.value!.Location.Longitude, lat: data.value!.Location.Latitude }));

  const updateMap = async () => {
    // TODO: move this out? i don't like it being called on every update
    // currently error if we move it out of `onMounted` specifically
    const maptilerToken = await useMaptilerToken();
    const map = new maplibregl.Map({
      container: 'map',
      style: `https://api.maptiler.com/maps/streets/style.json?key=${maptilerToken}`,
      center: center.value,
      zoom: 10,
    });

    new maplibregl.Marker()
      .setLngLat(center.value)
      .addTo(map);
  };

  onMounted(() => updateMap());
  onUpdated(() => updateMap());
</script>

<style>
  @import '@/../node_modules/maplibre-gl/dist/maplibre-gl.css';
</style>
