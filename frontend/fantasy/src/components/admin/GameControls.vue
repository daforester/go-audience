<template>
  <div class="card-panel mb-3">
    <div class="d-flex flex-wrap gap-2 align-items-center mb-2">
      <h2 class="h5 m-0 me-auto">
        {{ phaseLabel }}
        <small v-if="round && gameCode > GameCodes.Intake" class="text-muted">
          &middot; Round {{ round.Number }}
        </small>
      </h2>
      <button class="btn btn-outline-danger btn-sm" @click="resetGame">
        <FontAwesomeIcon :icon="['fas', 'rotate-left']" />
        Reset Game
      </button>
    </div>

    <div class="d-flex flex-wrap gap-2 mb-2">
      <button v-if="gameCode === GameCodes.Lobby" class="btn btn-primary" @click="setGameCode(GameCodes.Intake)">
        Open the Convention
      </button>
      <button v-if="canStartRound"
              class="btn btn-primary"
              :disabled="!!readiness || !game.NextTopicID"
              @click="startRound()"
      >
        <FontAwesomeIcon :icon="['fas', 'bullhorn']" />
        Reveal Next Topic
      </button>
      <button v-if="gameCode === GameCodes.TopicReveal" class="btn btn-primary" @click="setGameCode(GameCodes.Responses)">
        <FontAwesomeIcon :icon="['fas', 'microphone']" />
        Directors Respond
      </button>
      <button v-if="gameCode === GameCodes.Responses" class="btn btn-primary" @click="setGameCode(GameCodes.AudienceVote)">
        Open Audience Vote
      </button>
      <button v-if="gameCode === GameCodes.AudienceVote" class="btn btn-primary" @click="setGameCode(GameCodes.RoundResult)">
        <FontAwesomeIcon :icon="['fas', 'trophy']" />
        Reveal Round Winner
      </button>
      <button v-if="gameCode === GameCodes.RoundResult" class="btn btn-outline-primary" @click="setGameCode(GameCodes.Intake)">
        Back to Intake
      </button>
      <button v-if="gameCode === GameCodes.Intake || gameCode === GameCodes.RoundResult"
              class="btn btn-warning"
              @click="setGameCode(GameCodes.FinalTally)"
      >
        <FontAwesomeIcon :icon="['fas', 'crown']" />
        Final Tally
      </button>
    </div>

    <div v-if="gameCode === GameCodes.Responses && round" class="d-flex flex-wrap gap-2 mb-2 align-items-center">
      <span>Presenting:</span>
      <button v-for="team in roundTeams"
              :key="team.ID"
              class="btn btn-sm"
              :class="round.PresentingTeam === team.ID ? 'btn-dark' : 'btn-outline-dark'"
              @click="setPresenting(team.ID)"
      >
        <span class="team-swatch" :style="{ background: store.teamColour(team.ID) }" />
        {{ team.Name }}
      </button>
    </div>

    <p v-if="canStartRound && readiness" class="text-danger small m-0">
      {{ readiness }}
    </p>
    <p v-else-if="canStartRound && !game.NextTopicID" class="text-muted small m-0">
      Approve a submission (or pick one below) to reveal a topic.
    </p>

    <div class="d-flex flex-wrap gap-3 mt-2">
      <div class="form-check form-switch">
        <input id="submissionsEnabled"
               class="form-check-input"
               type="checkbox"
               :checked="game.SubmissionsEnabled"
               @change="toggle('/game/submissions', ($event.target as HTMLInputElement).checked)"
        >
        <label class="form-check-label" for="submissionsEnabled">Audience submissions open</label>
      </div>
      <div class="form-check form-switch">
        <input id="showSubmitterNames"
               class="form-check-input"
               type="checkbox"
               :checked="game.ShowSubmitterNames"
               @change="toggle('/game/submitter-names', ($event.target as HTMLInputElement).checked)"
        >
        <label class="form-check-label" for="showSubmitterNames">Show who submitted</label>
      </div>
    </div>

    <div v-if="ErrorMessage" class="alert alert-danger mt-2 mb-0 py-1">
      {{ ErrorMessage }}
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import PlayerView from '@/views/Player';
import { GameCodes } from '@/types/game';

export const PhaseLabels: Record<number, string> = {
  [GameCodes.Lobby]: 'Lobby',
  [GameCodes.Intake]: 'Intake: collecting submissions',
  [GameCodes.TopicReveal]: 'Topic revealed',
  [GameCodes.Responses]: 'Directors responding',
  [GameCodes.AudienceVote]: 'Audience voting',
  [GameCodes.RoundResult]: 'Round result',
  [GameCodes.FinalTally]: 'Final tally',
};

export default defineComponent({
  name: 'GameControls',
  components: {
    FontAwesomeIcon,
  },
  mixins: [PlayerView],
  computed: {
    phaseLabel(): string {
      return PhaseLabels[this.gameCode] ?? '';
    },
    canStartRound(): boolean {
      return this.gameCode === GameCodes.Intake || this.gameCode === GameCodes.RoundResult;
    },
    readiness(): string {
      if (this.game.Teams.length < 2) {
        return 'Add at least two director teams.';
      }
      const empty = this.game.Teams.find((t) => t.Members.length === 0);
      if (empty) {
        return `${empty.Name} needs at least one con director.`;
      }
      return '';
    },
  },
  methods: {
    setGameCode(code: number) {
      this.gameRequest('post', '/game/code', { Code: code }).catch(() => {});
    },
    startRound(submissionID = '') {
      this.gameRequest('post', '/game/round/start', { SubmissionID: submissionID }).catch(() => {});
    },
    setPresenting(teamID: string) {
      this.gameRequest('post', '/game/round/presenting', { TeamID: teamID }).catch(() => {});
    },
    toggle(path: string, enabled: boolean) {
      this.gameRequest('post', path, { Enabled: enabled }).catch(() => {});
    },
    resetGame() {
      if (window.confirm('Reset the game to the lobby? Scores and rounds are cleared.')) {
        this.setGameCode(GameCodes.Lobby);
      }
    },
  },
});
</script>
