<template>
  <div>
    <IconDance
      v-if="
        [
          90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 550, 560,
          570,
        ].includes(gameCode)
      "
    />
    <Sunburst v-else />
    <AudienceSuggestions
      v-if="
        game.SuggestionsEnabled && [90, 120, 150, 180, 550].includes(gameCode)
      "
    />
    <div
      class="d-flex vh-100 vw-100 main-display position-absolute overflow-hidden"
    >
      <div class="m-auto">
        <GameState v-if="gameCode === 500" :game="game" />
        <PlayerRoster v-if="gameCode === 20" :game="game" />
        <DisplayMessage
          v-if="gameCode >= 30 && gameCode < 50"
          :game="game"
          effect="shimmer"
          message="Recruiting Heroes"
        />
        <PlayerReveal v-if="gameCode === 50" :game="game" type="warrior" />
        <PlayerReveal v-if="gameCode === 60" :game="game" type="bard" />
        <PlayerReveal v-if="gameCode === 70" :game="game" type="wizard" />
        <GodsReveal v-if="gameCode === 80" :game="game" />
        <DisplayMessage
          v-if="gameCode === 90"
          :game="game"
          effect="toaster"
          message="Visiting the Armoury"
        />
        <LiveVoting v-if="gameCode === 100" category="weapon" :game="game" />
        <VoteResult
          v-if="gameCode === 110"
          category="weapon"
          :game="game"
          label="Your warrior weapon shall be..."
        />
        <DisplayMessage
          v-if="gameCode === 120"
          :game="game"
          effect="toaster"
          message="Stringing Mandolins"
        />
        <LiveVoting v-if="gameCode === 130" category="song" :game="game" />
        <VoteResult
          v-if="gameCode === 140"
          category="song"
          :game="game"
          label="Your famous ballad is called..."
        />
        <DisplayMessage
          v-if="gameCode === 150"
          :game="game"
          effect="toaster"
          message="Scribing Scrolls"
        />
        <LiveVoting v-if="gameCode === 160" category="spells" :game="game" />
        <VoteResult
          v-if="gameCode === 170"
          category="spells"
          :game="game"
          label="This magic spell will surely help..."
        />
        <DisplayMessage
          v-if="gameCode === 180"
          :game="game"
          effect="toaster"
          message="Creating History"
        />
        <LiveVoting v-if="gameCode === 190" category="quest" :game="game" />
        <VoteResult
          v-if="gameCode === 200"
          category="quest"
          :game="game"
          label="Your quest shall be..."
        />
        <Dice v-if="gameCode === 510" />
        <DisplayMessage
          v-if="gameCode === 550"
          :game="game"
          effect="toaster"
          message="Consulting the Gods"
        />
        <LiveVoting v-if="gameCode === 560" :game="game" />
        <VoteResult v-if="gameCode === 570" :game="game" />
        <GameEnd v-if="gameCode === 600" :game="game" />
      </div>
    </div>
    <div v-if="musicCopyrightNotice" class="musicCopyright">
      {{ musicCopyrightNotice }}
    </div>
  </div>
</template>

<script>
import Dice from "../components/display/Dice.vue";
import GameEnd from "../components/display/GameEnd.vue";
import GameState from "../components/display/GameState.vue";
import DisplayMessage from "../components/display/DisplayMessage.vue";
import GodsReveal from "../components/display/GodsReveal.vue";
import PlayerReveal from "../components/display/PlayerReveal.vue";
import PlayerRoster from "../components/display/PlayerRoster.vue";
import PlayerView from "./Player.vue";
import LiveVoting from "../components/display/LiveVoting.vue";
import VoteResult from "../components/display/VoteResult.vue";
import Sunburst from "../components/backgrounds/Sunburst.vue";
import IconDance from "../components/backgrounds/IconDance.vue";
import AudienceSuggestions from "../components/backgrounds/AudienceSuggestions.vue";

import musicCathedral from "@/assets/music/46_Cathedral.mp3";
import musicHalfling from "@/assets/music/133_Halfling_Festival.mp3";
import musicWizards from "@/assets/music/174_Wizards_Tower.mp3";
import musicRoyal from "@/assets/music/175_Royal_Court.mp3";
import musicTavern from "@/assets/music/177_Tavern_Music.mp3";
import musicTinkerers from "@/assets/music/219_Tinkerers_Workshop.mp3";
import musicMedieval from "@/assets/music/Anon_Medieval_Dance_Tunes.mp3";

