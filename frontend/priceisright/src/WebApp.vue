<template>
  <div id="WebAppBody">
    <div class="">
      <TokenManager ref="TokenManager"/>
      <router-view/>
    </div>
  </div>
</template>

<script>
import "@/libs/fontawesome/minimal";
import TokenManager from "@/components/TokenManager.vue";

export default {
  name: "WebApp",
  components: {
    TokenManager,
  },
};
</script>

<style lang="scss">
@font-face {
  font-family: "Hollywood Hills";
  src: url("@/assets/fonts/HollywoodHills.woff2") format("woff2"),
  url('@/assets/fonts/HollywoodHills.woff') format('woff');
  font-weight: normal;
  font-style: normal;
  font-display: swap;
}

@font-face {
  font-family: "LEMON MILK";
  src: url("@/assets/fonts/LEMONMILK-Bold.woff2") format("woff2"),
  url('@/assets/fonts/LEMONMILK-Bold.woff') format('woff');
  font-weight: bold;
  font-style: normal;
  font-display: swap;
}

body {
  margin: 0;
  padding: 0;
}

#WebAppBody {
  @import "@/libs/scss/variables.scss";
  @import "@/libs/scss/controls.scss";
  background: transparent;
  color: #999;
  font-family: sans-serif;
  font-size: 1.4em;
}
</style>
