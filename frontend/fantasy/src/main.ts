import { createApp } from 'vue';
import { createPinia } from 'pinia';
import persistedStatePlugin from 'pinia-plugin-persistedstate';

import App from '@/WebApp.vue';
import router from '@/router';

import mitt from 'mitt';
import useLocalStore from './stores/local';
const emitter = mitt();

// Relative URLs (e.g. "/api") are served from the same host as the page
const apiUri = new URL(import.meta.env.VITE_API_URL, window.location.href).href.replace(/\/$/, '');
const wsUrl = new URL(import.meta.env.VITE_WS_URL, window.location.href);
if (wsUrl.protocol === 'http:' || wsUrl.protocol === 'https:') {
  wsUrl.protocol = wsUrl.protocol === 'https:' ? 'wss:' : 'ws:';
}
const wsUri = wsUrl.href.replace(/\/$/, '');

const app = createApp(App);

app.use(createPinia().use(persistedStatePlugin));
app.use(router);
app.provide('API_URL', apiUri);
app.provide('EMITTER', emitter);
app.provide('WS_URL', wsUri);
app.config.globalProperties.emitter = emitter;
app.config.globalProperties.localStore = useLocalStore();
app.config.globalProperties.apiUri = apiUri;
app.config.globalProperties.wsUri = wsUri;

app.mount('#app');
