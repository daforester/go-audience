<template>
  <div class="team-bars" :style="{ '--teams': teams.length, '--bar-height': height }">
    <div v-for="team in teams" :key="team.ID" class="team" :class="{ winner: winners.includes(team.ID) }">
      <div class="bar-wrap">
        <div class="bar"
             :style="{ height: `${percent(team.ID)}%`, background: store.teamColour(team.ID) }"
             role="progressbar"
             :aria-valuenow="votes[team.ID] || 0"
             aria-valuemin="0"
             :aria-valuemax="total"
             :aria-label="`${team.Name} votes`"
        />
      </div>
      <div class="count">
        {{ votes[team.ID] || 0 }}
      </div>
      <div class="name">
        <FontAwesomeIcon v-if="winners.includes(team.ID)" :icon="['fas', 'crown']" class="crown" />
        {{ team.Name }}
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import useGameStore from '@/stores/game';
import type { Team } from '@/types/game';

// TeamBars shows one vote bar per team; laid out for any number of teams
export default defineComponent({
  name: 'TeamBars',
  components: {
    FontAwesomeIcon,
  },
  props: {
    teams: {
      type: Array as PropType<Team[]>,
      required: true,
    },
    votes: {
      type: Object as PropType<Record<string, number>>,
      required: true,
    },
    winners: {
      type: Array as PropType<string[]>,
      default: () => [],
    },
    height: {
      type: String,
      default: '40vh',
    },
  },
  data() {
    return {
      store: useGameStore(),
    };
  },
  computed: {
    total(): number {
      return Object.values(this.votes).reduce((a, b) => a + b, 0);
    },
  },
  methods: {
    percent(teamID: string): number {
      if (this.total === 0) {
        return 0;
      }
      return Math.round(((this.votes[teamID] || 0) / this.total) * 100);
    },
  },
});
</script>

<style scoped lang="scss">
.team-bars {
  display: grid;
  gap: 2vw;
  grid-template-columns: repeat(var(--teams), minmax(0, 1fr));
  margin: 0 auto;
  max-width: min(90vw, 60rem);
}
.team {
  text-align: center;
}
.bar-wrap {
  align-items: flex-end;
  background: rgba(255, 255, 255, 0.7);
  border: 2px solid #2b2233;
  border-radius: 0.8rem;
  display: flex;
  height: var(--bar-height);
  overflow: hidden;
}
.bar {
  transition: height 0.4s ease-out;
  width: 100%;
}
.count {
  font-size: calc(1rem + 2vw);
  font-weight: 800;
}
.name {
  font-size: calc(0.8rem + 1vw);
  font-weight: 700;
  overflow-wrap: anywhere;
}
.winner .bar-wrap {
  box-shadow: 0 0 2rem #f2a516;
}
.crown {
  color: #f2a516;
}
</style>
