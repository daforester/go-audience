<template>
  <div class="container">
    <div>
      <div class="row mb-2">
        <button v-if="gameCode !== 1"
                class="btn btn-dark mr-2"
                type="button"
                @click="setGameCode(1)"
        >
          Coming Soon
        </button>
        <button v-else
                class="btn btn-dark mr-2"
                type="button"
                disabled
        >
          Coming Soon
        </button>
        <button v-if="gameCode !== 2"
                class="btn btn-light mr-2"
                type="button"
                @click="setGameCode(2)"
        >
          Starting Soon
        </button>
        <button v-else
                class="btn btn-light mr-2"
                type="button"
                disabled
        >
          Starting Soon
        </button>
        <button v-if="gameCode < 10"
                class="btn btn-primary mr-2"
                type="button"
                @click="setGameCode(10)"
        >
          Start Game
        </button>
        <button v-else
                class="btn btn-primary mr-2"
                type="button"
                disabled
        >
          Start Game
        </button>
        <button class="btn btn-danger mr-2" type="button" @click="setGameCode(100)">
          End Game
        </button>
        <button v-if="game.BadgeCheck"
                class="btn btn-success mr-2"
                type="button"
                @click="toggleBadgeCheck()"
        >
          Badge Check
        </button>
        <button v-else
                class="btn btn-danger mr-2"
                type="button"
                @click="toggleBadgeCheck()"
        >
          Badge Check
        </button>
      </div>
      <div v-if="gameCode >= 10 && gameCode < 100" class="row mb-2">
        <button v-if="gameCode !== 10"
                class="btn btn-primary mr-2"
                type="button"
                @click="setGameCode(10)"
        >
          Please Wait
        </button>
        <button v-else
                class="btn btn-primary mr-2"
                type="button"
                disabled
        >
          Please Wait
        </button>
        <button v-if="gameCode !== 20"
                class="btn btn-primary mr-2"
                type="button"
                @click="activateRound"
        >
          Enable Guesses
        </button>
        <button v-else
                class="btn btn-primary mr-2"
                type="button"
                disabled
        >
          Enable Guesses
        </button>
        <button class="btn btn-primary mr-2" type="button" @click="setGameCode(30)">
          Stop Guesses
        </button>
        <button class="btn btn-primary mr-2" type="button" @click="setGameCode(40)">
          Show Result
        </button>
      </div>
      <div v-if="gameCode >= 10 && gameCode < 100" class="row mb-2">
        <button class="btn btn-primary mr-2" type="button" @click="nextRound()">
          Next Round
        </button>
        <button class="btn btn-primary mr-2" type="button" @click="previousRound()">
          Previous Round
        </button>
      </div>
    </div>
    <div v-if="ErrorMessage">
      {{ ErrorMessage }}
    </div>
    <div>
      <h5 class="d-inline">
        Password
      </h5>
      <p class="d-inline">
        <TextInput class="textInput"
                   :value="game.Password"
                   @change="savePassword($event)"
        />
      </p>
    </div>
    <div>
      <h5>
        Current Round
      </h5>
      <p v-if="game.Rounds.length" class="d-inline">
        {{ game.RoundNumber + 1 }} / {{ game.Rounds.length }}
      </p>
      <p v-else>
        None created
      </p>
    </div>
    <nav>
      <ol>
        <li @click="Section = 'Players'">
          Players
        </li>
        <li @click="Section = 'Guesses'">
          Guesses
        </li>
        <li @click="Section = 'Rounds'">
          Rounds
        </li>
      </ol>
    </nav>
    <PlayerAssignment v-if="Section === 'Players'"
                      :game="game"
                      @playerNameChanged="playerNameChanged"
                      @badgeNumberChanged="badgeNumberChanged"
    />
    <Guesses v-else-if="Section === 'Guesses'"
             ref="guesses"
             :game="game"
             :round="game.RoundNumber"
             @guessesCleared="guessesCleared"
             @optionChanged="guessChanged"
    />
    <Rounds v-else-if="Section === 'Rounds'"
            ref="guesses"
            :game="game"
    />
  </div>
</template>

<script>
import Axios from 'axios';
import PlayerView from './Player.vue';
import PlayerAssignment from '../components/gm/PlayerAssignment.vue';
import Guesses from '../components/gm/Guesses.vue';
import Rounds from '../components/gm/Rounds.vue';
import TextInput from '../components/inputs/TextInput.vue';

export default {
  name: 'LeaderView',
  components: {
    PlayerAssignment,
    TextInput,
    Guesses,
    Rounds,
  },
  extends: PlayerView,
  data() {
    return {
      Section: 'Players',
      ErrorMessage: '',
    };
  },
  methods: {
    activateRound() {
      // Assign Voted Option Property
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/round/activate`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          RoundNumber: -1,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.setGameCode(20);
    },
    assignResult() {
      this.assignOption('Generic');
    },
    assignOption(property) {
      // Assign Voted Option Property
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/option/assign`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Property: property,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
    },
    guessesCleared() {
      this.$emit('guessesCleared');
    },
    nextRound() {
      this.setRound(this.game.RoundNumber + 1);
    },
    previousRound() {
      this.setRound(this.game.RoundNumber - 1);
    },
    resetOptions() {
      this.$refs.votingOptions.clearOptions();
    },
    savePassword(event) {
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/password`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Password: event.target.value,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.$emit('passwordChanged', event.target.value);
    },
    setGameCode(code) {
      // Set Game Code
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/code`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Code: code,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
    },
    setRound(r) {
      // Set Game Code
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/round/set`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          RoundNumber: r,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
    },
    toggleBadgeCheck() {
      // Turn badge Number Check On or Off
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/badgecheck`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Enabled: !this.game.BadgeCheck,
        },
      }).then(() => {
        this.$emit('badgeCheckChanged', !this.game.BadgeCheck);
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
    },
  },
};

</script>

<style scoped>
.textInput {
  height: 1.5rem;
  min-width: 2rem;
}
nav ol {
  display: flex;
  list-style: none;
  margin: 0;
  padding: 0;
}

nav ol li {
  border: 1px solid grey;
  border-radius: 6px 6px 0 0;
  cursor: pointer;
  margin: 0;
  padding: 0.25rem 1rem;
}
</style>
