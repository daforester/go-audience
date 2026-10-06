<template>
  <div class="text-center lobby">
    <h1 class="title mb-3">
      Fantasy Furry Convention
    </h1>
    <div class="card-panel join mb-4">
      <JoinQr :room-key="game.Key" :size="320" class="lobby-qr" />
      <div>
        <div class="h4 m-0">
          Scan, or join at <strong>{{ joinHost }}</strong>
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
    </div>
    <Scoreboard v-if="game.Teams.length" :teams="game.Teams" :show-wins="!!game.Round" large />
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import JoinQr from './JoinQr.vue';
import Scoreboard from './Scoreboard.vue';
import { joinHost } from '@/libs/join';
import type { Game } from '@/types/game';

export default defineComponent({
  name: 'LobbyDisplay',
  components: {
    FontAwesomeIcon,
    JoinQr,
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
      return joinHost();
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
.join {
  align-items: center;
  display: inline-flex;
  gap: 2rem;
  padding: 1rem 2.5rem 1rem 1rem;
  text-align: left;
}
// About a quarter of the screen height: big enough to scan from the back of the room
.lobby-qr {
  height: clamp(160px, 26vh, 320px);
  width: auto;
}
.room-code {
  font-size: calc(2rem + 6vw);
  font-weight: 800;
  letter-spacing: 0.15em;
}
</style>
