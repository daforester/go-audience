<template>
  <div class="d-flex justify-content-center flex-wrap">
    <div v-if="!game.loading" class="text-center">
      <h1 class="popper">
        {{ label }}
      </h1>
      <h1 class="popper">
        {{ winningOption }}
      </h1>
      <div class="mt-5 toaster toaster-once toaster-delay-3 w-100">
        <h5>Voting Breakdown</h5>
        <div class="m-2">
          <FontAwesomeIcon :icon="['fal', 'eye']" />
          <div class="progress mt-1">
            <div ref="audienceProgressBar"
                 class="progress-bar bg-danger"
                 role="progressbar"
                 :aria-valuenow="audienceVotes"
                 aria-valuemin="0"
                 aria-valuemax="100"
            >
              <!-- -->
            </div>
          </div>
        </div>
        <div class="m-2">
          <FontAwesomeIcon :icon="['fal', 'bolt']" />
          <div class="progress mt-1">
            <div ref="godProgressBar"
                 class="progress-bar bg-success"
                 role="progressbar"
                 :aria-valuenow="godVotes"
                 aria-valuemin="0"
                 aria-valuemax="100"
            >
              <!-- -->
            </div>
          </div>
        </div>
        <div class="m-2">
          <FontAwesomeIcon :icon="['fal', 'swords']" />
          <FontAwesomeIcon :icon="['fal', 'mandolin']" />
          <FontAwesomeIcon :icon="['fal', 'wand-magic']" />
          <div class="progress mt-1">
            <div ref="partyProgressBar"
                 class="progress-bar"
                 role="progressbar"
                 :aria-valuenow="partyVotes"
                 aria-valuemin="0"
                 aria-valuemax="100"
            >
              <!-- -->
            </div>
          </div>
        </div>
        <div>
          Total: {{ votePercent }}%
        </div>
        <h2 class="mt-5">
          Runners Up...
        </h2>
        <div v-for="index in 2" :key="`runnerUp${index}`">
          <template v-if="index < sortedOptions.length">
            {{ sortedOptions[index].Name }}
            {{ percentOfTotal(sortedOptions[index].Total) }}%
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

export default {
  name: 'VoteResult',
  components: {
    FontAwesomeIcon,
  },
  props: {
    category: {
      type: String,
      required: false,
      default() {
        return 'options';
      },
    },
    game: {
      type: Object,
      required: true,
    },
    label: {
      type: String,
      required: false,
      default() {
        return 'Your fate is sealed...';
      },
    },
  },
  computed: {
    audienceVotes() {
      try {
        const index = this.winningOptionIndex();
        const total = this.game.Votes.reduce((a, v) => a + (v.Audience ? v.Audience : 0), 0);
        const votes = this.game.Votes[index].Audience
          ? this.game.Votes[index].Audience : 0;
        return (votes / total) * 100;
      } catch (e) {
        return 0;
      }
    },
    godVotes() {
      try {
        const index = this.winningOptionIndex();
        const total = this.game.Votes.reduce((a, v) => a + (v.God ? v.God : 0), 0);
        const votes = this.game.Votes[index].God
          ? this.game.Votes[index].God : 0;
        return (votes / total) * 100;
      } catch (e) {
        return 0;
      }
    },
    partyVotes() {
      try {
        const index = this.winningOptionIndex();
        const total = this.game.Votes.reduce((a, v) => a
            + (v.Warrior ? v.Warrior : 0)
            + (v.Bard ? v.Bard : 0)
            + (v.Wizard ? v.Wizard : 0),
        0);
        const votes = (this.game.Votes[index].Warrior
          ? this.game.Votes[index].Warrior : 0)
            + (this.game.Votes[index].Bard
              ? this.game.Votes[index].Bard : 0)
            + (this.game.Votes[index].Wizard
              ? this.game.Votes[index].Wizard : 0);
        return (votes / total) * 100;
      } catch (e) {
        return 0;
      }
    },
    sortedOptions() {
      const wi = this.winningOptionIndex();
      return Array.from(this.game.Votes).map((o, i) => {
        o.Name = this.game.Options[i];
        if (i === wi) {
          o.Total += 1;
        }
        return o;
      }).sort((a, b) => b.Total - a.Total);
    },
    totalVotes() {
      return this.game.Votes.reduce((t, v) => t + v.Total, 0);
    },
    votePercent() {
      try {
        const index = this.winningOptionIndex();
        const total = this.totalVotes;
        return Math.round((this.game.Votes[index].Total / total) * 100);
      } catch (e) {
        return 0;
      }
    },
    winningOption() {
      try {
        const index = this.winningOptionIndex();
        return this.game.Options[index];
      } catch (e) {
        return '';
      }
    },
  },
  mounted() {
    if (!this.game.loading) {
      this.resizeBars();
    }
  },
  updated() {
    if (!this.game.loading) {
      this.resizeBars();
    }
  },
  methods: {
    percentOfTotal(x) {
      return Math.round((x / this.totalVotes) * 100);
    },
    resizeBars() {
      if (typeof this.$refs.audienceProgressBar !== 'undefined') {
        this.$refs.audienceProgressBar.style.width = `${this.$refs.audienceProgressBar.getAttribute('aria-valuenow')}%`;
      }
      if (typeof this.$refs.partyProgressBar !== 'undefined') {
        this.$refs.partyProgressBar.style.width = `${this.$refs.partyProgressBar.getAttribute('aria-valuenow')}%`;
      }
      if (typeof this.$refs.godProgressBar !== 'undefined') {
        this.$refs.godProgressBar.style.width = `${this.$refs.godProgressBar.getAttribute('aria-valuenow')}%`;
      }
    },
    winningOptionIndex() {
      if (this.game.Options.length === 0 || this.game.Votes.length === 0) {
        throw new Error('options or votes are zero length');
      }
      const wv = Math.max(...this.game.Votes.map((x) => x.Total));
      let wi = 0;
      this.game.Votes.forEach((v, i) => {
        if (v.Total === wv) {
          wi = i;
        }
      });
      return wi;
    },
  },
};

</script>

<style scoped>
@import "../../libs/scss/animations/dance.scss";
@import "../../libs/scss/animations/popper.scss";
@import "../../libs/scss/animations/toaster.scss";

.text-center {
  font-size: 1.5rem;
}

</style>
