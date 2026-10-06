<template>
  <div class="container py-3">
    <GameControls :game-code="gameCode" />

    <div v-if="round && round.Topic && gameCode > GameCodes.Intake" class="card-panel mb-3">
      <h2 class="h5">
        Current Topic
      </h2>
      <TopicCard :topic="round.Topic" :round-number="round.Number" compact />
      <div class="d-flex justify-content-center gap-4 mt-2">
        <span><FontAwesomeIcon :icon="['fas', 'thumbs-up']" /> {{ round.Topic.Ups }}</span>
        <span><FontAwesomeIcon :icon="['fas', 'thumbs-down']" /> {{ round.Topic.Downs }}</span>
        <span>Net {{ round.Topic.Net }}</span>
        <span>{{ round.TotalVotes }} votes cast</span>
      </div>
      <div v-if="round.Votes" class="mt-2">
        <div v-for="team in roundTeams" :key="team.ID" class="d-flex align-items-center gap-2">
          <span class="team-swatch" :style="{ background: store.teamColour(team.ID) }" />
          <span class="team-label">{{ team.Name }}</span>
          <div class="progress flex-grow-1" role="progressbar" :aria-label="`${team.Name} votes`">
            <div class="progress-bar"
                 :style="{ width: `${votePercent(team.ID)}%`, background: store.teamColour(team.ID) }"
            />
          </div>
          <span>{{ round.Votes[team.ID] || 0 }}</span>
          <FontAwesomeIcon v-if="round.WinnerTeamIDs.includes(team.ID)" :icon="['fas', 'crown']" class="text-warning" />
        </div>
      </div>
    </div>

    <div class="row">
      <div class="col-lg-7">
        <div class="card-panel mb-3">
          <h2 class="h5">
            Submissions
          </h2>
          <SubmissionQueue :game-code="gameCode" allow-play />
        </div>
      </div>
      <div class="col-lg-5">
        <TeamEditor :game-code="gameCode" />
        <PlayerAssignment :game-code="gameCode" />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import PlayerView from './Player';
import GameControls from '@/components/admin/GameControls.vue';
import PlayerAssignment from '@/components/admin/PlayerAssignment.vue';
import TeamEditor from '@/components/admin/TeamEditor.vue';
import TopicCard from '@/components/display/TopicCard.vue';
import SubmissionQueue from '@/components/staff/SubmissionQueue.vue';

export default defineComponent({
  name: 'AdminView',
  components: {
    FontAwesomeIcon,
    GameControls,
    PlayerAssignment,
    SubmissionQueue,
    TeamEditor,
    TopicCard,
  },
  mixins: [PlayerView],
  methods: {
    votePercent(teamID: string): number {
      const total = this.round?.TotalVotes || 0;
      if (!total || !this.round?.Votes) {
        return 0;
      }
      return Math.round(((this.round.Votes[teamID] || 0) / total) * 100);
    },
  },
});
</script>

<style scoped lang="scss">
.team-label {
  min-width: 8rem;
}
</style>
