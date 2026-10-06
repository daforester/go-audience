<template>
  <div>
    <div v-if="Player.UUID">
      <WebSocket v-if="Player.BadgeNumber && Player.UUID"
                 :wsuri="wsuri"
                 @sockmsg="sockMsg"
                 @sockopen="sockOpen"
      />
      <LeaderView v-if="Player.Class.toUpperCase() === 'LEADER'"
                  ref="playerView"
                  :game-code="GameCode"
                  :game="Game"
                  :player="Player"
                  :round="Round"
                  :calling="Calling"
                  @callPlayer="callPlayer"
                  @gameReloading="Game.loading = true"
                  @gameUpdated="updateGame"
                  @guessesCleared="clearGuesses"
                  @guessAdded="guessAdded"
                  @guessChanged="guessChanged"
                  @playerUpdated="updatePlayer"
                  @playerNameChanged="playerNameChanged"
                  @gamePasswordChanged="gamePasswordChanged"
                  @badgeNumberChanged="badgeNumberChanged"
                  @roundsUpdated="roundsUpdated"
                  @roundActivated="roundActivated"
                  @currentRound="currentRound"
                  @badgeCheckChanged="setBadgeCheck"
      />
      <DisplayView v-else-if="Player.Class.toUpperCase() === 'DISPLAY'"
                   ref="playerView"
                   :game-code="GameCode"
                   :game="Game"
                   :player="Player"
                   :round="Round"
                   :calling="Calling"
                   @callPlayer="callPlayer"
                   @gameReloading="Game.loading = true"
                   @gameUpdated="updateGame"
                   @guessAdded="guessAdded"
                   @guessChanged="guessChanged"
                   @playerUpdated="updatePlayer"
                   @playerNameChanged="playerNameChanged"
                   @badgeNumberChanged="badgeNumberChanged"
                   @roundsUpdated="roundsUpdated"
                   @roundActivated="roundActivated"
                   @currentRound="currentRound"
      />
      <AudienceView v-else-if="Player.Class.toUpperCase() === 'AUDIENCE'"
                    ref="playerView"
                    :game-code="GameCode"
                    :player="Player"
                    :suggestions-enabled="Game.SuggestionsEnabled"
                    :round="Round"
                    :calling="Calling"
                    @callPlayer="callPlayer"
                    @gameReloading="Game.loading = true"
                    @gameUpdated="updateGame"
                    @playerUpdated="updatePlayer"
                    @playerNameChanged="playerNameChanged"
                    @badgeNumberChanged="badgeNumberChanged"
                    @roundsUpdated="roundsUpdated"
                    @roundActivated="roundActivated"
                    @currentRound="currentRound"
                    @quitGame="quitGame"
      />
    </div>
    <div v-else>
      <LoginView @joinedGame="joinedGame" />
    </div>
  </div>
</template>
<script>
import Axios from 'axios';
import WebSocket from '../components/WebSocket.vue';
import LoginView from '../components/Login.vue';
import LeaderView from './Leader.vue';
import DisplayView from './Display.vue';
import AudienceView from './Audience.vue';
import useLocalStore from "@/stores/local";

