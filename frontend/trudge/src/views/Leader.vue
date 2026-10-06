<template>
  <div>
    <div>
      <div class="row mb-2">
        <div class="col">
          <button v-if="gameCode === 20" class="btn btn-primary" @click="setGameCode(30)">
            Start Game
          </button>
          <template v-else>
            <button class="btn btn-danger" @click="setGameCode(20)">
              Reset Game
            </button>
            <button class="btn btn-danger" @click="setGameCode(600)">
              End Game
            </button>
            <button class="btn btn-success"
                    :disabled="game.SuggestionsEnabled ? 'disabled' : false"
                    @click="toggleSuggestions(true)"
            >
              Enable Suggestions
            </button>
            <button class="btn btn-warning"
                    :disabled="!game.SuggestionsEnabled ? 'disabled' : false"
                    @click="toggleSuggestions(false)"
            >
              Disable Suggestions
            </button>
          </template>
        </div>
      </div>
      <div v-if="gameCode >= 30 && gameCode <= 570" class="row mb-2">
        <div class="col">
          <button class="btn btn-primary" @click="setGameCode(50)">
            Show Warrior
          </button>
          <button class="btn btn-primary" @click="setGameCode(60)">
            Show Bard
          </button>
          <button class="btn btn-primary" @click="setGameCode(70)">
            Show Wizard
          </button>
          <button class="btn btn-primary" @click="setGameCode(80)">
            Show Gods
          </button>
        </div>
      </div>
      <div v-if="gameCode >= 30 && gameCode <= 570" class="row mb-2">
        <div class="col">
          <button class="btn btn-primary" @click="resetOptions(); setGameCode(90)">
            Request Weapons
          </button>
          <button class="btn btn-primary" @click="setGameCode(100)">
            Vote on Weapons
          </button>
          <button class="btn btn-primary" @click="assignWeapon(); setGameCode(110)">
            Show Weapon Result
          </button>
        </div>
      </div>
      <div v-if="gameCode >= 30 && gameCode <= 570" class="row mb-2">
        <div class="col">
          <button class="btn btn-primary" @click="resetOptions(); setGameCode(120)">
            Request Songs
          </button>
          <button class="btn btn-primary" @click="setGameCode(130)">
            Vote on Songs
          </button>
          <button class="btn btn-primary" @click="assignSong(); setGameCode(140)">
            Show Song Result
          </button>
        </div>
      </div>
      <div v-if="gameCode >= 30 && gameCode <= 570" class="row mb-2">
        <div class="col">
          <button class="btn btn-primary" @click="resetOptions(); setGameCode(150)">
            Request Spells
          </button>
          <button class="btn btn-primary" @click="setGameCode(160)">
            Vote on Spells
          </button>
          <button class="btn btn-primary" @click="assignSpell(); setGameCode(170)">
            Show Spell Result
          </button>
        </div>
      </div>
      <div v-if="gameCode >= 30 && gameCode <= 570" class="row mb-2">
        <div class="col">
          <button class="btn btn-primary" @click="resetOptions(); setGameCode(180)">
            Request Quests
          </button>
          <button class="btn btn-primary" @click="setGameCode(190)">
            Vote on Quest
          </button>
          <button class="btn btn-primary" @click="assignQuest(); setGameCode(200)">
            Show Quest Result
          </button>
        </div>
      </div>
      <div v-if="gameCode >= 30 && gameCode <= 570" class="row mb-2">
        <div class="col">
          <button class="btn btn-primary" @click="setGameCode(500)">
            Main Display
          </button>
          <button class="btn btn-primary" @click="setGameCode(510)">
            Dice
          </button>
          <button class="btn btn-primary" @click="resetOptions(); setGameCode(550)">
            Request Options
          </button>
          <button class="btn btn-primary" @click="setGameCode(560)">
            Vote on Options
          </button>
          <button class="btn btn-primary" @click="assignResult(); setGameCode(570)">
            Show Result
          </button>
        </div>
      </div>
    </div>
    <div v-if="ErrorMessage">
      {{ ErrorMessage }}
    </div>
    <div>
      <h5 class="d-inline">
        Quest
      </h5>
      <p class="d-inline">
        <TextInput class="textInput"
                   :value="game.Quest"
                   @change="saveQuest($event)"
        />
      </p>
    </div>
    <VotingOptions ref="votingOptions"
                   :game="game"
                   @optionsCleared="optionsCleared"
                   @optionChanged="optionChanged"
    />
    <PlayerAssignment :game="game"
                      @playerNameChanged="playerNameChanged"
                      @characterNameChanged="characterNameChanged"
                      @characterRaceChanged="characterRaceChanged"
                      @playerExperienceChanged="playerExperienceChanged"
                      @playerHealthChanged="playerHealthChanged"
    />
  </div>
</template>

<script>
import Axios from 'axios';
import PlayerView from './Player.vue';
import PlayerAssignment from '../components/gm/PlayerAssignment.vue';
import VotingOptions from '../components/gm/VotingOptions.vue';
import TextInput from '../components/inputs/TextInput.vue';
import useLocalStore from "@/stores/local";

export default {
  name: 'LeaderView',
  components: {
    PlayerAssignment,
    TextInput,
    VotingOptions,
  },
  extends: PlayerView,
  data() {
    return {
      ErrorMessage: '',
      localStore: useLocalStore(),
    };
  },
  methods: {
    assignWeapon() {
      this.assignOption('Weapon');
    },
    assignSong() {
      this.assignOption('Song');
    },
    assignSpell() {
      this.assignOption('Spell');
    },
    assignQuest() {
      this.assignOption('Quest');
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
    optionsCleared() {
      this.$emit('optionsCleared');
    },
    resetOptions() {
      this.$refs.votingOptions.clearOptions();
    },
    saveQuest(event) {
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/quest/change`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Quest: event.target.value,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.$emit('questChanged', event.target.value);
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
    toggleSuggestions(setting) {
      // Turn Audience Suggestions On or Off
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/suggestions`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Enabled: setting,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
    },
  },
};

</script>

<style scoped lang="scss">
.textInput {
  height: 1.5rem;
  min-width: 2rem;
}
</style>
