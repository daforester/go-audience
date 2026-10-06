import { createApp } from "vue";
import { createPinia } from "pinia";
import persistedStatePlugin from "pinia-plugin-persistedstate";

import App from "./WebApp.vue";
import router from "./router";

import mitt from "mitt";
const emitter = mitt();

const app = createApp(App);

app.use(createPinia().use(persistedStatePlugin));
app.use(router);
app.provide("API_URL", import.meta.env.VITE_API_URL);
app.provide("IMAGE_URL", import.meta.env.VITE_IMAGE_URL);
app.provide("WS_URL", import.meta.env.VITE_WS_URL);
app.config.globalProperties.emitter = emitter;
app.config.globalProperties.apiUri = import.meta.env.VITE_API_URL;
app.config.globalProperties.imageUri = import.meta.env.VITE_IMAGE_URL;
app.config.globalProperties.wsUri = import.meta.env.VITE_WS_URL;

app.mount("#app");