const Music = [
  {
    attribution: true,
    copyright: "tabletopaudio.com (CC BY-NC-ND 4.0)",
    file: musicCathedral,
    name: "Cathedral",
  },
  {
    attribution: true,
    copyright: "tabletopaudio.com (CC BY-NC-ND 4.0)",
    file: musicHalfling,
    name: "Halfling Festival",
  },
  {
    attribution: true,
    copyright: "tabletopaudio.com (CC BY-NC-ND 4.0)",
    file: musicWizards,
    name: "Wizards Tower",
  },
  {
    attribution: true,
    copyright: "tabletopaudio.com (CC BY-NC-ND 4.0)",
    file: musicRoyal,
    name: "Royal Court",
  },
  {
    attribution: true,
    copyright: "tabletopaudio.com (CC BY-NC-ND 4.0)",
    file: musicTavern,
    name: "Tavern Music",
  },
  {
    attribution: true,
    copyright: "tabletopaudio.com (CC BY-NC-ND 4.0)",
    file: musicTinkerers,
    name: "Tinkerers Workshop",
  },
  {
    attribution: false,
    copyright: "Public Domain",
    file: musicMedieval,
    name: "Medieval Dance Tunes",
  },
];

export default {
  name: "DisplayView",
  components: {
    AudienceSuggestions,
    IconDance,
    Sunburst,
    VoteResult,
    LiveVoting,
    Dice,
    GameState,
    DisplayMessage,
    GodsReveal,
    PlayerReveal,
    PlayerRoster,
    GameEnd,
  },
  extends: PlayerView,
  data() {
    return {
      currentTrack: -1,
      music: null,
      musicCopyrightNotice: "",
    };
  },
  watch: {
    gameCode(v) {
      this.codeChanged(v);
    },
    volume(v) {
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
    codeChanged(v) {
      if (v >= 20 && v <= 30) {
        this.playTrack(4); // Medieval Dance Tunes
      } else if (v >= 90 && v <= 100) {
        // Obtain Weapon
        this.playTrack(5); // Tinkerers Workshop
      } else if (v >= 120 && v <= 130) {
        // Obtain Song
        this.playTrack(1); // Half-ling Festival
      } else if (v >= 150 && v <= 160) {
        // Obtain Spell
        this.playTrack(2); // Wizards Tower
      } else if (v >= 180 && v <= 190) {
        // Obtain Quest
        this.playTrack(3); // Royal Court
      } else if (v >= 550 && v <= 560) {
        // Obtain Options
        this.playTrack(4); // Tavern Music
      } else {
        this.stopMusic();
      }
      if ([110, 140, 170, 200, 570].includes(v)) {
        // Update game because of results
        this.refreshGame();
      }
      /*
      else if (v === 510) { // Roll Dice
        this.playTrack(0); // Cathedral
      }
      */
    },
    playTrack(index) {
      if (this.currentTrack === index) {
        return;
      }
      this.stopMusic();
      const m = Music[index];
      this.music = new Audio(`${m.file}`);
      this.music.volume = this.volume / 100;
      this.music.loop = true;
      this.music.ondurationchange = () => {
        this.music.currentTime = Math.floor(
          Math.random() * this.music.duration
        );
      };
      this.music.play();
      if (m.attribution) {
        this.musicCopyrightNotice = `Music: ${m.name} - ${m.copyright}`;
      }
      this.currentTrack = index;
    },
    realVolume(v) {
      return (Math.exp(v / 100) - 1) / (Math.E - 1);
    },
    stopMusic() {
      if (this.music) {
        this.music.pause();
        this.music = null;
      }
      this.currentTrack = -1;
      this.musicCopyrightNotice = "";
    },
    volumeChanged(v) {
      if (this.music) {
        this.music.volume = this.realVolume(v);
      }
    },
  },
};
</script>

<style scoped lang="scss">
.main-display {
  top: 0;
  left: 0;
}
.musicCopyright {
  bottom: 10px;
  left: 0;
  margin: 0 auto;
  position: absolute;
  text-align: center;
  width: 100%;
}
</style>
