import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import { endpointErrors } from './src/libs/endpoints.ts'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // API endpoints are build configuration: from the environment, .env.[mode] or
  // .env.[mode].local. Never fall back to a hardcoded address.
  const env = { ...loadEnv(mode, process.cwd(), 'VITE_'), ...process.env }
  const errors = endpointErrors(env.VITE_API_URL, env.VITE_WS_URL)
  if (errors.length > 0 && mode !== 'test') {
    throw new Error(
      `Fantasy endpoints are not configured for "${mode}":\n  - ${errors.join('\n  - ')}\n` +
        'Set VITE_API_URL and VITE_WS_URL (see .env.production.example).',
    )
  }

  return {
    plugins: [
      vue(),
      vueDevTools(),
    ],
    server: {
      // Bind mounts from Windows/macOS hosts do not deliver file events into containers;
      // the dev compose file turns polling on
      watch: process.env.FANTASY_DEV_POLLING ? { usePolling: true, interval: 300 } : undefined,
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    css: {
      preprocessorOptions: {
        scss: {
          // Bootstrap 5.3 still uses @import internally; keep its deprecation noise quiet
          quietDeps: true,
          silenceDeprecations: ['import', 'global-builtin', 'color-functions', 'if-function'],
        },
      },
    },
  }
})
