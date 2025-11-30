import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig, type UserConfigExport } from 'vite'

export default ({ mode }: { mode: string }): UserConfigExport => {
  // The dev server listens on port 8080, use it during development with vite
  if (!process.env.VITE_API_ENDPOINT) {
    if (mode === 'development') {
      process.env.VITE_API_ENDPOINT = 'http://localhost:8080/api/v1'
    } else {
      process.env.VITE_API_ENDPOINT = '/api/v1'
    }
  }

  return defineConfig({
    plugins: [react(), tailwindcss()],
    root: 'web',
    base: process.env.VITE_BASE_PATH,
    build: {
      outDir: '../internal/web/public',
      emptyOutDir: true,
      sourcemap: true,
    },
  })
}
