<template>
  <div>
    <Sunburst />
    <div class="container audience py-3">
      <ReactTopic v-if="gameCode === GameCodes.Responses" :game-code="gameCode" />
      <VoteTeam v-else-if="gameCode === GameCodes.AudienceVote" :game-code="gameCode" />
      <template v-else>
        <div v-if="gameCode === GameCodes.TopicReveal && round && round.Topic" class="mb-3">
          <TopicCard :topic="round.Topic" :round-number="round.Number" compact />
          <p class="text-center mt-2">
            The con directors are huddling... get ready to react!
          </p>
        </div>
        <div v-else-if="gameCode === GameCodes.RoundResult && round" class="card-panel text-center mb-3">
          <FontAwesomeIcon :icon="['fas', 'trophy']" class="text-warning" />
          <template v-if="round.WinnerTeamIDs.length">
            {{ round.WinnerTeamIDs.map((id) => store.teamName(id)).join(' & ') }} won the round!
          </template>
          <template v-else>
            No winner this round.
          </template>
        </div>
        <div v-else-if="gameCode === GameCodes.FinalTally && game.Final" class="card-panel text-center mb-3">
          <FontAwesomeIcon :icon="['fas', 'crown']" class="text-warning" />
          Thanks for coming to the Fantasy Furry Convention!
        </div>

        <div v-if="game.SubmissionsEnabled" class="card-panel">
          <NewSubmission :game-code="gameCode" />
          <MySubmissions />
        </div>
        <PleaseWait v-else-if="gameCode < GameCodes.TopicReveal || gameCode === GameCodes.FinalTally"
                    :name="player.Name"
        >
          <MySubmissions />
        </PleaseWait>
      </template>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import PlayerView from './Player';
import MySubmissions from '@/components/audience/MySubmissions.vue';
import NewSubmission from '@/components/audience/NewSubmission.vue';
import ReactTopic from '@/components/audience/ReactTopic.vue';
import VoteTeam from '@/components/audience/VoteTeam.vue';
import Sunburst from '@/components/backgrounds/Sunburst.vue';
import TopicCard from '@/components/display/TopicCard.vue';
import PleaseWait from '@/components/player/PleaseWait.vue';

export default defineComponent({
  name: 'AudienceView',
  components: {
    FontAwesomeIcon,
    MySubmissions,
    NewSubmission,
    PleaseWait,
    ReactTopic,
    Sunburst,
    TopicCard,
    VoteTeam,
  },
  mixins: [PlayerView],
});
</script>

<style scoped lang="scss">
.audience {
  max-width: 34rem;
}
</style>
