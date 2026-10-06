<template>
  <div>
    <Sunburst />
    <div class="d-flex vw-100 main-display position-absolute overflow-hidden">
      <div class="m-auto w-100 px-4">
        <Lobby v-if="gameCode === GameCodes.Lobby" :game="game" />

        <template v-else-if="gameCode === GameCodes.Intake">
          <DisplayMessage effect="toaster" message="Complaints Desk Open" />
          <h2 class="text-center mb-2">
            Join at <strong>{{ joinHost }}</strong> with code <strong>{{ game.Key }}</strong>
            and tell the con directors what's wrong
          </h2>
          <p v-if="game.SubmissionCount" class="text-center h3 mb-4">
            <FontAwesomeIcon :icon="['fas', 'bullhorn']" />
            {{ plural('submission', game.SubmissionCount) }} received
          </p>
          <Scoreboard v-if="game.Teams.length" :teams="game.Teams" :show-wins="!!round" large />
        </template>

        <template v-else-if="gameCode === GameCodes.TopicReveal && round && round.Topic">
          <h2 class="text-center mb-3 popper">
            The audience demands an answer...
          </h2>
          <TopicCard :topic="round.Topic" :round-number="round.Number" class="popper" />
        </template>

        <template v-else-if="gameCode === GameCodes.Responses && round && round.Topic">
          <TopicCard :topic="round.Topic" :round-number="round.Number" compact class="mb-4" />
          <div v-if="round.PresentingTeam" class="presenting text-center mb-4">
            <FontAwesomeIcon :icon="['fas', 'microphone']" class="dance me-3" />
            <span :style="{ color: store.teamColour(round.PresentingTeam) }">
              {{ store.teamName(round.PresentingTeam) }}
            </span>
            respond
          </div>
          <ReactionMeter :topic="round.Topic" />
          <p class="text-center h5 mt-2">
            Is this a good problem? Thumbs up or down on your phone
          </p>
        </template>

        <template v-else-if="gameCode === GameCodes.AudienceVote && round">
          <h1 class="text-center mb-3">
            Vote for the best response!
          </h1>
          <TeamBars :teams="roundTeams" :votes="round.Votes || {}" />
          <p class="text-center h4 mt-3">
            {{ plural('vote', round.TotalVotes) }}
          </p>
        </template>

        <template v-else-if="gameCode === GameCodes.RoundResult && round">
          <h1 class="text-center mb-3 popper">
            <template v-if="round.WinnerTeamIDs.length">
              <FontAwesomeIcon :icon="['fas', 'trophy']" class="trophy" />
              {{ round.WinnerTeamIDs.map((id) => store.teamName(id)).join(' & ') }}
              {{ round.WinnerTeamIDs.length > 1 ? 'share' : 'win' }} round {{ round.Number }}!
            </template>
            <template v-else>
              Nobody voted. Nobody wins.
            </template>
          </h1>
          <TeamBars :teams="roundTeams"
                    :votes="round.Votes || {}"
                    :winners="round.WinnerTeamIDs"
                    height="24vh"
                    class="mb-4"
          />
          <Scoreboard :teams="game.Teams" :show-members="false" large />
        </template>

        <GameEnd v-else-if="gameCode === GameCodes.FinalTally" :key="`end-${game.Key}`" :game="game" />
      </div>
    </div>
    <div v-if="musicCopyrightNotice" class="musicCopyright">
      {{ musicCopyrightNotice }}
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import PlayerView from './Player';
import { plural } from '@/libs/plural';
import Sunburst from '@/components/backgrounds/Sunburst.vue';
import DisplayMessage from '@/components/display/DisplayMessage.vue';
import GameEnd from '@/components/display/GameEnd.vue';
import Lobby from '@/components/display/Lobby.vue';
import ReactionMeter from '@/components/display/ReactionMeter.vue';
import Scoreboard from '@/components/display/Scoreboard.vue';
import TeamBars from '@/components/display/TeamBars.vue';
import TopicCard from '@/components/display/TopicCard.vue';
import { Music, TrackForCode } from '@/libs/music';

export default defineComponent({
  name: 'DisplayView',
  components: {
    DisplayMessage,
    FontAwesomeIcon,
    GameEnd,
    Lobby,
    ReactionMeter,
    Scoreboard,
    Sunburst,
    TeamBars,
    TopicCard,
  },
  mixins: [PlayerView],
  props: {
    volume: {
      type: Number,
      required: false,
      default: 50,
    },
  },
  data() {
    return {
      currentTrack: -1,
      music: null as HTMLAudioElement | null,
      musicCopyrightNotice: '',
    };
  },
  computed: {
    joinHost(): string {
      return window.location.host;
    },
  },
  watch: {
    gameCode(v: number) {
      this.codeChanged(v);
    },
    volume(v: number) {
      this.volumeChanged(v);
    },
  },
  created() {
    this.codeChanged(this.gameCode);
  },
  unmounted() {
    this.stopMusic();
  },
  methods: {
    plural,
    codeChanged(v: number) {
      const track = TrackForCode[v];
      if (track !== undefined && Music[track]) {
        this.playTrack(track);
      } else if (v !== this.GameCodes.Responses) {
        // Keep the topic music going while teams respond, otherwise stop
        this.stopMusic();
      }
    },
    playTrack(index: number) {
      if (this.currentTrack === index) {
        return;
      }
      const m = Music[index];
      if (!m) {
        return;
      }
      this.stopMusic();
      const music = new Audio(`${m.file}`);
      music.volume = this.realVolume(this.volume);
      music.loop = true;
      music.ondurationchange = () => {
        music.currentTime = Math.floor(Math.random() * music.duration);
      };
      // Browsers block audio until someone interacts with the page
      music.play().catch(() => {});
      this.music = music;
      if (m.attribution) {
        this.musicCopyrightNotice = `Music: ${m.name} - ${m.copyright}`;
      }
      this.currentTrack = index;
    },
    realVolume(v: number) {
      return (Math.exp(v / 100) - 1) / (Math.E - 1);
    },
    stopMusic() {
      if (this.music) {
        this.music.pause();
        this.music = null;
      }
      this.currentTrack = -1;
      this.musicCopyrightNotice = '';
    },
    volumeChanged(v: number) {
      if (this.music) {
        this.music.volume = this.realVolume(v);
      }
    },
  },
});
</script>

<style scoped lang="scss">
@use '@/libs/scss/animations/dance.scss';
@use '@/libs/scss/animations/popper.scss';

.main-display {
  left: 0;
  min-height: calc(100vh - 3.5rem);
  top: 3.5rem;
}
.presenting {
  font-size: calc(1.5rem + 3vw);
  font-weight: 800;
}
.trophy {
  color: #f2a516;
}
.musicCopyright {
  bottom: 10px;
  left: 0;
  margin: 0 auto;
  position: fixed;
  text-align: center;
  width: 100%;
}
</style>
