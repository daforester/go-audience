import { defineComponent } from 'vue';
import Axios, { type Method } from 'axios';
import useLocalStore from '@/stores/local';

export const SERVICE = 'Fantasy';

// GameRequest gives components an authorised call to the game API, as each Trudge
// component did by hand
export default defineComponent({
  data() {
    return {
      ErrorMessage: '',
      requesting: false,
    };
  },
  methods: {
    getToken() {
      return useLocalStore().managedTokens.find((token) => token.Service === SERVICE);
    },
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    gameRequest<T = any>(method: Method, path: string, data?: unknown): Promise<T> {
      const token = this.getToken();
      if (!token) {
        return Promise.reject(new Error('not in a game'));
      }

      this.ErrorMessage = '';
      this.requesting = true;

      return Axios({
        method,
        url: `${this.apiUri}${path}`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data,
      })
        .then((r) => {
          this.requesting = false;
          return r.data.Response as T;
        })
        .catch((e) => {
          this.requesting = false;
          if (typeof e.response === 'object' && e.response.data?.ErrorMessage) {
            this.ErrorMessage = e.response.data.ErrorMessage;
          } else {
            this.ErrorMessage = 'Could not reach the game, try again';
          }
          throw e;
        });
    },
  },
});
