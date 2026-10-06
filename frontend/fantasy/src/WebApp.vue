<template>
  <div id="WebAppBody">
    <div class="container-fluid p-0">
      <TokenManager ref="TokenManager" />
      <router-view />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import '@/libs/fontawesome/minimal';
import TokenManager from '@/components/TokenManager.vue';

export default defineComponent({
  name: 'WebApp',
  components: {
    TokenManager,
  },
});
</script>

<style lang="scss">
@use '@/libs/scss/app.scss';

#WebAppBody {
  background: transparent;
  color: #2b2233;
  font-size: 1.2em;
  min-height: 100vh;
}
</style>
