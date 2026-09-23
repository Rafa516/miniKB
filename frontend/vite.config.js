import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  base: '/miniKB/',
  test: {
    environment: 'jsdom',
    setupFiles: './src/test/setup.js',
    globals: true,
    // isolate: false reaproveita o mesmo worker entre arquivos de teste, em
    // vez de criar um novo por arquivo — mais estável e bem mais rápido em
    // ambientes com poucos recursos (ex.: containers sandboxed usados para
    // rodar os testes durante o desenvolvimento).
    pool: 'threads',
    isolate: false,
  },
})
