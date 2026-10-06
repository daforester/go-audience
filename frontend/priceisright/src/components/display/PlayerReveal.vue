<template>
  <div class="d-flex justify-content-center flex-wrap">
    <div v-for="player in players" :key="player.ID" class="m-4 toaster text-center">
      <div v-if="player.Class.toUpperCase() === 'WARRIOR'"
           class="dance d-inline-block"
           title="Warrior"
      >
        <FontAwesomeIcon :icon="['fal', 'swords']" />
      </div>
      <div v-else-if="player.Class.toUpperCase() === 'BARD'"
           class="dance d-inline-block"
           title="Bard"
      >
        <FontAwesomeIcon :icon="['fal', 'mandolin']" />
      </div>
      <div v-else-if="player.Class.toUpperCase() === 'WIZARD'"
           class="dance d-inline-block"
           title="Wizard"
      >
        <FontAwesomeIcon :icon="['fal', 'wand-magic']" />
      </div>
      <div>
        <p class="popper name">
          {{ player.CharacterName }}
        </p>
        <p class="popper the">
          the
        </p>
        <p class="popper race">
          {{ player.CharacterRace }}
        </p>
        <p class="popper class">
          {{ player.Class }}
        </p>
      </div>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

export default {
  name: 'PlayerReveal',
  components: {
    FontAwesomeIcon,
  },
  props: {
    game: {
      type: Object,
      required: true,
    },
    type: {
      type: String,
      required: false,
      default() {
        return '';
      },
    },
  },
  computed: {
    players() {
      return this.game.Players.filter((p) => p.Class.toUpperCase() === this.type.toUpperCase());
    },
  },
};

</script>

<style scoped>
@import "../../libs/scss/animations/dance.scss";
@import "../../libs/scss/animations/popper.scss";

.text-center {
  font-size: 1.5rem;
}
.name {
  font-size: 3rem;
}
.the {
  font-size: 1.5rem;
}
.race {
  font-size: 2.5rem;
}
.class {
  font-size: 2rem;
}
</style>
