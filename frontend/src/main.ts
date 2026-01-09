import './assets/main.css';
import 'primeicons/primeicons.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';
import PrimeVue from 'primevue/config';
import Material from '@primevue/themes/material';

import App from './App.vue';
import router from './router';

const pinia = createPinia();
const app = createApp(App);

app.use(PrimeVue, {
  ripple: true,
  theme: {
    preset: Material,
  },
});

// add eruda for debugging on mobile browsers
if (import.meta.env.DEV) {
  import('eruda').then((eruda) => eruda.default.init());
}

app.use(pinia);
app.use(router);

app.mount('#app');
