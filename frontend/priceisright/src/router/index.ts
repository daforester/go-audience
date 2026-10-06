import { createRouter, createWebHistory } from "vue-router";
import Game from "@/views/Game.vue";
import Logout from "@/views/Logout.vue";

//const Game = () => import("@/views/Game.vue");
//const Logout = () => import("@/views/Logout.vue");

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      component: Game,
    },
    {
      path: "/logout",
      component: Logout,
    },
  ],
});

export default router;
