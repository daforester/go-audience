<template>
  <div>
    <div>
      <h1 class="text-center">
        The Quest
      </h1>
      <p class="text-center h5">
        {{ game.Quest }}
      </p>
    </div>
    <div>
      <h1 class="text-center">
        The Party
      </h1>
      <ul class="list-unstyled d-flex justify-content-center">
        <li v-for="player in party"
            :key="player.UUID"
            class="playerIcon bg-white border border-dark rounded rounded-lg h1 m-2 p-4 text-center"
        >
          <div v-if="player.Class.toUpperCase() === 'WARRIOR'"
               class="dance d-inline-block"
               title="Warrior"
          >
            <FontAwesomeIcon :icon="['fal', 'swords']" />
          </div>
          <div v-if="player.Class.toUpperCase() === 'BARD'"
               class="dance d-inline-block"
               title="Bard"
          >
            <FontAwesomeIcon :icon="['fal', 'mandolin']" />
          </div>
          <div v-if="player.Class.toUpperCase() === 'WIZARD'"
               class="dance d-inline-block"
               title="Wizard"
          >
            <FontAwesomeIcon :icon="['fal', 'wand-magic']" />
          </div>
          <div class="playerClass mb-2">
            <span>{{ player.Class }}</span>
          </div>
          <div class="characterName small">
            {{ player.CharacterName }}
          </div>
          <div class="characterRace small">
            <span class="small">the</span>
            <br>
            <span>{{ player.CharacterRace }}</span>
          </div>
          <div v-if="player.Class.toUpperCase() === 'WARRIOR'"
               class="d-inline-block equipment"
          >
            Wielding: {{ formatArray(player.Equipment) }}
          </div>
          <div v-if="player.Class.toUpperCase() === 'BARD'"
               class="d-inline-block equipment"
          >
            Performing: {{ formatArray(player.Equipment) }}
          </div>
          <div v-if="player.Class.toUpperCase() === 'WIZARD'"
               class="d-inline-block equipment"
          >
            Casting: {{ formatArray(player.Equipment) }}
          </div>
          <div class="text-center">
            <template v-for="h in player.Health" :key="`health${player.UUID}_${h}`">
              <FontAwesomeIcon :icon="['fas', 'heart']"
                               class="m-1 text-danger"
              />
            </template>
            <template v-for="h in 5 - player.Health" :key="`death${player.UUID}_${h}`">
              <FontAwesomeIcon :icon="['fal', 'heart']"
                               class="m-1 text-dark"
              />
            </template>
          </div>
        </li>
      </ul>
    </div>
    <div>
      <h1 class="text-center">
        The Gods
      </h1>
      <ul class="list-unstyled d-flex justify-content-center flex-wrap">
        <li v-for="player in gods"
            :key="player.UUID"
            class="playerIcon bg-white border border-dark rounded rounded-lg h1 m-2 p-4 text-center"
        >
          <div class="dance d-inline-block" title="God">
            <FontAwesomeIcon :icon="['fal', 'bolt']" />
          </div>
          <div class="characterName small">
            {{ player.CharacterName }}
          </div>
          <div class="characterRace small">
            <span class="small">God of</span>
            <br>
            <span>{{ player.CharacterRace }}</span>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import useLocalStore from "@/stores/local";

export default {
  name: 'GameState',
  components: {
    FontAwesomeIcon,
  },
  data() {
    return {
      localStore: useLocalStore(),
    };
  },
  props: {
    game: {
      type: Object,
      required: true,
    },
  },
  computed: {
    party() {
      return this.game.Players.filter((p) => p.Class.toUpperCase() === 'WARRIOR' || p.Class.toUpperCase() === 'WIZARD' || p.Class.toUpperCase() === 'BARD');
    },
    gods() {
      return this.game.Players.filter((p) => p.Class.toUpperCase() === 'GOD');
    },
  },
  methods: {
    formatArray(a) {
      return a.toString();
    },
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'PriceIsRight');
    },
  },
};

</script>

<style scoped>
@import "../../libs/scss/animations/dance.scss";
.equipment {
  font-size: 1.5rem;
}
</style>
