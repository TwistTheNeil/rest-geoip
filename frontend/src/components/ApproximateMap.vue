<template>
  <div id='map' style='height: 100%; width: 100%; flex: 1;'></div>
</template>

<script setup lang='ts'>
  import { onMounted, onUpdated, computed } from 'vue';
  import maplibregl from 'maplibre-gl';
  import { storeToRefs } from 'pinia';

  import { useConfigStore } from '@/stores/config';
  import { useMaxmindDataStore } from '@/stores/maxmindDataStore';

  const configStore = useConfigStore();
  const { maptilerToken } = storeToRefs(configStore);
  const maxmindDataStore = useMaxmindDataStore();
  const { data } = storeToRefs(maxmindDataStore);

  const center = computed(() => ({ lng: data.value!.Location.Longitude, lat: data.value!.Location.Latitude }));

  const updateMap = async () => {
    const map = new maplibregl.Map({
      container: 'map',
      style: `https://api.maptiler.com/maps/streets/style.json?key=${maptilerToken.value}`,
      center: center.value,
      zoom: 10,
    });

    new maplibregl.Marker()
      .setLngLat(center.value)
      .addTo(map);
  };

  onMounted(updateMap);
  onUpdated(updateMap);
</script>

<style>
  @import '@/../node_modules/maplibre-gl/dist/maplibre-gl.css';
</style>
