<template>
  <div>
    <h1>Top Guesses</h1>
    <ul v-if="top && top.length > 0"
        class="list-unstyled"
    >
      <li v-for="(guess, i) in top" :key="i" class="row mb-1">
        {{ guess.Player.Name }} ({{ guess.Player.BadgeNumber }})
        <TextInput class="col"
                   input-class="form-control"
                   :value="guess.Value"
                   @change="saveGuess($event, i)"
        />
        <div class="col p-2">
          ({{ guess.Timestamp }})
        </div>
        <div class="col">
          <button class="btn btn-info" type="button" @click="summonPlayer(guess)">
            Summon
          </button>
        </div>
      </li>
    </ul>
    <button type="button" class="btn btn-primary" @click="topGuesses">
      Refresh
    </button>
    <button type="button" class="btn btn-info" @click="hidePlayer">
      Hide
    </button>
    <h1>Audience Guesses</h1>
    <ul v-if="game.Rounds.length > 0 &&
          game.RoundNumber < game.Rounds.length &&
          Array.isArray(game.Rounds[game.RoundNumber].Guesses) &&
          game.Rounds[game.RoundNumber].Guesses.length > 0"
        class="list-unstyled"
    >
      <li v-for="(guess, i) in game.Rounds[game.RoundNumber].Guesses" :key="i" class="row mb-1">
        {{ guess.Player.Name }} ({{ guess.Player.BadgeNumber }})
        <TextInput class="col"
                   input-class="form-control"
                   :value="guess.Value"
                   @change="saveGuess($event, i)"
        />
        <div class="col p-2">
          ({{ guess.Timestamp }})
        </div>
        <div class="col">
          <button class="btn btn-info" type="button" @click="summonPlayer(guess)">
            Summon
          </button>
        </div>
      </li>
    </ul>
    <button class="btn btn-danger" type="button" @click="clearGuesses">
      Clear
    </button>
  </div>
</template>

<script>
import Axios from "axios";
import TextInput from "../inputs/TextInput.vue";
import useLocalStore from "@/stores/local";

export default {
  name: 'Guesses',
  components: {
    TextInput,
  },
  props: {
    game: {
      type: Object,
      required: true,
    },
    round: {
      type: Number,
      required: true,
    },
  },
  data() {
    return {
      top: [],
      localStore: useLocalStore(),
    };
  },
  watch: {
    round() {
      this.topGuesses();
    },
  },
  mounted() {
    this.topGuesses();
  },
  methods: {
    clearGuesses() {
      // Clear Options & Votes
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/guesses/clear`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
      }).then(() => {
        this.$emit('guessesCleared');
      }).catch(() => {
        // Error with reset
      });
    },
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'PriceIsRight');
    },
    hidePlayer() {
      // Clear Options & Votes
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/summon`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Player: null,
          Value: 0,
        },
      }).then(() => {
        // Okay
      }).catch(() => {
        // Error with reset
      });
    },
    saveGuess(event, index) {
      // Change Option
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/guess/change/${index}`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Round: this.game.RoundNumber,
          Value: parseInt(event.target.value, 10),
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.$emit('optionChanged', index, event.target.value);
    },
    summonPlayer(guess) {
      // Clear Options & Votes
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/summon`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Player: guess.Player,
          Value: guess.Value,
        },
      }).then(() => {
        // Okay
      }).catch(() => {
        // Error with reset
      });
    },
    topGuesses() {
      // Clear Options & Votes
      const token = this.getToken();

      Axios({
        method: 'get',
        url: `${this.apiUri}/game/round/guesses/${this.round}`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
      }).then((r) => {
        this.top = r.data.Response.Guesses;
      }).catch(() => {
        // Error with reset
      });
    },
  },
};

</script>

<style scoped>

</style>
