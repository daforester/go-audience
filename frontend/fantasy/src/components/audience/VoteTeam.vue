<template>
  <div v-if="round">
    <h1 class="h3 text-center">
      Which team handled it best?
    </h1>
    <p v-if="round.Topic" class="text-center small text-muted text">
      &ldquo;{{ round.Topic.Text }}&rdquo;
    </p>
    <div class="d-grid gap-3">
      <button v-for="team in roundTeams"
              :key="team.ID"
              class="btn big-button team-button"
              :class="{ chosen: round.MyVote === team.ID }"
              :style="buttonStyle(team.ID)"
              :disabled="requesting"
              :aria-pressed="round.MyVote === team.ID"
              @click="vote(team.ID)"
      >
        <FontAwesomeIcon v-if="round.MyVote === team.ID" :icon="['fas', 'check']" />
        {{ team.Name }}
      </button>
    </div>
    <p class="text-center small mt-2">
      <template v-if="round.MyVote">You can change your vote until voting closes.</template>
      {{ round.TotalVotes }} votes so far
    </p>
    <div v-if="ErrorMessage" class="form-text text-danger text-center">
      {{ ErrorMessage }}
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import PlayerView from '@/views/Player';

export default defineComponent({
  name: 'VoteTeam',
  components: {
    FontAwesomeIcon,
  },
  mixins: [PlayerView],
  methods: {
    buttonStyle(teamID: string) {
      const colour = this.store.teamColour(teamID);
      const chosen = this.round?.MyVote === teamID;
      return {
        borderColor: colour,
        background: chosen ? colour : 'white',
        color: chosen ? 'white' : colour,
      };
    },
    vote(teamID: string) {
      this.gameRequest('post', '/game/round/vote', { TeamID: teamID })
        .then(() => {
          this.store.setMyVote(teamID);
        })
        .catch(() => {});
    },
  },
});
</script>

<style scoped lang="scss">
.team-button {
  border-width: 3px;
  font-weight: 800;
}
.text {
  overflow-wrap: anywhere;
}
</style>
