import { ref, reactive } from "vue";
import { defineStore } from "pinia";

export const useLocalStore = defineStore(
  "localStore",
  () => {
    const clearLocalStorage = ref(false);
    const managedTokens = ref(<object[]>[]);
    const version = ref(1.3);

    const addToken = (token: object) => {
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
      debug: true,
      storage: localStorage,
      paths: ["clearLocalStorage", "managedTokens", "version"],
    },
  }
);

export default useLocalStore;
