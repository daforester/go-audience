import { createRouter, createWebHistory } from "vue-router";
import Game from "@/views/Game.vue";
import Logout from "@/views/Logout.vue";

const Advert = () => import("@/views/Advert.vue");
const MapMaker = () => import("@/views/MapMaker.vue");

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: Game,
    },
    {
      path: "/advert",
      component: Advert,
    },
    {
      path: "/logout",
      component: Logout,
    },
    {
      path: "/mapmaker",
      component: MapMaker,
    },
  ],
});

export default router;
