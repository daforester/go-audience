<template>
  <div class="text-center">
    <div class="dance display-3 mb-3">
      <FontAwesomeIcon :icon="['fas', icon]" />
    </div>
    <div v-if="name" class="card-panel h2 mb-3">
      {{ name }}
    </div>
    <slot />
    <DisplayMessage effect="toaster" :message="message" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import DisplayMessage from '../display/DisplayMessage.vue';

export default defineComponent({
  name: 'PleaseWait',
  components: {
    DisplayMessage,
    FontAwesomeIcon,
  },
  props: {
    icon: {
      type: String,
      required: false,
      default: 'paw',
    },
    name: {
      type: String,
      required: false,
      default: '',
    },
    message: {
      type: String,
      required: false,
      default: 'Please Wait...',
    },
  },
});
</script>

<style scoped lang="scss">
@use '@/libs/scss/animations/dance.scss';
</style>
