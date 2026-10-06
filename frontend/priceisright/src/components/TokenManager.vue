<template>
  <div />
</template>
<script>
import Axios from "axios";
import useLocalStore from "@/stores/local";

export default {
  name: "TokenManager",
  props: {
    tokens: {
      type: Array,
      required: false,
      default() {
        return [];
      },
    },
  },
  data() {
    return {
      localStore: useLocalStore(),
      Tokens: [],
    };
  },
  emits: ["SessionExpired"],
  created() {
    if (
      this.localStore.clearLocalStorage === true &&
      !sessionStorage.getItem("clearLocalStorage")
    ) {
      this.localStore.clearLocalStorage = true;
    }

    this.replaceTokens(this.tokens);
    this.Tokens = this.localStore.managedTokens;
    this.Tokens.forEach((token) => {
      this.startTimer(token, 1);
    });
  },
  mounted() {
    this.emitter.on("refreshTokensByService", (service) => {
      const t = this.getTokenByService(service);
      if (t) {
        this.refreshToken(t);
      }
    });

    this.emitter.on("removeTokensByService", (service) => {
      this.removeTokensByService(service);
    });

    this.emitter.on("newToken", (token, clearLocalStorage) => {
      const newTokens = this.replaceTokens([token]);
      newTokens.forEach((t) => {
        this.startTimer(t);
      });

      if (
        clearLocalStorage !== null &&
        typeof clearLocalStorage === "boolean"
      ) {
        this.localStore.clearLocalStorage = clearLocalStorage;
      }
    });
  },
  methods: {
    clone(i) {
      return JSON.parse(JSON.stringify(i));
    },
    newToken() {
      return {
        AccessToken: "", // Data
        AccessExpires: 0, // Unix timestamp
        RefreshToken: "", // Data
        RefreshExpires: 0, // Unix time stamp
        Vendor: "", // 3rd party data
        Endpoint: "/api/token/refresh", // E.g. /api/token/refresh
        Method: "header", // Request method: header, post
        Header: "Authorization", // For header method, name of header
        HeaderPrefix: "Bearer ", // For header method, any prefix to token data
        Service: "", // Name of service the token is for
        Timer: null, // Internal Use
      };
    },
    getTokenByService(service) {
      return this.Tokens.find((t) => t.Service === service);
    },
    removeTokensByService(service) {
      this.Tokens.forEach((token) => {
        if (token.Service === service) {
          this.removeToken(token);
          this.$emit("SessionExpired", token);
        }
      });
    },
    parseTokenResponse(token, r) {
      let response = r;
      if (typeof r.data.Response !== "undefined") {
        response.data.Response ||= {};
        response = response.data.Response;
      }
      if (typeof response.Tokens !== "undefined") {
        response.Tokens ||= [];
        response.Tokens.forEach((data) => this.parseToken(token, data));
      }
      if (typeof response.Token !== "undefined") {
        if (typeof response.Token === "object") {
          this.parseToken(token, response.Token);
        } else {
          token.AccessToken = response.AccessToken;
          if (typeof response.Expires !== "undefined") {
            token.AccessExpires = response.Expires;
          }
        }
      }
      if (
        typeof response.AccessToken !== "undefined" &&
        typeof response.AccessExpires !== "undefined"
      ) {
        token.AccessToken = response.AccessToken;
        token.AccessExpires = response.AccessExpires;
      }
      if (
        typeof response.RefreshToken !== "undefined" &&
        typeof response.RefreshExpires !== "undefined"
      ) {
        token.RefreshToken = response.RefreshToken;
        token.RefreshExpires = response.RefreshExpires;
      }
      this.tokenUpdated(token);
    },
    parseToken(token, data) {
      data.Token ||= "";
      data.Expires ||= 0;
      data.Type ||= "Access";

      if (data.Type.toLowerCase() === "access") {
        token.AccessToken = data.Token;
        token.AccessExpires = data.Expires;
      } else if (data.Type.toLowerCase() === "refresh") {
        token.RefreshToken = data.Token;
        token.RefreshExpires = data.Expires;
      }
    },
    refreshToken(token) {
      if (token.Method === "post") {
        let data;
        if (token.RefreshToken) {
          data = token.RefreshToken;
        } else {
          data = token.AccessToken;
        }
        if (token.Header) {
          const jsonData = {};
          jsonData[token.Header] = data;
          data = jsonData;
        }
        return Axios({
          method: "post",
          url: token.Endpoint,
          withCredentials: true,
          data,
        })
          .then((r) => {
            this.parseTokenResponse(token, r);
          })
          .catch((e) => {
            if (typeof e.response === "object" && e.response.status === 401) {
              // Invalid Token
              this.removeToken(token);
              this.emitter.emit("sessionExpired", token);
            } else {
              this.retryTokenRefresh(token);
            }
          });
      }

      if (token.Method === "header") {
        const headers = {};
        if (token.RefreshToken) {
          headers[token.Header] = `${token.HeaderPrefix}${token.RefreshToken}`;
        } else {
          headers[token.Header] = `${token.HeaderPrefix}${token.AccessToken}`;
        }
        return Axios({
          method: "get",
          url: token.Endpoint,
          headers,
          withCredentials: true,
        })
          .then((r) => {
            this.parseTokenResponse(token, r);
          })
          .catch((e) => {
            if (typeof e.response === "object" && e.response.status === 401) {
              // Invalid Token
              this.removeToken(token);
              this.emitter.emit("sessionExpired", token);
            } else {
              this.retryTokenRefresh(token);
            }
          });
      }

      return new Promise((resolve, reject) => {
        reject(new Error("unsupported token method"));
      });
    },
    replaceTokens(tokens) {
      const newTokens = [];
      if (tokens.length > 0) {
        tokens.forEach((token) => {
          if (token.AccessToken !== "") {
            const newToken = this.newToken();
            Object.entries(newToken).forEach(([k]) => {
              if (typeof token[k] !== "undefined") {
                newToken[k] = token[k];
              }
            });

            const indexOf = this.Tokens.findIndex(
              (t) => t.Service === token.Service
            );
            if (indexOf >= 0) {
              this.removeToken(this.Tokens[indexOf]);
            }
            this.Tokens.push(newToken);
            newTokens.push(newToken);
          }
        });
      }

      return newTokens;
    },
    retryTokenRefresh(token) {
      const rms = this.refreshTimeRemaining(token);
      if (rms > 0) {
        // Update failed but Token still valid - Retry in 10 seconds
        token.Timer = setTimeout(() => {
          this.startTimer(token);
        }, 10000);
      } else {
        // Session Expired
        this.removeToken(token);
        this.emitter.emit("sessionExpired", token);
      }
    },
    startTimer(token, msOverride) {
      if (token.AccessExpires === 0 && token.RefreshExpires === 0) {
        this.removeToken(token);
        this.emitter.emit("SessionExpired", token);
        return;
      }

      // Refresh 8 seconds before expiry
      let ms = this.accessTimeRemaining(token);
      const rms = this.refreshTimeRemaining(token);

      if (ms < 0 && rms < 0) {
        this.removeToken(token);
        this.emitter.emit("SessionExpired", token);
        return;
      }

      // Renew before expires
      if (ms - 10 < 0) {
        ms = 1;
      }

      // Set Timer
      token.Timer = setTimeout(() => {
        this.refreshToken(token);
      }, msOverride || ms - 10);
    },
    accessTimeRemaining(token) {
      return token.AccessExpires * 1000 - new Date().getTime() - 8000;
    },
    refreshTimeRemaining(token) {
      return token.RefreshExpires * 1000 - new Date().getTime() - 8000;
    },
    removeToken(token) {
      if (token.Timer !== null) {
        clearTimeout(token.Timer);
      }
      const indexOf = this.Tokens.findIndex((t) => t === token);
      if (indexOf >= 0) {
        this.Tokens.splice(indexOf, 1);
      }
    },
    tokenUpdated(token, clearLocalStorage) {
      this.startTimer(token);
      if (
        clearLocalStorage !== null &&
        typeof clearLocalStorage === "boolean"
      ) {
        this.localStore.clearLocalStorage = clearLocalStorage;
      }

      this.emitter.emit("tokenUpdated", token);
    },
  },
};
</script>
