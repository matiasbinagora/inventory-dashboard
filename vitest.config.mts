// The .mts extension makes this ESM configuration explicit while package.json
// remains CommonJS by default.
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: { environment: 'jsdom', setupFiles: ['./test/setup.ts'], include: ['features/**/*.test.ts', 'features/**/*.test.tsx'] },
})
