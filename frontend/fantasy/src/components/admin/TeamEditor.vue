<template>
  <div class="card-panel mb-3">
    <h2 class="h5">
      <FontAwesomeIcon :icon="['fas', 'users']" />
      Director Teams
      <small class="text-muted">{{ game.Teams.length }} / {{ maxTeams }}</small>
    </h2>
    <ul class="list-unstyled mb-2">
      <li v-for="team in game.Teams" :key="team.ID" class="team-row d-flex flex-wrap align-items-center gap-2 mb-2">
        <span class="team-swatch" :style="{ background: store.teamColour(team.ID) }" />
        <input class="form-control form-control-sm team-name"
               :value="team.Name"
               :aria-label="`Name of ${team.Name}`"
               maxlength="32"
               @change="renameTeam(team.ID, ($event.target as HTMLInputElement).value)"
        >
        <span class="small text-nowrap">{{ plural('win', team.RoundWins) }}</span>
        <button class="btn btn-outline-danger btn-sm"
                :disabled="roundInProgress"
                :title="`Remove ${team.Name}`"
                @click="removeTeam(team.ID, team.Name)"
        >
          <FontAwesomeIcon :icon="['fas', 'trash']" />
        </button>
        <span class="members small text-muted">
          {{ team.Members.map((m) => m.Name).join(', ') || 'No directors yet' }}
        </span>
      </li>
    </ul>
    <form v-if="game.Teams.length < maxTeams" class="d-flex gap-2" @submit.prevent="addTeam">
      <input v-model="NewTeamName"
             class="form-control form-control-sm team-name"
             placeholder="New team name"
             aria-label="New team name"
             maxlength="32"
             required
      >
      <button class="btn btn-primary btn-sm" :disabled="roundInProgress">
        <FontAwesomeIcon :icon="['fas', 'plus']" />
        Add Team
      </button>
    </form>
    <div v-if="ErrorMessage" class="alert alert-danger mt-2 mb-0 py-1">
      {{ ErrorMessage }}
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import PlayerView from '@/views/Player';
import { plural } from '@/libs/plural';
import { GameCodes } from '@/types/game';

export default defineComponent({
  name: 'TeamEditor',
  components: {
    FontAwesomeIcon,
  },
  mixins: [PlayerView],
  data() {
    return {
      NewTeamName: '',
    };
  },
  computed: {
    maxTeams(): number {
      return this.game.MaxTeams || 3;
    },
    roundInProgress(): boolean {
      return [GameCodes.TopicReveal, GameCodes.Responses, GameCodes.AudienceVote].includes(
        this.gameCode as 100 | 110 | 120,
      );
    },
  },
  methods: {
    plural,
    addTeam() {
      this.gameRequest('post', '/game/team/add', { Name: this.NewTeamName })
        .then(() => {
          this.NewTeamName = '';
        })
        .catch(() => {});
    },
    renameTeam(id: string, name: string) {
      this.gameRequest('post', `/game/team/change/${id}`, { Name: name }).catch(() => {});
    },
    removeTeam(id: string, name: string) {
      if (window.confirm(`Remove ${name}? Its directors go back to the audience.`)) {
        this.gameRequest('post', `/game/team/remove/${id}`).catch(() => {});
      }
    },
  },
});
</script>

<style scoped lang="scss">
.team-row .team-name {
  flex: 1 1 10rem;
  min-width: 0;
}
form .team-name {
  max-width: 16rem;
}
.members {
  flex-basis: 100%;
  padding-left: 1.4rem;
}
</style>
