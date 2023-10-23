<template>
  <div>
    <NavBar />
    <div class="container justify-content-center align-items-center">
      <div v-if="isReady">
        <SearchBar :currentIPAddress="data!.IP" @query-request="fetchIPAddressDetails" />
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
                      <span v-for="name in cityNames">
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
                <ApproximateMap :longitude="data?.Location.Longitude" :latitude="data?.Location.Latitude" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang='ts'>
  import { onMounted, ref, computed } from 'vue';
  import type { Ref } from 'vue';

  import NavBar from '@/components/NavBar.vue';
  import { useFetch_GetMaxmindData } from '@/composables/useFetch';
  import SearchBar from '@/components/SearchBar.vue';
  import type { MaxmindBackendResponse } from '@/types/MaxmindBackend';
  import ApproximateMap from '@/components/ApproximateMap.vue';

  const data:Ref<MaxmindBackendResponse | null> = ref(null);
  const error:Ref<string | null> = ref(null);

  const isReady = computed(() => {
    return data.value !== null;
  });

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

  const fetchCurrentIPAddressDetails = async () => {
    const { data: lata, error: lrror } = await useFetch_GetMaxmindData('/api/geoip');

    data.value = lata.value;
    error.value = lrror.value;

    if (error.value) {
      console.log(error.value);
    } else {
      console.log(data.value);
    }
  };

  const fetchIPAddressDetails = async (ipAddress: string) => {
    data.value = null;
    error.value = null;

    // TODO: check for valid ip address
    const { data: lata, error: lrror } = await useFetch_GetMaxmindData(`/api/geoip/${ipAddress}`);

    data.value = lata.value;
    error.value = lrror.value;


    if (error.value) {
      console.log(error.value);
    } else {
      console.log(data.value);
    }
  };

  onMounted(() => {
    fetchCurrentIPAddressDetails();
  });

</script>
