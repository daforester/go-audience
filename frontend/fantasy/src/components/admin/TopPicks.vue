<template>
  <div class="card-panel mb-3">
    <h2 class="h5">
      <FontAwesomeIcon :icon="['fas', 'star']" class="text-warning" />
      Top Picks
      <small class="text-muted">the staff's favourites, in the order they would play</small>
    </h2>

    <div v-if="pinned" class="alert alert-dark py-2 d-flex align-items-center gap-2">
      <FontAwesomeIcon :icon="['fas', 'thumbtack']" />
      <span class="flex-grow-1 text">
        <strong>Your pick plays next:</strong> {{ pinned.Text }}
      </span>
      <button class="btn btn-sm btn-light" @click="pin('')">
        Unpin
      </button>
    </div>

    <p v-if="picks.length === 0" class="text-muted small m-0">
      No approved submissions yet. Staff approve and star submissions; the best appear here.
    </p>
    <ol class="list-unstyled m-0">
      <li v-for="(s, i) in picks" :key="s.ID" class="d-flex align-items-start gap-2 mb-2">
        <span class="rank">{{ i + 1 }}</span>
        <span class="stars text-warning" :title="`${s.StaffVotes} staff stars`">
          <FontAwesomeIcon :icon="['fas', 'star']" /> {{ s.StaffVotes }}
        </span>
        <span class="flex-grow-1 text">
          <span v-if="s.ID === game.NextTopicID" class="badge text-bg-primary me-1">Up next</span>
          {{ s.Text }}
          <small v-if="s.AuthorName" class="text-muted">&middot; {{ s.AuthorName }}</small>
        </span>
        <span class="d-flex flex-wrap gap-1 justify-content-end">
          <button v-if="canStartRound" class="btn btn-primary btn-sm" @click="play(s.ID)">
            <FontAwesomeIcon :icon="['fas', 'play']" />
            Play now
          </button>
          <button v-if="s.ID !== game.PinnedTopicID" class="btn btn-outline-dark btn-sm" @click="pin(s.ID)">
            <FontAwesomeIcon :icon="['fas', 'thumbtack']" />
            Play next
          </button>
          <button class="btn btn-outline-danger btn-sm" @click="veto(s.ID)">
            <FontAwesomeIcon :icon="['fas', 'xmark']" />
            Veto
          </button>
        </span>
      </li>
    </ol>

    <div v-if="ErrorMessage" class="alert alert-danger mt-2 mb-0 py-1">
      {{ ErrorMessage }}
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import PlayerView from '@/views/Player';
import { GameCodes, type Submission } from '@/types/game';

// TopPicks shows the admin the staff's highest starred submissions, with the admin's
// override (play now, or pin to play next) and veto
export default defineComponent({
  name: 'TopPicks',
  components: {
    FontAwesomeIcon,
  },
  mixins: [PlayerView],
  computed: {
    picks(): Submission[] {
      const byID = new Map((this.game.Submissions || []).map((s) => [s.ID, s]));
      return (this.game.TopPickIDs || []).map((id) => byID.get(id)).filter((s): s is Submission => !!s);
    },
    pinned(): Submission | undefined {
      return (this.game.Submissions || []).find((s) => s.ID === this.game.PinnedTopicID);
    },
    canStartRound(): boolean {
      return this.gameCode === GameCodes.Intake || this.gameCode === GameCodes.RoundResult;
    },
  },
  methods: {
    play(id: string) {
      this.gameRequest('post', '/game/round/start', { SubmissionID: id }).catch(() => {});
    },
    pin(id: string) {
      this.gameRequest('post', '/game/submission/pin', { SubmissionID: id }).catch(() => {});
    },
    veto(id: string) {
      this.gameRequest('post', `/game/submission/status/${id}`, { Status: 'Rejected' }).catch(() => {});
    },
  },
});
</script>

<style scoped lang="scss">
.rank {
  font-weight: 800;
  min-width: 1.2rem;
}
.stars {
  min-width: 2.5rem;
  white-space: nowrap;
}
.text {
  overflow-wrap: anywhere;
}
</style>
