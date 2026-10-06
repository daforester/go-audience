<template>
  <div>
    <h1 class="text-center">
      <template v-if="effect">
        <template v-if="!wholeWord">
          <span v-for="(letter, k) in splitMessage"
                :key="k"
                :class="effect"
          ><template v-if="letter === ' '">
        &nbsp;
          </template><template v-else>{{ letter }}</template></span>
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

<script>
export default {
  name: 'DisplayMessage',
  props: {
    effect: {
      type: String,
      required: false,
      default() {
        return '';
      },
    },
    message: {
      type: String,
      required: false,
      default() {
        return '';
      },
    },
    wholeWord: {
      type: Boolean,
      required: false,
      default() {
        return false;
      },
    },
  },
  computed: {
    splitMessage() {
      return this.message.split('');
    },
  },
};

</script>

<style scoped lang="scss">
@import "../../libs/scss/animations/burst.scss";
@import "../../libs/scss/animations/fader.scss";
@import "../../libs/scss/animations/ripple.scss";
@import "../../libs/scss/animations/shimmer.scss";
@import "../../libs/scss/animations/suggestion.scss";
@import "../../libs/scss/animations/toaster.scss";
h1 {
  font-size: calc(1.375rem + 3vw) !important;
}
</style>
