<template>
  <div class="container py-3">
    <div v-if="round && round.Topic && gameCode > GameCodes.Intake && gameCode < GameCodes.FinalTally"
         class="mb-3"
    >
      <TopicCard :topic="round.Topic" :round-number="round.Number" compact />
    </div>
    <div class="card-panel">
      <h2 class="h5">
        <FontAwesomeIcon :icon="['fas', 'id-badge']" />
        Staff Triage
      </h2>
      <p class="small text-muted">
        Approve the submissions the con directors should face, reject anything unsuitable, and
        star your favourites. The approved submission with the most stars is played next.
      </p>
      <SubmissionQueue :game-code="gameCode" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import PlayerView from './Player';
import TopicCard from '@/components/display/TopicCard.vue';
import SubmissionQueue from '@/components/staff/SubmissionQueue.vue';

export default defineComponent({
  name: 'StaffView',
  components: {
    FontAwesomeIcon,
    SubmissionQueue,
    TopicCard,
  },
  mixins: [PlayerView],
});
</script>
