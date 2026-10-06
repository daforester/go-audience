<template>
  <div></div>
</template>

<script>
import Axios from 'axios';
import useLocalStore from "@/stores/local";

export default {
  name: 'PlayerView',
  props: {
    gameCode: {
      type: Number,
      required: true,
    },
    game: {
      type: Object,
      required: false,
      default() {
        return {
          Key: '',
          SuggestionsEnabled: false,
        };
      },
    },
    calling: {
      type: Object,
      required: false,
      default() {
        return {
          BadgeNumber: '',
          Name: '',
          UUID: '',
        };
      },
    },
    player: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      localStore: useLocalStore(),
      sockMessages: [
        { m: 'NEW_PLAYER', f: this.refreshGame },
        { m: 'PLAYER_QUIT', f: this.refreshGame },
        { m: 'PLAYER_UPDATE', f: this.playerUpdate },
        { m: 'OPTION_ADDED', f: this.refreshGame },
        { m: 'OPTION_ASSIGNED', f: this.refreshGame },
        {
          m: 'ACTIVATE_ROUND',
          f: (data) => {
            this.activateRound(data.Round);
          },
        },
        {
          m: 'CHECK_BADGE',
          f: (data) => {
            this.setBadgeCheck(data.Enabled);
          },
        },
        {
          m: 'CLEAR_GUESSES',
          f: (data) => {
            this.clearGuesses(data.Round);
          },
        },
        {
          m: 'GUESS_ADDED',
          f: (data) => {
            this.guessAdded(data.Guess);
          },
        },
        {
          m: 'GUESS_CHANGED',
          f: (data) => {
            this.guessChanged(data.Round, data.Index, data.Value);
          },
        },
        {
          m: 'GAME_PASSWORD_CHANGED',
          f: (data) => {
            this.gamePasswordChanged(data.Password);
          },
        },
        {
          m: 'PLAYER_NAME_CHANGED',
          f: (data) => {
            this.playerNameChanged(data.Index, data.Name);
          },
        },
        {
          m: 'PLAYER_BADGE_CHANGED',
          f: (data) => {
            this.badgeNumberChanged(data.Index, data.BadgeNumber);
          },
        },
        {
          m: 'ROUNDS_UPDATED',
          f: (data) => {
            this.roundsUpdated(data.Rounds);
          },
        },
        {
          m: 'CURRENT_ROUND',
          f: (data) => {
            this.currentRound(data);
          },
        },
        {
          m: 'COME_ON_DOWN',
          f: (data) => {
            if (data.Name) {
              this.callPlayer(data.UUID, data.Name, data.BadgeNumber, data.Value);
            } else {
              this.callPlayer(data.UUID, '', '', 0);
            }
          },
        },
      ],
    };
  },
  methods: {
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'PriceIsRight');
    },
    refreshGame() {
      const token = this.getToken();
      if (token) {
        this.$emit('gameReloading');
        Axios({
          method: 'get',
          url: `${this.apiUri}/game`,
          headers: {
            Authorization: `Bearer ${token.AccessToken}`,
          },
        }).then((r) => {
          if (r.data.Response.Game) {
            this.$emit('gameUpdated', r.data.Response.Game);
          }
        });
      }
    },
    activateRound(round) {
      this.$emit('roundActivated', round);
    },
    callPlayer(uuid, name, badge, value) {
      this.$emit('callPlayer', {
        UUID: uuid, Name: name, BadgeNumber: badge, Value: value
      });
    },
    clearGuesses(round) {
      this.$emit('guessesCleared', round);
    },
    currentRound(data) {
      this.$emit('currentRound', data.Round, data.RoundNumber);
    },
    badgeNumberChanged(index, value) {
      this.$emit('badgeNumberChanged', index, value);
    },
    gamePasswordChanged(value) {
      this.$emit('gamePasswordChanged', value);
    },
    guessAdded(guess) {
      this.$emit('guessAdded', guess);
    },
    guessChanged(round, index, value) {
      this.$emit('guessChanged', round, index, value);
    },
    playerUpdate(data) {
      this.$emit('playerUpdated', data);
    },
    playerNameChanged(index, value) {
      this.$emit('playerNameChanged', index, value);
    },
    roundsUpdated(rounds) {
      this.$emit('roundsUpdated', rounds);
    },
    setBadgeCheck(enabled) {
      this.$emit('badgeCheckChanged', enabled);
    },
    sockMsg(msg) {
      try {
        const data = JSON.parse(msg);
        const instruction = this.sockMessages.find((m) => m.m === data.Command);
        if (instruction) {
          if (data.Data) {
            instruction.f.call(this, data.Data);
          } else {
            instruction.f.call(this);
          }
        }
      } catch (e) {
        // console.log(e);
      }
    },
  },
};

</script>

<style scoped>

</style>
