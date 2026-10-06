<template>
  <div class="topic text-center" :class="{ compact }">
    <div class="kind text-uppercase">
      <FontAwesomeIcon :icon="['fas', kindIcon]" />
      {{ kindLabel }}
      <template v-if="roundNumber">&middot; Round {{ roundNumber }}</template>
    </div>
    <blockquote class="text mb-2">
      &ldquo;{{ topic.Text }}&rdquo;
    </blockquote>
    <div v-if="topic.AuthorName" class="author">
      &mdash; {{ topic.AuthorName }}
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { KindIcons } from '@/libs/kinds';
import { SubmissionKindLabels, type Submission } from '@/types/game';

export default defineComponent({
  name: 'TopicCard',
  components: {
    FontAwesomeIcon,
  },
  props: {
    topic: {
      type: Object as PropType<Submission>,
      required: true,
    },
    roundNumber: {
      type: Number,
      default: 0,
    },
    compact: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    kindLabel(): string {
      return SubmissionKindLabels[this.topic.Kind] ?? this.topic.Kind;
    },
    kindIcon(): string {
      return KindIcons[this.topic.Kind] ?? 'bullhorn';
    },
  },
});
</script>

<style scoped lang="scss">
.topic {
  background: rgba(255, 255, 255, 0.95);
  border: 3px solid #2b2233;
  border-radius: 1.2rem;
  margin: 0 auto;
  max-width: min(90vw, 70rem);
  padding: 1.5rem 2rem;
}
.kind {
  color: #7b3fbf;
  font-weight: 700;
  letter-spacing: 0.1em;
}
.text {
  font-size: calc(1.2rem + 2.4vw);
  font-weight: 700;
  overflow-wrap: anywhere;
}
.author {
  font-size: calc(1rem + 0.8vw);
  font-style: italic;
}
.compact {
  padding: 0.75rem 1.25rem;
  .text {
    font-size: calc(1rem + 1.2vw);
  }
}
</style>