export default {
  name: "GameView",
  components: {
    AudienceView,
    DisplayView,
    LeaderView,
    LoginView,
    WebSocket,
  },
  data() {
    return {
      localStore: useLocalStore(),
      GameCode: 0,
      Game: {
        Password: '',
        Players: [],
        Volume: 0,
        RoundNumber: 0,
        CurrentRound: {},
        Rounds: [],
        BadgeCheck: true,
        loading: false,
      },
      Player: {
        Class: '',
        Name: '',
        BadgeNumber: '',
        UUID: '',
      },
      Round: {
        Name: '',
        Image: '',
        UUID: '',
      },
      Calling: {
        Name: '',
        BadgeNumber: '',
        UUID: '',
        Value: 0,
      },
    };
  },
  computed: {
    wsuri() {
      return `${this.wsUri}/${this.Player.UUID}`;
    },
  },
  created() {
    this.emitter.on('sessionExpired', (token) => {
      if (token.Service === 'PriceIsRight') {
        // Token expired so leave game
        this.quitGame();
      }
    });
    this.localStore.managedTokens.forEach((token) => {
      if (token.Service === "PriceIsRight") {
        // We have a token so probably...
        // Already in a game!
        // Update token first
        this.emitter.emit("refreshTokensByService", "PriceIsRight");
        const rejoinGame = () => {
          this.emitter.off("tokenUpdated", rejoinGame);
          this.Player.UUID = ".";
          Axios({
            method: "get",
            url: `${this.apiUri}/game/rejoin`,
            headers: {
              Authorization: `Bearer ${token.AccessToken}`,
            },
          })
            .then((r) => {
              if (r.data.Response.Game) {
                this.Game = r.data.Response.Game;
                this.GameCode = this.Game.GameCode;
                this.Game.loading = false;
              } else if (r.data.Response.GameCode) {
                this.GameCode = r.data.Response.GameCode;
              }
            this.Player = r.data.Response.Player;
            })
            .catch(() => {
              this.Player.UUID = "";
          });
        }
        this.emitter.on("tokenUpdated", rejoinGame);
      }
    });
  },
  mounted() {

  },
  methods: {
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'PriceIsRight');
    },
    callPlayer(calling) {
      this.Calling = calling;
    },
    clearGuesses(round) {
      let r = 0;
      if (typeof round === 'number') {
        r = round;
      } else {
        r = this.Game.RoundNumber;
      }

      if (r < this.Game.Rounds.length && r >= 0) {
        this.Game.Rounds[r].Guesses = [];
      }
    },
    currentRound(round, roundNumber) {
      this.Game.CurrentRound = round;
      this.Round = round;
      this.Game.RoundNumber = roundNumber;
    },
    changeVolume() {
      // Sync Volume Change
      const token = this.getToken();
      if (token) {
        const data = {
          Volume: this.Game.Volume,
        };
        Axios({
          method: "post",
          url: `${this.apiUri}/game/volume`,
          headers: {
            Authorization: `Bearer ${token.AccessToken}`,
          },
          data,
        });
      }
    },
    joinedGame(player, game) {
      if (game !== undefined) {
        this.Game = game;
      }
      this.Player = player;
    },
    roundActivated(round) {
      this.Round = round;
      const er = this.Game.Rounds.find((r) => r.UUID === round.UUID);
      const ir = this.Game.Rounds.findIndex((r) => r.UUID === round.UUID);
      if (er) {
        this.Game.CurrentRound = er;
        this.Game.RoundNumber = ir;
      }
    },
    sockMsg(msg) {
      if (msg.substr(0, 4) === 'CODE') {
        this.GameCode = parseInt(msg.substr(5), 10);
        return;
      }
      if (msg.substr(0, 4) === 'QUIT') {
        this.quitGame();
        return;
      }
      this.$refs.playerView.sockMsg(msg);
    },
    sockOpen() {
      //
    },
    quitGame() {
      // Tell Game we quit & reset
      const token = this.getToken();
      if (token) {
        Axios({
          method: "get",
          url: `${this.apiUri}/game/quit`,
          headers: {
            Authorization: `Bearer ${token.AccessToken}`,
          },
        });
      }
      this.resetGame();
    },
    fetchGame() {
      // Obtain Game Data
      const token = this.getToken();
      if (token) {
        this.Game.loading = true;
        Axios({
          method: "get",
          url: `${this.apiUri}/game`,
          headers: {
            Authorization: `Bearer ${token.AccessToken}`,
          },
        }).then((r) => {
          if (r.data.Response.Game) {
            this.Game = r.data.Response.Game;
            this.Game.loading = false;
          }
        });
      }
    },
    getTokenManager() {
      if (this.$parent.$refs.TokenManager) {
        return this.$parent.$refs.TokenManager;
      }
      return undefined;
    },
    resetGame() {
      // Return to base state
      this.emitter.emit("removeTokensByService", "PriceIsRight");
      this.Player = {
        UUID: '',
        RoomKey: '',
      };
      this.Game = {
        Key: '',
        Players: [],
        loading: false,
      };
      this.GameCode = 0;
    },
    setBadgeCheck(enabled) {
      this.Game.BadgeCheck = enabled;
    },
    updateGame(game) {
      this.Game = game;
      this.Game.loading = false;
    },
    updatePlayer(player) {
      if (player.UUID === this.Player.UUID) {
        this.Player = player;
        const tokenUpdateWatcher = () => {
          if (this.Player.Class === 'Leader' || this.Player.Class === 'Display') {
            this.fetchGame();
          }
          this.emitter.off('tokenUpdated', tokenUpdateWatcher);
        };
        this.emitter.on('tokenUpdated', tokenUpdateWatcher);
        this.updateToken();
      }
      const index = this.Game.Players.findIndex((p) => player.UUID === p.UUID);
      if (index >= 0) {
        // Inject new player data
        this.Game.Players.splice(index, 1, player);
      } else {
        this.Game.Players.push(player);
      }
    },
    guessAdded(guess) {
      const round = this.Game.Rounds.find((r) => (r.UUID === guess.RoundUUID));
      round.Guesses.push(guess);
    },
    guessChanged(round, index, value) {
      this.Game.Rounds[round].Guesses[index].Value = value;
    },
    gamePasswordChanged(value) {
      this.Game.Password = value;
    },
    playerNameChanged(playerIndex, value) {
      if (this.Game.Players.length > playerIndex) {
        let player = this.Game.Players[playerIndex];
        if (player == null) {
          player = {};
        }
        player.Name = value;
        this.Game.Players.splice(playerIndex, 1, player);
      }
    },
    badgeNumberChanged(playerIndex, value) {
      if (this.Game.Players.length > playerIndex) {
        let player = this.Game.Players[playerIndex];
        if (player == null) {
          player = {};
        }
        player.BadgeNumber = value;
        this.Game.Players.splice(playerIndex, 1, player);
      }
    },
    roundsUpdated(value) {
      if (this.Game) {
        this.Game.Rounds = value;
      }
    },
    updateToken() {
      this.emitter.emit('refreshTokensByService', 'PriceIsRight');
    },
  },
};

</script>
<style scoped>
nav {
  z-index: 100;
}
.onHover .btn, .onHover .volumeControl {
  visibility: hidden;
}
.onHover:hover .btn, .onHover:hover .volumeControl {
  visibility: visible;
}
</style>
