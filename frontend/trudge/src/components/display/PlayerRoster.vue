<template>
  <div>
    <h1 class="text-center">
      Time for an adventure
    </h1>
    <h2 class="text-center">
      Sign Up Today!
    </h2>
    <h3 class="bg-white text-center m-2 p-2 border border-dark rounded rounded-lg mt-4">
      Room Code: {{ game.Key }}
    </h3>

    <div v-if="party.length" :key="`party${party.length}`">
      <h2 class="text-center">
        The Party
      </h2>
      <ul class="list-unstyled row">
        <li v-for="player in party" :key="player.UUID" class="playerIcon col h1 text-center">
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
          <div class="playerName">
            {{ player.Name }}
          </div>
        </li>
      </ul>
    </div>

    <div v-if="gods.length" class="mt-5" :key="`gods${gods.length}`">
      <h2 class="text-center">
        The Gods
      </h2>
      <ul class="list-unstyled row">
        <li v-for="player in gods" :key="player.UUID" class="playerIcon col h1 text-center">
          <div class="dance d-inline-block" title="God">
            <FontAwesomeIcon :icon="['fal', 'bolt']" />
          </div>
          <div class="playerName">
            {{ player.Name }}
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

export default {
  name: 'PlayerRoster',
  components: {
    FontAwesomeIcon,
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
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'Trudge');
    },
  },
};

</script>

<style scoped lang="scss">
@import '../../libs/scss/animations/dance.scss';
</style>
