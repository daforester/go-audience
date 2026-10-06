<template>
  <div class="background vw-100 vh-100 position-absolute overflow-hidden">
    <div class="text-center">
      <div class="mt-5">
        <template v-for="(icon, index) in ['bolt', 'swords', 'bolt', 'swords', 'bolt']" :key="`bgIcon.1.${index}`">
          <FontAwesomeIcon 
                           :icon="['fal', icon]"
                           class="bg-icon dance text-light"
          />
        </template>
      </div>
      <div class="ml-5 pl-5">
        <template v-for="(icon, index) in ['mandolin', 'wand-magic', 'mandolin', 'wand-magic']" :key="`bgIcon.2.${index}`">
          <FontAwesomeIcon 
                           :icon="['fal', icon]"
                           class="bg-icon dance text-light"
          />
        </template>
      </div>
      <div>
        <template v-for="(icon, index) in ['bolt', 'swords', 'bolt', 'swords', 'bolt']" :key="`bgIcon.3.${index}`">
          <FontAwesomeIcon
                           :icon="['fal', icon]"
                           class="bg-icon dance text-light"
          />
        </template>
      </div>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

export default {
  name: 'IconDance',
  components: {
    FontAwesomeIcon,
  },
};

</script>

<style scoped lang="scss">
@import "../../libs/scss/animations/dance.scss";

.background {
  top: 0;
  left: 0;
  z-index: -100;
}

.bg-icon {
  font-size: 12vw;
  margin: 3rem;
}

</style>
