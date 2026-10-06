import { ref } from 'vue';
import { defineStore } from 'pinia';
import type { Tokens } from '@/types/game';
import 'pinia-plugin-persistedstate';

export interface ManagedToken extends Tokens {
  Method: string;
  Header: string;
  HeaderPrefix: string;
  Timer: ReturnType<typeof setTimeout> | null;
}

export const useLocalStore = defineStore(
  'localStore',
  () => {
    const clearLocalStorage = ref(false);
    const managedTokens = ref(<ManagedToken[]>[]);
    const version = ref(1.0);

    const addToken = (token: ManagedToken) => {
      managedTokens.value.push(token);
    };

    return {
      clearLocalStorage,
      managedTokens,
      version,
      addToken,
    };
  },
  {
    persist: {
      storage: localStorage,
      pick: ['clearLocalStorage', 'managedTokens', 'version'],
    },
  },
);

export default useLocalStore;
