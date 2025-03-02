import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig({
  plugins: [react()],
  root: resolve('./static'),
  build: {
    outDir: resolve('./static/dist'),
    emptyOutDir: true,
    assetsInlineLimit: 0,
  },
  publicDir: resolve('./static/img'),
  server: {
    proxy: {
      '/api': 'http://localhost:3000',
    },
  },
});