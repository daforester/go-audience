<template>
  <h2 class="text-center">
    <template v-if="effect">
      <template v-if="!wholeWord && !perWord">
        <span v-for="(letter, k) in splitMessage"
              :key="k"
              :class="effect"
        ><template v-if="letter === ' '">
      &nbsp;
        </template><template v-else>{{ letter }}</template></span>
      </template>
      <template v-else-if="!wholeWord">
        <span v-for="(word, k) in splitWord"
              :key="k"
              :class="effect"
        >
          {{ word }}
        </span>
      </template>
      <template v-else>
        <span :class="effect">{{ message }}</span>
      </template>
    </template>
    <template v-else>
      {{ message }}
    </template>
  </h2>
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
    perWord: {
      type: Boolean,
      required: false,
      default() {
        return false;
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
    splitWord() {
      return this.message.split(' ');
    },
  },
};

</script>

<style lang="scss" scoped>
@import "../../libs/scss/animations/burst.scss";
@import "../../libs/scss/animations/comedown.scss";
@import "../../libs/scss/animations/fader.scss";
@import "../../libs/scss/animations/ripple.scss";
@import "../../libs/scss/animations/shimmer.scss";
@import "../../libs/scss/animations/suggestion.scss";
@import "../../libs/scss/animations/toaster.scss";

h2 span {
  display: inline-block;
  vertical-align: middle;
}

.smaller {
  margin: 0;
  padding: 0;
}

.smaller span {
  display: inline;
  font-size: 6rem;
  padding: 0;
}
</style>
