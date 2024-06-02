import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/yanxi': {
        target: 'http://localhost:51233/',  // 要代理的目标地址
        changeOrigin: true,
         rewrite: (path) => path.replace(/^\/yanxi/, '')
      }
    }
  },
  resolve: {
    alias: {
      '@': resolve('./src')
    },

  }
})
