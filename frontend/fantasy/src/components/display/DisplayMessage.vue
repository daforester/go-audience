<template>
  <div>
    <h1 class="text-center">
      <template v-if="effect">
        <template v-if="!wholeWord">
          <span v-for="(letter, k) in splitMessage"
                :key="k"
                :class="effect"
          ><template v-if="letter === ' '">&nbsp;</template><template v-else>{{ letter }}</template></span>
        </template>
        <template v-else>
          <span :class="effect">{{ message }}</span>
        </template>
      </template>
      <template v-else>
        {{ message }}
      </template>
    </h1>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'DisplayMessage',
  props: {
    effect: {
      type: String,
      required: false,
      default: '',
    },
    message: {
      type: String,
      required: false,
      default: '',
    },
    wholeWord: {
      type: Boolean,
      required: false,
      default: false,
    },
  },
  computed: {
    splitMessage(): string[] {
      return this.message.split('');
    },
  },
});
</script>

<style scoped lang="scss">
@use '@/libs/scss/animations/burst.scss';
@use '@/libs/scss/animations/fader.scss';
@use '@/libs/scss/animations/ripple.scss';
@use '@/libs/scss/animations/shimmer.scss';
@use '@/libs/scss/animations/toaster.scss';
h1 {
  font-size: calc(1.375rem + 3vw) !important;
}
</style>
