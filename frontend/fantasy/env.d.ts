/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string;
  readonly VITE_WS_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

import type { Emitter } from 'mitt';
import type { useLocalStore } from '@/stores/local';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Events = Record<string | symbol, any>;

declare module 'vue' {
  interface ComponentCustomProperties {
    emitter: Emitter<Events>;
    localStore: ReturnType<typeof useLocalStore>;
    apiUri: string;
    wsUri: string;
  }
}

export {};
