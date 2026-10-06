<template>
  <div class="scoreboard" :style="{ '--teams': teams.length }">
    <div v-for="team in teams"
         :key="team.ID"
         class="team card-panel"
         :class="{ presenting: presenting === team.ID }"
         :style="{ borderColor: store.teamColour(team.ID) }"
    >
      <div class="name" :style="{ color: store.teamColour(team.ID) }">
        <FontAwesomeIcon v-if="presenting === team.ID" :icon="['fas', 'microphone']" />
        {{ team.Name }}
      </div>
      <div v-if="showMembers" class="members small">
        {{ team.Members.map((m) => m.Name).join(', ') }}
      </div>
      <div class="wins">
        <FontAwesomeIcon :icon="['fas', 'trophy']" />
        {{ team.RoundWins }}
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import useGameStore from '@/stores/game';
import type { Team } from '@/types/game';

export default defineComponent({
  name: 'TeamScoreboard',
  components: {
    FontAwesomeIcon,
  },
  props: {
    teams: {
      type: Array as PropType<Team[]>,
      required: true,
    },
    presenting: {
      type: String,
      default: '',
    },
    showMembers: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      store: useGameStore(),
    };
  },
});
</script>

<style scoped lang="scss">
.scoreboard {
  display: grid;
  gap: 1rem;
  grid-template-columns: repeat(var(--teams), minmax(0, 1fr));
  margin: 0 auto;
  max-width: min(94vw, 70rem);
}
.team {
  border-width: 3px;
  text-align: center;
}
.name {
  font-size: 1.4em;
  font-weight: 800;
  overflow-wrap: anywhere;
}
.wins {
  font-size: 1.6em;
  font-weight: 800;
}
.presenting {
  transform: scale(1.06);
  transition: transform 0.3s;
}
</style>
