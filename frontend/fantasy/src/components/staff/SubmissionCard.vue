<template>
  <div class="card-panel mb-2 submission" :class="{ 'up-next': upNext }">
    <div class="d-flex align-items-start gap-2">
      <button class="btn btn-sm vote"
              :class="submission.MyStaffVote ? 'btn-warning' : 'btn-outline-warning'"
              :disabled="!canVote"
              :title="submission.MyStaffVote ? 'Remove your upvote' : 'Upvote'"
              @click="$emit('vote', submission.ID)"
      >
        <FontAwesomeIcon :icon="[submission.MyStaffVote ? 'fas' : 'far', 'star']" />
        <span class="d-block">{{ submission.StaffVotes }}</span>
      </button>
      <div class="flex-grow-1">
        <div class="small text-muted">
          <FontAwesomeIcon :icon="['fas', kindIcon]" />
          {{ kindLabel }}
          <template v-if="submission.AuthorName">&middot; {{ submission.AuthorName }}</template>
          <span v-if="upNext" class="badge text-bg-primary ms-1">Up next</span>
          <span v-if="submission.Status === 'Played'" class="badge text-bg-secondary ms-1">
            Played &middot; {{ submission.Net > 0 ? '+' : '' }}{{ submission.Net }}
          </span>
        </div>
        <div class="text">
          {{ submission.Text }}
        </div>
      </div>
    </div>
    <div v-if="submission.Status !== 'Played'" class="d-flex flex-wrap gap-1 mt-2 justify-content-end">
      <button v-if="submission.Status !== 'Approved'"
              class="btn btn-success btn-sm"
              @click="$emit('status', submission.ID, 'Approved')"
      >
        <FontAwesomeIcon :icon="['fas', 'check']" />
        Approve
      </button>
      <button v-if="submission.Status !== 'Rejected'"
              class="btn btn-outline-danger btn-sm"
              @click="$emit('status', submission.ID, 'Rejected')"
      >
        <FontAwesomeIcon :icon="['fas', 'xmark']" />
        Reject
      </button>
      <button v-if="submission.Status !== 'Pending'"
              class="btn btn-outline-secondary btn-sm"
              @click="$emit('status', submission.ID, 'Pending')"
      >
        Back to queue
      </button>
      <button v-if="canPlay"
              class="btn btn-primary btn-sm"
              @click="$emit('play', submission.ID)"
      >
        <FontAwesomeIcon :icon="['fas', 'play']" />
        Play now
      </button>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { KindIcons } from '@/libs/kinds';
import { SubmissionKindLabels, type Submission } from '@/types/game';

export default defineComponent({
  name: 'SubmissionCard',
  components: {
    FontAwesomeIcon,
  },
  props: {
    submission: {
      type: Object as PropType<Submission>,
      required: true,
    },
    upNext: {
      type: Boolean,
      default: false,
    },
    canPlay: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['vote', 'status', 'play'],
  computed: {
    kindLabel(): string {
      return SubmissionKindLabels[this.submission.Kind] ?? this.submission.Kind;
    },
    kindIcon(): string {
      return KindIcons[this.submission.Kind] ?? 'bullhorn';
    },
    canVote(): boolean {
      return this.submission.Status === 'Pending' || this.submission.Status === 'Approved';
    },
  },
});
</script>

<style scoped lang="scss">
.vote {
  min-width: 3rem;
}
.text {
  overflow-wrap: anywhere;
}
.up-next {
  border-color: var(--bs-primary);
  box-shadow: 0 0 0 2px var(--bs-primary);
}
</style>
