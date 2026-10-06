import 'es6-promise/auto';
import get from 'lodash.get';
import Vue from 'vue';
import Vuex from 'vuex';
import VueRouter from 'vue-router';
import { routes } from './routes';
import { config } from './store';
import WebApp from './views/WebApp.vue';

export default function Trudge(options = {}, componentConfig = {}) {
  Vue.use(Vuex);
  Vue.use(VueRouter);

  Vue.mixin({
    data() {
      return {
        apiUri: componentConfig.apiUri || API_URI,
        baseUri: componentConfig.baseUri || BASE_URI,
        wsUri: componentConfig.wsUri || WS_URI,
      };
    },
  });

  const store = new Vuex.Store(config);
  const router = new VueRouter({
    // eslint-disable-next-line no-undef
    base: componentConfig.baseUri || BASE_URI,
    mode: 'history',
    routes,
  });

  return new Vue({
    render: (h) => h(WebApp, componentConfig),
    router,
    store,
  }).$mount(get(options, 'el', '#webapp'));
}

if ('serviceWorker' in navigator) {
  navigator.serviceWorker
    // eslint-disable-next-line no-undef
    .register(`${BASE_URI}service-worker.js`)
    .then((registration) => {
      registration.onupdatefound = () => {
        if (navigator.serviceWorker.controller) {
          const installingWorker = registration.installing;
          installingWorker.onstatechange = () => {
            switch (installingWorker.state) {
              case 'installed':
                // New SW Installed
                break;
              case 'redundant':
                // SW Became Redundant
                break;
              default:
                // Other State Change
                break;
            }
          };
        }
      };
    }).catch(() => {
      // Discard
    });
}
