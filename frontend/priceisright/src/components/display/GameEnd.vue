<template>
  <div class="endGame">
    <div v-if="revealStep === 0" class="popper text-center">
      <div class="endTitle">
        <span class="h1">
          Congratulations to...
        </span>
      </div>
    </div>
    <div v-else-if="revealStep === 1" class="popper text-center">
      <h1 class="awardTitle">
        <span>
          Leader of the Pack
        </span>
      </h1>
      <div v-if="player.Class.toUpperCase() === 'WARRIOR'"
           class="dance h1"
           title="Warrior"
      >
        <FontAwesomeIcon :icon="['fal', 'swords']" />
      </div>
      <div v-if="player.Class.toUpperCase() === 'BARD'"
           class="dance h1"
           title="Bard"
      >
        <FontAwesomeIcon :icon="['fal', 'mandolin']" />
      </div>
      <div v-if="player.Class.toUpperCase() === 'WIZARD'"
           class="dance h1"
           title="Wizard"
      >
        <FontAwesomeIcon :icon="['fal', 'wand-magic']" />
      </div>
      <div v-if="player.CharacterName" class="characterName h1">
        {{ player.CharacterName }}
      </div>
      <div v-if="player.CharacterRace" class="characterRace h1">
        <span class="small">
          the
        </span>
        <br>
        <span>
          {{ player.CharacterRace }}
        </span>
      </div>
      <div v-if="player.Experience" class="characterXP h2">
        <span>
          {{ player.Experience }} XP
        </span>
      </div>
    </div>
    <div v-else-if="revealStep === 2" class="popper text-center">
      <h1 class="awardTitle h1">
        <span>
          Most Followers
        </span>
      </h1>
      <div class="dance h1" title="God">
        <FontAwesomeIcon :icon="['fal', 'bolt']" />
      </div>
      <div v-if="player.CharacterName" class="characterName h1">
        {{ player.CharacterName }}
      </div>
      <div v-if="player.CharacterRace" class="characterRace h1">
        <span class="small">
          God of
        </span>
        <br>
        <span>
          {{ player.CharacterRace }}
        </span>
      </div>
      <div v-if="player.Experience" class="characterXP h2">
        <span>
          Mana: {{ player.Experience }}
        </span>
      </div>
    </div>
    <div v-else-if="revealStep === 3" class="popper text-center">
      <h1 class="awardTitle">
        <span>
          Dedicated Follower
        </span>
      </h1>
      <div class="playerName h1">
        {{ player.Name }}
      </div>
      <div class="characterXP h2">
        <span>
          {{ player.Experience }} XP
        </span>
      </div>
    </div>
    <div v-else class="popper text-center">
      <h1>
        Game Over
      </h1>
      <p class="h2">
        Thank you for playing
      </p>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

export default {
  name: 'GameState',
  components: {
    FontAwesomeIcon,
  },
  props: {
    game: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      revealStep: 0,
      player: {
        CharacterName: '',
        CharacterRace: '',
        Class: '',
        Experience: 0,
        Name: '',
      },
    };
  },
  mounted() {
    setInterval(() => {
      this.nextStep();
    }, 12000);
  },
  methods: {
    nextStep() {
      switch (this.revealStep) {
        default:
        case 0:
          this.player = this.getTopPartyMember();
          break;
        case 1:
          this.player = this.getTopGod();
          break;
        case 2:
          this.player = this.getTopAudience();
          break;
      }
      this.revealStep += 1;
    },
    getTopPlayer(classes) {
      let player = {
        CharacterName: '',
        CharacterRace: '',
        Class: '',
        Experience: 0,
        Name: '',
      };
      this.game.Players.forEach((p) => {
        if (classes.includes(p.Class.toUpperCase())) {
          if (p.Experience > player.Experience) {
            player = p;
          } else if (p.Experience === player.Experience) {
            if (Math.floor(Math.random() * Math.floor(1))) {
              player = p;
            }
          }
        }
      });
      return player;
    },
    getTopAudience() {
      return this.getTopPlayer(['AUDIENCE']);
    },
    getTopGod() {
      return this.getTopPlayer(['GOD']);
    },
    getTopPartyMember() {
      return this.getTopPlayer(['WARRIOR', 'BARD', 'WIZARD']);
    },
  },
};

</script>

<style scoped>
@import "../../libs/scss/animations/dance.scss";
@import "../../libs/scss/animations/popper.scss";

h1, .h1 {
  font-size: 4rem !important;
}

h2, .h2 {
  font-size: 3rem !important;
}
</style>
