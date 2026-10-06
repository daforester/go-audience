<template>
  <div />
</template>
<script lang="ts">
import { defineComponent } from 'vue';
import { SERVICE } from '@/mixins/api/GameRequest';

export default defineComponent({
  name: 'LogoutView',
  created() {
    this.emitter.emit('removeTokensByService', SERVICE);
    this.$router.push({ path: '/' });
  },
});
</script>
