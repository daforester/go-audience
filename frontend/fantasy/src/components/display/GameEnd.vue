<template>
  <div class="endGame text-center">
    <div v-if="revealStep === 0" class="popper">
      <h1 class="award">
        And the awards go to...
      </h1>
    </div>
    <div v-else-if="revealStep === 1" class="popper">
      <h1 class="award">
        Best Audience Problem
      </h1>
      <template v-if="final.WinningSubmission">
        <TopicCard :topic="final.WinningSubmission" />
        <div class="h2 mt-3">
          <FontAwesomeIcon :icon="['fas', 'thumbs-up']" />
          {{ final.WinningSubmission.Net > 0 ? '+' : '' }}{{ final.WinningSubmission.Net }}
        </div>
      </template>
      <h2 v-else>
        No problems were played. Suspicious.
      </h2>
    </div>
    <div v-else class="popper">
      <h1 class="award">
        <FontAwesomeIcon :icon="['fas', 'crown']" class="crown dance" />
        {{ winners.length > 1 ? 'Joint Con Directors of the Year' : 'Con Directors of the Year' }}
      </h1>
      <div v-for="team in winners" :key="team.ID" class="winner" :style="{ color: store.teamColour(team.ID) }">
        {{ team.Name }}
        <div class="members h3">
          {{ team.Members.map((m) => m.Name).join(', ') }}
        </div>
      </div>
      <h2 v-if="winners.length === 0">
        Nobody won. The convention has been cancelled.
      </h2>
      <Scoreboard v-if="revealStep >= 3" :teams="game.Teams" :show-members="false" class="mt-4 toaster toaster-once" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import Scoreboard from './Scoreboard.vue';
import TopicCard from './TopicCard.vue';
import useGameStore from '@/stores/game';
import type { Final, Game, Team } from '@/types/game';

export default defineComponent({
  name: 'GameEnd',
  components: {
    FontAwesomeIcon,
    Scoreboard,
    TopicCard,
  },
  props: {
    game: {
      type: Object as PropType<Game>,
      required: true,
    },
  },
  data() {
    return {
      revealStep: 0,
      store: useGameStore(),
      timer: null as ReturnType<typeof setInterval> | null,
    };
  },
  computed: {
    final(): Final {
      return this.game.Final || { WinnerTeamIDs: [] };
    },
    winners(): Team[] {
      return this.final.WinnerTeamIDs.map((id) => this.store.teamById[id]).filter((t): t is Team => !!t);
    },
  },
  mounted() {
    // Reveal the awards one at a time
    this.timer = setInterval(() => {
      this.revealStep++;
      if (this.revealStep >= 3 && this.timer) {
        clearInterval(this.timer);
      }
    }, 5000);
  },
  unmounted() {
    if (this.timer) {
      clearInterval(this.timer);
    }
  },
});
</script>

<style scoped lang="scss">
@use '@/libs/scss/animations/dance.scss';
@use '@/libs/scss/animations/popper.scss';
@use '@/libs/scss/animations/toaster.scss';

.award {
  font-size: calc(1.5rem + 3vw);
  font-weight: 800;
  margin-bottom: 1.5rem;
}
.winner {
  font-size: calc(2rem + 4vw);
  font-weight: 800;
}
.members {
  color: #2b2233;
}
.crown {
  color: #f2a516;
  display: block;
  margin: 0 auto;
}
</style>
