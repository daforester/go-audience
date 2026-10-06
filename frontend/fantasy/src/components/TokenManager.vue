<template>
  <div />
</template>
<script lang="ts">
/* eslint-disable @typescript-eslint/no-explicit-any */
import { defineComponent, type PropType } from 'vue';
import Axios, { type AxiosResponse } from 'axios';
import useLocalStore, { type ManagedToken } from '@/stores/local';

export default defineComponent({
  name: 'TokenManager',
  props: {
    tokens: {
      type: Array as PropType<ManagedToken[]>,
      required: false,
      default() {
        return [];
      },
    },
  },
  emits: ['sessionExpired'],
  data() {
    return {
      localStore: useLocalStore(),
      Tokens: [] as ManagedToken[],
    };
  },
  created() {
    if (this.localStore.clearLocalStorage === true && !sessionStorage.getItem('clearLocalStorage')) {
      this.localStore.clearLocalStorage = true;
    }

    this.replaceTokens(this.tokens);
    this.Tokens = this.localStore.managedTokens;
    this.Tokens.forEach((token) => {
      this.startTimer(token, 1);
    });
  },
  mounted() {
    this.emitter.on('refreshTokensByService', (service: string) => {
      const t = this.getTokenByService(service);
      if (t) {
        this.refreshToken(t);
      }
    });

    this.emitter.on('removeTokensByService', (service: string) => {
      this.removeTokensByService(service);
    });

    this.emitter.on('newToken', (token: ManagedToken, clearLocalStorage?: boolean) => {
      const newTokens = this.replaceTokens([token]);
      newTokens.forEach((t) => {
        this.startTimer(t);
      });

      if (clearLocalStorage !== null && typeof clearLocalStorage === 'boolean') {
        this.localStore.clearLocalStorage = clearLocalStorage;
      }
    });
  },
  beforeUnmount() {
    this.emitter.off('refreshTokensByService');
    this.emitter.off('removeTokensByService');
    this.emitter.off('newToken');
    this.Tokens.forEach((token) => {
      if (token.Timer !== null) {
        clearTimeout(token.Timer);
      }
    });
  },
  methods: {
    newToken(): ManagedToken {
      return {
        AccessToken: '', // Data
        AccessExpires: 0, // Unix timestamp
        RefreshToken: '', // Data
        RefreshExpires: 0, // Unix time stamp
        Vendor: '', // 3rd party data
        Endpoint: '/api/token/refresh', // E.g. /api/token/refresh
        Method: 'header', // Request method: header, post
        Header: 'Authorization', // For header method, name of header
        HeaderPrefix: 'Bearer ', // For header method, any prefix to token data
        Service: '', // Name of service the token is for
        Timer: null, // Internal Use
      };
    },
    getTokenByService(service: string) {
      return this.Tokens.find((t) => t.Service === service);
    },
    removeTokensByService(service: string) {
      [...this.Tokens].forEach((token) => {
        if (token.Service === service) {
          this.removeToken(token);
          this.$emit('sessionExpired', token);
        }
      });
    },
    parseTokenResponse(token: ManagedToken, r: AxiosResponse) {
      let response: any = r;
      if (typeof r.data.Response !== 'undefined') {
        response = r.data.Response || {};
      }
      if (typeof response.Tokens !== 'undefined') {
        (response.Tokens || []).forEach((data: any) => this.parseToken(token, data));
      }
      if (typeof response.Token !== 'undefined') {
        if (typeof response.Token === 'object') {
          this.parseToken(token, response.Token);
        } else {
          token.AccessToken = response.AccessToken;
          if (typeof response.Expires !== 'undefined') {
            token.AccessExpires = response.Expires;
          }
        }
      }
      if (typeof response.AccessToken !== 'undefined' && typeof response.AccessExpires !== 'undefined') {
        token.AccessToken = response.AccessToken;
        token.AccessExpires = response.AccessExpires;
      }
      if (
        typeof response.RefreshToken !== 'undefined' &&
        typeof response.RefreshExpires !== 'undefined'
      ) {
        token.RefreshToken = response.RefreshToken;
        token.RefreshExpires = response.RefreshExpires;
      }
      this.tokenUpdated(token);
    },
    parseToken(token: ManagedToken, data: any) {
      data.Token ||= '';
      data.Expires ||= 0;
      data.Type ||= 'Access';

      if (data.Type.toLowerCase() === 'access') {
        token.AccessToken = data.Token;
        token.AccessExpires = data.Expires;
      } else if (data.Type.toLowerCase() === 'refresh') {
        token.RefreshToken = data.Token;
        token.RefreshExpires = data.Expires;
      }
    },
    refreshFailed(token: ManagedToken, e: any) {
      if (typeof e.response === 'object' && e.response.status === 401) {
        // Invalid Token
        this.removeToken(token);
        this.emitter.emit('sessionExpired', token);
      } else {
        this.retryTokenRefresh(token);
      }
    },
    refreshToken(token: ManagedToken) {
      if (token.Method === 'post') {
        let data: any = token.RefreshToken ? token.RefreshToken : token.AccessToken;
        if (token.Header) {
          const jsonData: Record<string, string> = {};
          jsonData[token.Header] = data;
          data = jsonData;
        }
        return Axios({
          method: 'post',
          url: token.Endpoint,
          withCredentials: true,
          data,
        })
          .then((r) => this.parseTokenResponse(token, r))
          .catch((e) => this.refreshFailed(token, e));
      }

      if (token.Method === 'header') {
        const headers: Record<string, string> = {};
        if (token.RefreshToken) {
          headers[token.Header] = `${token.HeaderPrefix}${token.RefreshToken}`;
        } else {
          headers[token.Header] = `${token.HeaderPrefix}${token.AccessToken}`;
        }
        return Axios({
          method: 'get',
          url: token.Endpoint,
          headers,
          withCredentials: true,
        })
          .then((r) => this.parseTokenResponse(token, r))
          .catch((e) => this.refreshFailed(token, e));
      }

      return Promise.reject(new Error('unsupported token method'));
    },
    replaceTokens(tokens: ManagedToken[]) {
      const newTokens: ManagedToken[] = [];
      tokens.forEach((token) => {
        if (token.AccessToken !== '') {
          const newToken: any = this.newToken();
          Object.keys(newToken).forEach((k) => {
            if (typeof (token as any)[k] !== 'undefined') {
              newToken[k] = (token as any)[k];
            }
          });
          newToken.Timer = null;

          const existing = this.Tokens.find((t) => t.Service === token.Service);
          if (existing) {
            this.removeToken(existing);
          }
          this.Tokens.push(newToken);
          newTokens.push(newToken);
        }
      });

      return newTokens;
    },
    retryTokenRefresh(token: ManagedToken) {
      const rms = this.refreshTimeRemaining(token);
      if (rms > 0) {
        // Update failed but Token still valid - Retry in 10 seconds
        token.Timer = setTimeout(() => {
          this.startTimer(token);
        }, 10000);
      } else {
        // Session Expired
        this.removeToken(token);
        this.emitter.emit('sessionExpired', token);
      }
    },
    startTimer(token: ManagedToken, msOverride?: number) {
      if (token.AccessExpires === 0 && token.RefreshExpires === 0) {
        this.removeToken(token);
        this.emitter.emit('sessionExpired', token);
        return;
      }

      // Refresh 8 seconds before expiry
      let ms = this.accessTimeRemaining(token);
      const rms = this.refreshTimeRemaining(token);

      if (ms < 0 && rms < 0) {
        this.removeToken(token);
        this.emitter.emit('sessionExpired', token);
        return;
      }

      // Renew before expires
      if (ms - 10 < 0) {
        ms = 11;
      }

      if (token.Timer !== null) {
        clearTimeout(token.Timer);
      }

      // Set Timer
      token.Timer = setTimeout(
        () => {
          this.refreshToken(token);
        },
        msOverride || ms - 10,
      );
    },
    accessTimeRemaining(token: ManagedToken) {
      return token.AccessExpires * 1000 - new Date().getTime() - 8000;
    },
    refreshTimeRemaining(token: ManagedToken) {
      return token.RefreshExpires * 1000 - new Date().getTime() - 8000;
    },
    removeToken(token: ManagedToken) {
      if (token.Timer !== null) {
        clearTimeout(token.Timer);
      }
      const indexOf = this.Tokens.findIndex((t) => t === token);
      if (indexOf >= 0) {
        this.Tokens.splice(indexOf, 1);
      }
    },
    tokenUpdated(token: ManagedToken) {
      this.startTimer(token);
      this.emitter.emit('tokenUpdated', token);
    },
  },
});
</script>
