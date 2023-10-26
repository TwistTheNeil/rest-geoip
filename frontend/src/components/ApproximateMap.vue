<template>
  <!-- TODO: figure out why height: 100% doesn't work on mobile -->
  <div id='map' style='height: 500px'></div>
</template>

<script setup lang='ts'>
  import { onMounted } from 'vue';
  import maplibregl from 'maplibre-gl';

  import { useMaptilerToken } from '@/composables/useMaptilerToken';

  const props = defineProps<{
    longitude: number,
    latitude: number,
  }>();

  onMounted(async () => {
    const maptilerToken = await useMaptilerToken();
    const center = { lng: props.longitude, lat: props.latitude};

    const map = new maplibregl.Map({
      container: 'map',
      style: `https://api.maptiler.com/maps/streets/style.json?key=${maptilerToken}`,
      center,
      zoom: 10,
    });

    new maplibregl.Marker()
      .setLngLat(center)
      .addTo(map);
  });
</script>

<style>
  @import '@/../node_modules/maplibre-gl/dist/maplibre-gl.css';
</style>
