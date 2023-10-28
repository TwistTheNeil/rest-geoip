<template>
  <div>
    <div class="card">
      <div class="card-header">
        Details for IP Address: {{ data!.IP }}
      </div>
      <div class="card-body">
        <div class="row">
          <div class="col-lg-6">
            <table class="table table-borderless table-hover table-sm">
            <tbody>
              <tr> <td> Country </td> <td>{{ data!.Country.ISOCode }}</td>  </tr>
              <tr> <td> EU </td> <td>{{ data!.Country.IsInEuropeanUnion }}</td> </tr>
              <tr>
                <td> City </td>
                <td>
                  <span v-for="name in cityNames" :key="name">
                    {{ name }}
                    <br />
                  </span>
                </td>
              </tr>
              <tr> <td> Latitude </td> <td>{{ data!.Location.Latitude }}</td> </tr>
              <tr> <td> Longitude </td> <td>{{ data!.Location.Longitude }}</td> </tr>
              <tr> <td> Time Zone </td> <td>{{ data!.Location.TimeZone }}</td> </tr>
              <tr> <td> Approx. Zip Code </td> <td>{{ data!.Postal.Code }}</td> </tr>
              <tr> <td> Anonymous Proxy </td> <td>{{ data!.Traits.IsAnonymousProxy }}</td> </tr>
              <tr> <td> Satellite Provider </td> <td>{{ data!.Traits.IsSatelliteProvider }}</td> </tr>
              </tbody>
            </table>
          </div>
          <div class="col-lg-6">
            <ApproximateMap :longitude="data!.Location.Longitude" :latitude="data!.Location.Latitude" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang='ts'>
  import { computed } from 'vue';
  import type { Ref } from 'vue';

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

    const names: Array<string> = Object.keys(data!.value.City.Names).map((e) => `${e}:${data!.value!.City.Names[e]}`);
    return names;
  });
</script>
