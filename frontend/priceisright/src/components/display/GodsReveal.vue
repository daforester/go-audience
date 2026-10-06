<template>
  <div class="d-flex justify-content-center flex-wrap">
    <div v-for="god in gods" :key="god.ID" class="m-4 toaster text-center">
      <div class="dance d-inline-block" title="God">
        <FontAwesomeIcon :icon="['fal', 'bolt']" />
      </div>
      <h1>
        {{ god.CharacterName }}
      </h1>
      <h5>
        God of
      </h5>
      <h2>
        {{ god.CharacterRace }}
      </h2>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

export default {
  name: 'GodsReveal',
  components: {
    FontAwesomeIcon,
  },
  props: {
    game: {
      type: Object,
      required: true,
    },
  },
  computed: {
    gods() {
      return this.game.Players.filter((p) => p.Class.toUpperCase() === 'GOD');
    },
  },
};

</script>

<style scoped>
@import '../../libs/scss/animations/dance.scss';
@import "../../libs/scss/animations/toaster.scss";
.toaster {
  animation-iteration-count: 1;
}
.text-center {
  font-size: 1.5rem;
}
</style>
