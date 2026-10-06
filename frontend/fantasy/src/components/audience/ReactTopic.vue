<template>
  <div v-if="round && round.Topic">
    <TopicCard :topic="round.Topic" :round-number="round.Number" compact class="mb-3" />
    <p class="text-center">
      <template v-if="round.PresentingTeam">
        <span class="team-swatch" :style="{ background: store.teamColour(round.PresentingTeam) }" />
        {{ store.teamName(round.PresentingTeam) }} are presenting.
      </template>
      How good a problem is this?
    </p>
    <div class="d-flex gap-3">
      <button class="btn big-button"
              :class="reaction === 1 ? 'btn-success' : 'btn-outline-success'"
              :disabled="requesting"
              aria-label="Upvote this problem"
              :aria-pressed="reaction === 1"
              @click="react(1)"
      >
        <FontAwesomeIcon :icon="[reaction === 1 ? 'fas' : 'far', 'thumbs-up']" />
      </button>
      <button class="btn big-button"
              :class="reaction === -1 ? 'btn-danger' : 'btn-outline-danger'"
              :disabled="requesting"
              aria-label="Downvote this problem"
              :aria-pressed="reaction === -1"
              @click="react(-1)"
      >
        <FontAwesomeIcon :icon="[reaction === -1 ? 'fas' : 'far', 'thumbs-down']" />
      </button>
    </div>
    <div v-if="ErrorMessage" class="form-text text-danger text-center">
      {{ ErrorMessage }}
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import PlayerView from '@/views/Player';
import TopicCard from '@/components/display/TopicCard.vue';

export default defineComponent({
  name: 'ReactTopic',
  components: {
    FontAwesomeIcon,
    TopicCard,
  },
  mixins: [PlayerView],
  computed: {
    reaction(): number {
      return this.round?.MyReaction || 0;
    },
  },
  methods: {
    react(value: number) {
      // Pressing your current reaction again withdraws it
      const reaction = this.reaction === value ? 0 : value;
      this.gameRequest('post', '/game/round/react', { Reaction: reaction })
        .then(() => {
          this.store.setMyReaction(reaction);
        })
        .catch(() => {});
    },
  },
});
</script>
