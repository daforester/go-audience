<template>
  <div class="container py-3">
    <div v-if="presenting" class="alert presenting-banner text-center mb-3">
      <FontAwesomeIcon :icon="['fas', 'microphone']" class="dance" />
      Your team is presenting &mdash; over to you!
    </div>

    <div v-if="store.myTeam" class="card-panel mb-3 text-center" :style="{ borderColor: colour }">
      <div class="h3 mb-1" :style="{ color: colour }">
        {{ store.myTeam.Name }}
      </div>
      <div class="small">
        {{ store.myTeam.Members.map((m) => m.Name).join(', ') }}
      </div>
    </div>

    <template v-if="round && round.Topic && gameCode > GameCodes.Intake && gameCode < GameCodes.FinalTally">
      <TopicCard :topic="round.Topic" :round-number="round.Number" class="mb-3" />
      <p v-if="gameCode === GameCodes.TopicReveal" class="text-center h5">
        Huddle with your team: how will the con directors handle this?
      </p>
      <p v-else-if="gameCode === GameCodes.Responses && !presenting" class="text-center h5">
        <template v-if="round.PresentingTeam">
          {{ store.teamName(round.PresentingTeam) }} are presenting their answer.
        </template>
      </p>
      <p v-else-if="gameCode === GameCodes.AudienceVote" class="text-center h5">
        The audience are voting. {{ plural('vote', round.TotalVotes) }} so far...
      </p>
      <p v-else-if="gameCode === GameCodes.RoundResult" class="text-center h4">
        <template v-if="round.WinnerTeamIDs.length === 0">
          No votes, no winner this round.
        </template>
        <template v-else-if="round.WinnerTeamIDs.includes(player.TeamID)">
          <FontAwesomeIcon :icon="['fas', 'trophy']" class="text-warning" />
          Your team won this round!
        </template>
        <template v-else>
          {{ round.WinnerTeamIDs.map((id) => store.teamName(id)).join(' & ') }} won this round.
        </template>
      </p>
    </template>

    <div v-else-if="gameCode === GameCodes.FinalTally && game.Final" class="text-center mb-3">
      <h2>
        <FontAwesomeIcon :icon="['fas', 'crown']" class="text-warning" />
        {{ winnerText }}
      </h2>
    </div>

    <PleaseWait v-else icon="user-tie" message="Waiting for the next topic..." />

    <Scoreboard v-if="game.Teams.length" :teams="game.Teams" :presenting="round?.PresentingTeam || ''" class="mt-3" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import PlayerView from './Player';
import { plural } from '@/libs/plural';
import PleaseWait from '@/components/player/PleaseWait.vue';
import Scoreboard from '@/components/display/Scoreboard.vue';
import TopicCard from '@/components/display/TopicCard.vue';

export default defineComponent({
  name: 'DirectorView',
  components: {
    FontAwesomeIcon,
    PleaseWait,
    Scoreboard,
    TopicCard,
  },
  mixins: [PlayerView],
  methods: {
    plural,
  },
  computed: {
    colour(): string {
      return this.store.teamColour(this.player.TeamID) ?? '';
    },
    presenting(): boolean {
      return (
        this.gameCode === this.GameCodes.Responses &&
        !!this.round &&
        this.round.PresentingTeam === this.player.TeamID
      );
    },
    winnerText(): string {
      const winners = this.game.Final?.WinnerTeamIDs || [];
      if (winners.length === 0) {
        return 'Nobody won... somehow.';
      }
      if (winners.includes(this.player.TeamID)) {
        return winners.length > 1 ? 'Your team shares the crown!' : 'Your team runs the best convention!';
      }
      return `${winners.map((id) => this.store.teamName(id)).join(' & ')} win the convention.`;
    },
  },
});
</script>

<style scoped lang="scss">
@use '@/libs/scss/animations/dance.scss';

.card-panel {
  border-width: 3px;
}
.presenting-banner {
  background: #f2a516;
  font-size: 1.4em;
  font-weight: 800;
}
</style>
