<template>
  <div class="d-flex justify-content-center flex-wrap">
    <div v-for="(option, i) in game.Options"
         :key="i"
         ref="items"
         class="m-4 text-center toaster"
    >
      <h1 class="bg-white border border-dark p-2 rounded-lg target">
        {{ option }}
      </h1>
      <div class="progress d-none">
        <div ref="progressbar"
             class="progress-bar"
             role="progressbar"
             :aria-valuenow="votePercent(i)"
             aria-valuemin="0"
             aria-valuemax="100"
        >
          <!-- -->
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LiveVoting',
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
  },
  data() {
    return {
      totalVotes: 0,
    };
  },
  created() {
    if (this.game.Votes) {
      this.totalVotes = this.game.Votes.reduce((a, b) => a + (b.Total ? b.Total : 0), 0);
      this.emitter.on("voteAdded", this.voteAdded);
    }
  },
  unmounted() {
    this.emitter.off("voteAdded", this.voteAdded);
  },
  mounted() {
    this.setWidths();
  },
  methods: {
    setWidths() {
      if (this.$refs.progressbar) {
        this.$refs.progressbar.forEach((e, i) => {
          e.style.width = `${this.votePercent(i)}%`;
        });
      }
    },
    voteAdded(data) {
      this.totalVotes += 1;
      const el = this.$refs.items[data.index];
      if (el) {
        if (el.classList.contains('flasher')) {
          el.classList.remove('flasher');
        }
        el.classList.add('flasher');
        setTimeout(() => {
          el.classList.remove('flasher');
        }, 500);
      }
      this.setWidths();
    },
    votePercent(index) {
      const categoryTotal = (this.game.Votes[index].Total ? this.game.Votes[index].Total : 0);
      return (categoryTotal / this.totalVotes) * 100;
    },
  },
};

</script>

<style scoped lang="scss">
@import '../../libs/scss/animations/dance.scss';
@import "../../libs/scss/animations/toaster.scss";
.toaster {
  animation-iteration-count: 1;
}
.text-center {
  font-size: 1.5rem;
}

@keyframes flashBorder {
  from {
    box-shadow: 0 0 10px black;
    transform: scale(1.2);
  }
  to {
    box-shadow: unset;
    transform: scale(1);
  }
}

.flasher .target {
  animation-name: flashBorder;
  animation-duration:0.5s;
  animation-iteration-count: 1;
  animation-timing-function: linear;
}
</style>
