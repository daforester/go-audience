<template>
  <div>
    <ul class="nav nav-pills mb-2 flex-nowrap overflow-auto">
      <li v-for="tab in tabs" :key="tab.status" class="nav-item">
        <button class="nav-link"
                :class="{ active: Tab === tab.status }"
                @click="Tab = tab.status"
        >
          {{ tab.label }}
          <span class="badge text-bg-light">{{ counts[tab.status] }}</span>
        </button>
      </li>
    </ul>
    <div v-if="ErrorMessage" class="alert alert-danger py-1">
      {{ ErrorMessage }}
    </div>
    <p v-if="shown.length === 0" class="text-muted">
      {{ emptyMessage }}
    </p>
    <SubmissionCard v-for="s in shown"
                    :key="s.ID"
                    :submission="s"
                    :up-next="s.ID === game.NextTopicID"
                    :can-play="allowPlay && canStartRound"
                    @vote="vote"
                    @status="setStatus"
                    @play="play"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import PlayerView from '@/views/Player';
import SubmissionCard from './SubmissionCard.vue';
import { GameCodes, type Submission, type SubmissionStatus } from '@/types/game';

export default defineComponent({
  name: 'SubmissionQueue',
  components: {
    SubmissionCard,
  },
  mixins: [PlayerView],
  props: {
    allowPlay: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      Tab: 'Pending' as SubmissionStatus,
      tabs: [
        { status: 'Pending' as SubmissionStatus, label: 'To review' },
        { status: 'Approved' as SubmissionStatus, label: 'Approved' },
        { status: 'Rejected' as SubmissionStatus, label: 'Rejected' },
        { status: 'Played' as SubmissionStatus, label: 'Played' },
      ],
    };
  },
  computed: {
    submissions(): Submission[] {
      return this.game.Submissions || [];
    },
    counts(): Record<SubmissionStatus, number> {
      const counts = { Pending: 0, Approved: 0, Rejected: 0, Played: 0 };
      this.submissions.forEach((s) => {
        counts[s.Status]++;
      });
      return counts;
    },
    shown(): Submission[] {
      const list = this.submissions.filter((s) => s.Status === this.Tab);
      if (this.Tab === 'Approved') {
        // Same order the next topic is chosen in: most staff votes, then oldest
        return list.sort((a, b) => b.StaffVotes - a.StaffVotes || a.Created.localeCompare(b.Created));
      }
      if (this.Tab === 'Pending') {
        return list.sort((a, b) => a.Created.localeCompare(b.Created));
      }
      return list.reverse();
    },
    emptyMessage(): string {
      if (this.Tab === 'Pending') {
        return this.game.SubmissionsEnabled
          ? 'Nothing to review yet. Submissions will appear here as the audience sends them.'
          : 'Nothing to review. Audience submissions are closed.';
      }
      return 'Nothing here.';
    },
    canStartRound(): boolean {
      return this.gameCode === GameCodes.Intake || this.gameCode === GameCodes.RoundResult;
    },
  },
  methods: {
    vote(id: string) {
      this.gameRequest('post', `/game/submission/vote/${id}`).catch(() => {});
    },
    setStatus(id: string, status: SubmissionStatus) {
      this.gameRequest('post', `/game/submission/status/${id}`, { Status: status }).catch(() => {});
    },
    play(id: string) {
      this.gameRequest('post', '/game/round/start', { SubmissionID: id }).catch(() => {});
    },
  },
});
</script>
