import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    outDir: '../internal/router/dist',
    // TODO: mapbox/maplibre-gl.js is a huge library and
    // needs to be worked upon to downsize
    // https://github.com/mapbox/mapbox-gl-js/issues/6320
    chunkSizeWarningLimit: 1000,
  },
})
