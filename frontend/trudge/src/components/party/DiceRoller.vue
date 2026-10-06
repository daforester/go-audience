<template>
  <div>
    <button v-if="!Rolled" class="btn btn-primary" @click.prevent="rollDice()">
      Roll Dice
    </button>
    <div v-else>
      <PleaseWait :player="player" :game-code="gameCode" />
    </div>
  </div>
</template>

<script>
import Axios from 'axios';
import PleaseWait from '../player/PleaseWait.vue';

export default {
  name: 'DiceRoller',
  components: {
    PleaseWait,
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
  data() {
    return {
      Rolled: false,
    };
  },
  created() {
    this.emitter.on('rollDice', this.diceRolled);
  },
  unmounted() {
    this.emitter.off('rollDice', this.diceRolled);
  },
  methods: {
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'Trudge');
    },
    diceRolled() {
      this.Rolled = true;
    },
    rollDice() {
      // Obtain Game Data
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/dice/roll`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
      }).then(() => {
        this.Rolled = true;
      }).catch(() => {
        this.Rolled = false;
      });
    },
  },
};

</script>

<style scoped>

</style>
