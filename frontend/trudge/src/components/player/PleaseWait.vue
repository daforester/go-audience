<template>
  <div>
    <div v-if="gameCode > 20 && ['BARD', 'WARRIOR', 'WIZARD'].includes(player.Class.toUpperCase())">
      <div class="dance h1 mb-3 text-center">
        <span v-if="player.Class.toUpperCase() === 'WARRIOR'">
          <FontAwesomeIcon :icon="['fal', 'swords']" />
        </span>
        <span v-if="player.Class.toUpperCase() === 'BARD'">
          <FontAwesomeIcon :icon="['fal', 'mandolin']" />
        </span>
        <span v-if="player.Class.toUpperCase() === 'WIZARD'">
          <FontAwesomeIcon :icon="['fal', 'wand-magic']" />
        </span>
      </div>
      <div v-if="player.CharacterName"
           class="bg-white border border-dark mb-2 p-2 rounded rounded-lg text-center"
      >
        <div class="h2">
          {{ player.CharacterName }}
        </div>
        <div v-if="player.CharacterRace" class="h2">
          the
        </div>
        <div v-if="player.CharacterRace" class="h2 m-0">
          {{ player.CharacterRace }}
        </div>
      </div>
      <div v-if="player.Equipment"
           class="bg-white border border-dark h5 mb-2 p-2 rounded rounded-lg text-left"
      >
        <span v-if="player.Class.toUpperCase() === 'WARRIOR'">
          <strong>Weapon:</strong> {{ formatArray(player.Equipment) }}
        </span>
        <span v-if="player.Class.toUpperCase() === 'BARD'">
          <strong>Song:</strong> {{ formatArray(player.Equipment) }}
        </span>
        <span v-if="player.Class.toUpperCase() === 'WIZARD'">
          <strong>Spell:</strong> {{ formatArray(player.Equipment) }}
        </span>
      </div>
      <div class="bg-white border border-dark h2 mb-2 rounded rounded-lg text-center">
        <template v-for="i in player.Health" :key="`health${i}`">
          <FontAwesomeIcon :icon="['fas', 'heart']" class="m-1 text-danger"/>
        </template>
        <template v-for="i in 5 - player.Health" :key="`death${i}`">
          <FontAwesomeIcon :icon="['fal', 'heart']" class="m-1" />
        </template>
      </div>
      <div class="bg-white border border-dark h5 mb-2 p-2 rounded rounded-lg text-left">
        <strong>XP:</strong> {{ player.Experience }}
      </div>
    </div>
    <div v-if="gameCode > 20 && ['GOD'].includes(player.Class.toUpperCase())">
      <div class="dance h1 mb-3 text-center">
        <span v-if="player.Class.toUpperCase() === 'GOD'">
          <FontAwesomeIcon :icon="['fal', 'bolt']" />
        </span>
      </div>
      <div v-if="player.CharacterName"
           class="bg-white border border-dark mb-2 p-2 rounded rounded-lg text-center"
      >
        <div class="h2" v-if="player.CharacterName">
          {{ player.CharacterName }}
        </div>
        <div v-if="player.CharacterRace" class="h2">
          god of
        </div>
        <div v-if="player.CharacterRace" class="h2 m-0">
          {{ player.CharacterRace }}
        </div>
      </div>
      <div class="bg-white border border-dark h5 mb-2 p-2 rounded rounded-lg text-left">
        <strong>Mana:</strong> {{ player.Experience }}
      </div>
    </div>
    <div v-if="gameCode > 20 && ['AUDIENCE'].includes(player.Class.toUpperCase())">
      <div class="dance h1 mb-3 text-center">
        <span v-if="player.Class.toUpperCase() === 'AUDIENCE'">
          <FontAwesomeIcon :icon="['fal', 'eye']" />
        </span>
      </div>
      <div class="bg-white border border-dark mb-2 p-2 rounded rounded-lg text-center">
        <div class="h2 m-0" v-if="player.Name">
          {{ player.Name }}
        </div>
      </div>
      <div class="bg-white border border-dark h5 mb-2 p-2 rounded rounded-lg text-left">
        <strong>XP:</strong> {{ player.Experience }}
      </div>
    </div>
    <DisplayMessage effect="toaster"
                    message="Please Wait..."
    />
  </div>
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import DisplayMessage from '../display/DisplayMessage.vue';

export default {
  name: 'PleaseWait',
  components: {
    DisplayMessage,
    FontAwesomeIcon,
  },
  props: {
    gameCode: {
      type: Number,
      required: true,
    },
    player: {
      type: Object,
      required: true,
    },
  },
  methods: {
    formatArray(a) {
      if (a) {
        return a.toString();
      }
      return '';
    },
  },
};

</script>

<style scoped lang="scss">
@import "../../libs/scss/animations/dance.scss";
</style>
