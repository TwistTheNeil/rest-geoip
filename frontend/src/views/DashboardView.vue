<template>
  <div>
    <div class="progress" style="height: 3px;">
      <div class="progress-bar" role="progressbar" :style="progressValueWidth" :aria-valuenow="progressValue" aria-valuemin="0" aria-valuemax="100"></div>
    </div>
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
                <ApproximateMap :longitude="data!.Location.Longitude" :latitude="data!.Location.Latitude" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang='ts'>
  import { onMounted, ref, computed, watch } from 'vue';
  import type { Ref } from 'vue';
  import { useRoute } from 'vue-router';
  import { isIP } from 'is-ip';

  import NavBar from '@/components/NavBar.vue';
  import { useFetch_GetMaxmindData } from '@/composables/useFetch';
  import SearchBar from '@/components/SearchBar.vue';
  import type { MaxmindBackendResponse } from '@/types/MaxmindBackend';
  import ApproximateMap from '@/components/ApproximateMap.vue';

  const route = useRoute();
  const data: Ref<MaxmindBackendResponse | null> = ref(null);
  const error: Ref<string | null> = ref(null);
  const progressValue: Ref<string> = ref("0");

  const isReady = computed(() => {
    return data.value !== null;
  });

  const progressValueWidth = computed(() => {
    return `width: ${progressValue.value}%;`
  });

  watch(progressValue, async (newProgressValue) => {
    if (newProgressValue === '100') {
      setTimeout(() => progressValue.value = '0', 1000);
    }
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
    progressValue.value = '0';
    const { data: lata, error: lrror } = await useFetch_GetMaxmindData('/api/geoip');

    progressValue.value = '60';

    if (lrror.value) {
      console.log(lrror.value);
      return;
    } else {
      console.log(lata.value);
    }

    progressValue.value = '80';

    data.value = lata.value;
    error.value = lrror.value;

    progressValue.value = '100';
  };

  const fetchIPAddressDetails = async (ipAddress: string) => {
    if (!isIP(ipAddress)) {
      return;
    }

    progressValue.value = '0';

    data.value = null;
    error.value = null;

    const { data: lata, error: lrror } = await useFetch_GetMaxmindData(`/api/geoip/${ipAddress}`);

    progressValue.value = '60';

    if (lrror.value) {
      console.log(lrror.value);
      return;
    } else {
      console.log(lata.value);
    }

    progressValue.value = '80';
    console.log(progressValue.value);

    data.value = lata.value;
    error.value = lrror.value;

    progressValue.value = '100';
    console.log(progressValue.value);
  };

  onMounted(() => {
    if (!route.query.address || route.query.address === "") {
      fetchCurrentIPAddressDetails();
    } else {
      fetchIPAddressDetails(route.query.address as string);
    }
  });

</script>
