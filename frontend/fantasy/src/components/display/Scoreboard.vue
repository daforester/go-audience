<template>
  <div class="scoreboard" :class="{ large }" :style="{ '--teams': teams.length }">
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
      <div v-if="showWins" class="wins">
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
    showWins: {
      type: Boolean,
      default: true,
    },
    // Sized for reading across a room on the main display
    large: {
      type: Boolean,
      default: false,
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
.large {
  max-width: min(94vw, 110rem);
  .name {
    font-size: calc(1rem + 1.4vw);
  }
  .members {
    font-size: calc(0.8rem + 0.7vw) !important;
  }
  .wins {
    font-size: calc(1rem + 1.6vw);
  }
}
.presenting {
  transform: scale(1.06);
  transition: transform 0.3s;
}

// Phones: one team per row, name on the left and wins on the right
@media (max-width: 576px) {
  .scoreboard {
    gap: 0.5rem;
    grid-template-columns: minmax(0, 1fr);
  }
  .team {
    align-items: center;
    column-gap: 0.75rem;
    display: grid;
    grid-template-areas: 'name wins' 'members wins';
    grid-template-columns: minmax(0, 1fr) auto;
    padding: 0.5rem 0.75rem;
    text-align: left;
  }
  .name {
    font-size: 1.1em;
    grid-area: name;
    overflow-wrap: break-word;
  }
  .members {
    grid-area: members;
  }
  .wins {
    font-size: 1.3em;
    grid-area: wins;
  }
  .presenting {
    transform: none;
  }
}
</style>
