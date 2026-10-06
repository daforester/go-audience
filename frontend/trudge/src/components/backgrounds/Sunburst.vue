<template>
  <div class="background vw-100 vh-100 position-absolute overflow-hidden">
    <div class="spinner">
      <div class="sunburst">
        <div class="outer">
          <b v-for="x in 18" :key="x"></b>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

export default {
  name: 'Sunburst',
};

</script>

<style scoped lang="scss">
@import "../../libs/scss/animations/sunburst.scss";

.background {
  top: 0;
  left: 0;
  z-index: -100;
}

</style>
