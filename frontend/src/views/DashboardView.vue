<template>
  <div>
    <Panel :header="`Details for IP Address: ${data!.IP}`">
      <div style="display: flex; flex-wrap: wrap; gap: 1rem; min-height: 500px;">
        <div style="flex: 1; min-width: 300px;">
          <DataTable :value="tableData" size="small" :showHeaders="false" :showGridlines="false">
            <Column field="property" />
            <Column field="value">
              <template #body="slotProps">
                <div v-if="slotProps.data.property === 'City'" style="white-space: pre-line;">{{ slotProps.data.value }}</div>
                <div v-else>{{ slotProps.data.value }}</div>
              </template>
            </Column>
          </DataTable>
        </div>
        <div style="flex: 1; min-width: 300px; min-height: 400px;">
          <ApproximateMap :longitude="data!.Location.Longitude" :latitude="data!.Location.Latitude" />
        </div>
      </div>
    </Panel>
  </div>
</template>

<script setup lang='ts'>
  import { computed } from 'vue';
  import type { Ref } from 'vue';
  import DataTable from 'primevue/datatable';
  import Column from 'primevue/column';
  import Panel from 'primevue/panel';

  import ApproximateMap from '@/components/ApproximateMap.vue';
  import { useMaxmindDataStore } from '@/stores/maxmindDataStore';
  import type { MaxmindBackendResponse } from '@/types/MaxmindBackend';
  import { storeToRefs } from 'pinia';

  const maxmindDataStore = useMaxmindDataStore();
  const { data }: { data: Ref<MaxmindBackendResponse | null> } = storeToRefs(maxmindDataStore);

  const cityNames = computed(() => {
    if (data.value === null) {
      return null;
    }

    if (!data.value.City.Names) {
      return null;
    }

    const names: Array<string> = Object.keys(data!.value.City.Names).map((e) => `${e}:${data!.value!.City.Names![e]}`);
    return names;
  });

  const tableData = computed(() => {
    if (!data.value) return [];

    return [
      { property: 'Country', value: data.value.Country.ISOCode },
      { property: 'EU', value: String(data.value.Country.IsInEuropeanUnion) },
      { property: 'City', value: cityNames.value?.join('\n') || '' },
      { property: 'Latitude', value: String(data.value.Location.Latitude) },
      { property: 'Longitude', value: String(data.value.Location.Longitude) },
      { property: 'Time Zone', value: data.value.Location.TimeZone },
      { property: 'Approx. Zip Code', value: data.value.Postal.Code },
      { property: 'Anonymous Proxy', value: String(data.value.Traits.IsAnonymousProxy) },
      { property: 'Satellite Provider', value: String(data.value.Traits.IsSatelliteProvider) },
      { property: 'Subdivision ISO Code', value: data.value.Subdivisions?.map(i => i.IsoCode)?.join(', ') || '' },
      { property: 'Subdivision Geo Name ID', value: data.value.Subdivisions?.map(i => i.GeoNameID)?.join(', ') || '' },
    ];
  });
</script>
