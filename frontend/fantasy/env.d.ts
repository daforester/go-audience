/// <reference types="vite/client" />

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
