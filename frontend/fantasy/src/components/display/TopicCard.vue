<template>
  <div class="topic text-center" :class="[{ compact, collapsed: isCollapsed }, lengthClass]">
    <div class="kind text-uppercase">
      <FontAwesomeIcon :icon="['fas', kindIcon]" />
      {{ kindLabel }}
      <template v-if="roundNumber">&middot; Round {{ roundNumber }}</template>
    </div>
    <blockquote class="text mb-2">
      &ldquo;{{ topic.Text }}&rdquo;
    </blockquote>
    <button v-if="collapsible && isLong"
            class="btn btn-link btn-sm p-0 mb-1"
            :aria-expanded="!isCollapsed"
            @click="Expanded = !Expanded"
    >
      {{ isCollapsed ? 'Show all' : 'Show less' }}
    </button>
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

// Text sizes step down with length so a full 280 character topic still fits a 720p projector
const MEDIUM_LENGTH = 90;
const LONG_LENGTH = 170;

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
    // Long topics start clamped to a few lines, for screens where the topic is context
    collapsible: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      Expanded: false,
    };
  },
  computed: {
    kindLabel(): string {
      return SubmissionKindLabels[this.topic.Kind] ?? this.topic.Kind;
    },
    kindIcon(): string {
      return KindIcons[this.topic.Kind] ?? 'bullhorn';
    },
    isLong(): boolean {
      return this.topic.Text.length > MEDIUM_LENGTH;
    },
    isCollapsed(): boolean {
      return this.collapsible && this.isLong && !this.Expanded;
    },
    lengthClass(): string {
      if (this.topic.Text.length > LONG_LENGTH) {
        return 'long';
      }
      if (this.topic.Text.length > MEDIUM_LENGTH) {
        return 'medium';
      }
      return '';
    },
  },
  watch: {
    'topic.ID'() {
      this.Expanded = false;
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
  line-height: 1.25;
  overflow-wrap: anywhere;
}
.medium .text {
  font-size: calc(1.1rem + 1.8vw);
}
.long .text {
  font-size: calc(1rem + 1.35vw);
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
  &.long .text {
    font-size: calc(0.95rem + 0.8vw);
  }
}
.collapsed .text {
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
  display: -webkit-box;
  line-clamp: 3;
  overflow: hidden;
}

// Phones: the topic is read at arm's length, not across a room
@media (max-width: 576px) {
  .topic,
  .compact {
    border-width: 2px;
    max-width: 100%;
    padding: 0.75rem 1rem;
  }
  .text,
  .medium .text,
  .long .text,
  .compact .text,
  .compact.long .text {
    font-size: 1.15rem;
  }
  .kind {
    font-size: 0.8rem;
  }
  .author {
    font-size: 0.95rem;
  }
}
</style>
