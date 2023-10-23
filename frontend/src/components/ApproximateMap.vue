<template>
  <div id='map' style='height: 100%'></div>
</template>

<script setup lang='ts'>
  import { onMounted } from 'vue';

  import { useMaptilerToken } from '@/composables/useMaptilerToken';

  const props = defineProps<{
    longitude: number,
    latitude: number,
  }>();

  onMounted(async () => {
    const maptilerToken = await useMaptilerToken();
    const center = [props.longitude, props.latitude];

    const map = new maplibregl.Map({
      container: 'map',
      style: `https://api.maptiler.com/maps/streets/style.json?key=${maptilerToken}`,
      center,
      zoom: 10,
    });

    const marker = new maplibregl.Marker()
      .setLngLat(center)
      .addTo(map);
  });
</script>
