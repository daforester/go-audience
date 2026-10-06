<template>
  <div>
    <DisplayMessage effect="ripple"
                    message="Please Wait"
    />
  </div>
</template>

<script>
import DisplayMessage from '../display/DisplayMessage.vue';

export default {
  name: 'PleaseWait',
  components: {
    DisplayMessage,
  },
  props: {
    gameCode: {
      type: Number,
      required: true,
    },
    player: {
      type: Object,
      required: true,
    },
  },
  methods: {
    formatArray(a) {
      if (a) {
        return a.toString();
      }
      return '';
    },
  },
};

</script>

<style scoped>
</style>
