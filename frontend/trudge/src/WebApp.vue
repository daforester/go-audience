<template>
  <div id="WebAppBody">
    <div class="container">
      <TokenManager ref="TokenManager" />
      <router-view />
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
#WebAppBody {
  @import "@/libs/scss/variables.scss";
  @import "@/libs/scss/controls.scss";
  background: transparent;
  color: #333;
  font-family: sans-serif;
  font-size: 1.4em;
}
</style>
