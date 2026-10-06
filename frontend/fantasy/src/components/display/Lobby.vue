<template>
  <div class="text-center lobby">
    <h1 class="title mb-3">
      Fantasy Furry Convention
    </h1>
    <div class="card-panel d-inline-block mb-4 px-5">
      <div class="h4 m-0">
        Join at <strong>{{ joinHost }}</strong>
      </div>
      <div class="room-code">
        {{ game.Key }}
      </div>
      <div class="h5 m-0">
        <FontAwesomeIcon :icon="['fas', 'paw']" />
        {{ audienceCount }} in the audience
        <template v-if="staffCount"> &middot; {{ staffCount }} staff</template>
      </div>
    </div>
    <Scoreboard v-if="game.Teams.length" :teams="game.Teams" />
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import Scoreboard from './Scoreboard.vue';
import type { Game } from '@/types/game';

export default defineComponent({
  name: 'LobbyDisplay',
  components: {
    FontAwesomeIcon,
    Scoreboard,
  },
  props: {
    game: {
      type: Object as PropType<Game>,
      required: true,
    },
  },
  computed: {
    joinHost(): string {
      return window.location.host;
    },
    audienceCount(): number {
      return (this.game.Players || []).filter((p) => p.Class === 'Audience').length;
    },
    staffCount(): number {
      return (this.game.Players || []).filter((p) => p.Class === 'Staff').length;
    },
  },
});
</script>

<style scoped lang="scss">
.title {
  color: #7b3fbf;
  font-size: calc(1.5rem + 3vw);
  font-weight: 800;
}
.room-code {
  font-size: calc(2rem + 6vw);
  font-weight: 800;
  letter-spacing: 0.15em;
}
</style>
