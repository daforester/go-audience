<template>
  <div v-if="mine.length" class="mt-4">
    <h2 class="h6">
      Your submissions
    </h2>
    <ul class="list-unstyled small">
      <li v-for="s in mine" :key="s.ID" class="d-flex gap-2 mb-1">
        <span class="badge align-self-start" :class="badge(s.Status)">{{ label(s.Status) }}</span>
        <span class="text">{{ s.Text }}</span>
      </li>
    </ul>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import useGameStore from '@/stores/game';
import type { Submission, SubmissionStatus } from '@/types/game';

const Labels: Record<SubmissionStatus, string> = {
  Pending: 'With staff',
  Approved: 'Approved',
  Rejected: 'Not this time',
  Played: 'Played!',
};

const Badges: Record<SubmissionStatus, string> = {
  Pending: 'text-bg-light',
  Approved: 'text-bg-success',
  Rejected: 'text-bg-secondary',
  Played: 'text-bg-warning',
};

export default defineComponent({
  name: 'MySubmissions',
  data() {
    return {
      store: useGameStore(),
    };
  },
  computed: {
    mine(): Submission[] {
      return [...(this.store.Game.MySubmissions || [])].reverse();
    },
  },
  methods: {
    label(status: SubmissionStatus): string {
      return Labels[status];
    },
    badge(status: SubmissionStatus): string {
      return Badges[status];
    },
  },
});
</script>

<style scoped lang="scss">
.text {
  overflow-wrap: anywhere;
}
</style>
