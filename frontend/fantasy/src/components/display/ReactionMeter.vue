<template>
  <div class="meter mx-auto" role="meter" :aria-valuenow="topic.Net" :aria-label="`Audience reaction ${topic.Net}`">
    <div class="side up">
      <FontAwesomeIcon :icon="['fas', 'thumbs-up']" />
      {{ topic.Ups }}
    </div>
    <div class="track">
      <div class="fill" :style="{ width: `${upShare}%` }" />
    </div>
    <div class="side down">
      {{ topic.Downs }}
      <FontAwesomeIcon :icon="['fas', 'thumbs-down']" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import type { Submission } from '@/types/game';

export default defineComponent({
  name: 'ReactionMeter',
  components: {
    FontAwesomeIcon,
  },
  props: {
    topic: {
      type: Object as PropType<Submission>,
      required: true,
    },
  },
  computed: {
    upShare(): number {
      const total = this.topic.Ups + this.topic.Downs;
      if (total === 0) {
        return 50;
      }
      return Math.round((this.topic.Ups / total) * 100);
    },
  },
});
</script>

<style scoped lang="scss">
.meter {
  align-items: center;
  display: flex;
  font-size: calc(1rem + 1.5vw);
  font-weight: 800;
  gap: 1rem;
  max-width: min(90vw, 50rem);
}
.track {
  background: #d63c5e;
  border: 2px solid #2b2233;
  border-radius: 1rem;
  flex-grow: 1;
  height: 2.5rem;
  overflow: hidden;
}
.fill {
  background: #3aa45a;
  height: 100%;
  transition: width 0.3s ease-out;
}
.up {
  color: #3aa45a;
}
.down {
  color: #d63c5e;
}
</style>
