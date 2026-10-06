import { createApp } from 'vue';
import { createPinia } from 'pinia';
import persistedStatePlugin from 'pinia-plugin-persistedstate';

import App from '@/WebApp.vue';
import router from '@/router';

import mitt from 'mitt';
import useLocalStore from './stores/local';
import { resolveEndpoints } from './libs/endpoints';
const emitter = mitt();

const { apiUri, wsUri } = resolveEndpoints(import.meta.env.VITE_API_URL, import.meta.env.VITE_WS_URL);

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
