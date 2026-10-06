<template>
  <div>
    <div v-if="Player.UUID">
      <WebSocket v-if="Player.RoomKey && Player.UUID"
                 :wsuri="wsuri"
                 @sockmsg="sockMsg"
                 @sockopen="sockOpen"
      />
      <nav class="navbar navbar-light bg-light">
        <div :class="`${Player.Class.toUpperCase() === 'DISPLAY' ? 'onHover' : ''}`">
          <button class="btn btn-dark" @click="quitGame">
            Disconnect
          </button>
        </div>
        <span v-if="Player.Class.toUpperCase() === 'LEADER'">
          DC: {{ GameCode }}
        </span>
        <span v-if="Player.Class.toUpperCase() === 'DISPLAY'" class="onHover">
          <input v-model.number="Game.Volume"
                 type="range"
                 min="0"
                 max="100"
                 class="volumeControl"
                 aria-label="Music Volume"
                 @change="changeVolume"
          >
        </span>
        <span v-if="Player.Class.toUpperCase() === 'LEADER'">
          <input v-model.number="Game.Volume"
                 type="range"
                 min="0"
                 max="100"
                 class="volumeControl"
                 aria-label="Music Volume"
                 @change="changeVolume"
          >
        </span>
        <span v-if="Game.Key" class="text-uppercase">
          Room Code: {{ Game.Key }}
        </span>
        <span v-if="Player.Class.toUpperCase() === 'WARRIOR'" class="">
          Your role: Warrior
        </span>
        <span v-else-if="Player.Class.toUpperCase() === 'BARD'" class="">
          Your role: Bard
        </span>
        <span v-else-if="Player.Class.toUpperCase() === 'WIZARD'" class="">
          Your role: Wizard
        </span>
        <span v-else-if="Player.Class.toUpperCase() === 'GOD'" class="">
          Your role: God
        </span>
        <span v-else-if="Player.Class.toUpperCase() === 'AUDIENCE'" class="">
          Your role: Audience
        </span>
      </nav>
      <LeaderView v-if="Player.Class.toUpperCase() === 'LEADER'"
                  ref="playerView"
                  :game-code="GameCode"
                  :game="Game"
                  :player="Player"
                  @setSuggestions="setSuggestions"
                  @experienceUpdated="experienceUpdated"
                  @healthUpdated="healthUpdated"
                  @gameReloading="Game.loading = true"
                  @gameUpdated="updateGame"
                  @optionsCleared="clearOptions"
                  @optionChanged="optionChanged"
                  @playerUpdated="updatePlayer"
                  @playerNameChanged="playerNameChanged"
                  @playerExperienceChanged="playerExperienceChanged"
                  @playerHealthChanged="playerHealthChanged"
                  @questChanged="questChanged"
                  @characterEquipmentAdded="characterEquipmentAdded"
                  @characterEquipmentChanged="characterEquipmentChanged"
                  @characterEquipmentRemoved="characterEquipmentRemoved"
                  @characterNameChanged="characterNameChanged"
                  @characterRaceChanged="characterRaceChanged"
                  @setVolume="setVolume"
      />
      <DisplayView v-else-if="Player.Class.toUpperCase() === 'DISPLAY'"
                   ref="playerView"
                   :game-code="GameCode"
                   :game="Game"
                   :player="Player"
                   :volume="Game.Volume"
                   @setSuggestions="setSuggestions"
                   @experienceUpdated="experienceUpdated"
                   @healthUpdated="healthUpdated"
                   @gameReloading="Game.loading = true"
                   @gameUpdated="updateGame"
                   @optionChanged="optionChanged"
                   @playerUpdated="updatePlayer"
                   @playerNameChanged="playerNameChanged"
                   @playerExperienceChanged="playerExperienceChanged"
                   @playerHealthChanged="playerHealthChanged"
                   @questChanged="questChanged"
                   @characterEquipmentAdded="characterEquipmentAdded"
                   @characterEquipmentChanged="characterEquipmentChanged"
                   @characterEquipmentRemoved="characterEquipmentRemoved"
                   @characterNameChanged="characterNameChanged"
                   @characterRaceChanged="characterRaceChanged"
                   @setVolume="setVolume"
      />
      <GodView v-else-if="Player.Class.toUpperCase() === 'GOD'"
               ref="playerView"
               :game-code="GameCode"
               :player="Player"
               @setSuggestions="setSuggestions"
               @experienceUpdated="experienceUpdated"
               @healthUpdated="healthUpdated"
               @gameReloading="Game.loading = true"
               @gameUpdated="updateGame"
               @playerUpdated="updatePlayer"
               @playerNameChanged="playerNameChanged"
               @playerExperienceChanged="playerExperienceChanged"
               @playerHealthChanged="playerHealthChanged"
               @questChanged="questChanged"
               @characterEquipmentAdded="characterEquipmentAdded"
               @characterEquipmentChanged="characterEquipmentChanged"
               @characterEquipmentRemoved="characterEquipmentRemoved"
               @characterNameChanged="characterNameChanged"
               @characterRaceChanged="characterRaceChanged"
      />
      <PartyMemberView v-else-if="Player.Class.toUpperCase() === 'BARD'"
                       ref="playerView"
                       :game-code="GameCode"
                       :player="Player"
                       @setSuggestions="setSuggestions"
                       @experienceUpdated="experienceUpdated"
                       @healthUpdated="healthUpdated"
                       @gameReloading="Game.loading = true"
                       @gameUpdated="updateGame"
                       @playerUpdated="updatePlayer"
                       @playerNameChanged="playerNameChanged"
                       @playerExperienceChanged="playerExperienceChanged"
                       @playerHealthChanged="playerHealthChanged"
                       @questChanged="questChanged"
                       @characterEquipmentAdded="characterEquipmentAdded"
                       @characterEquipmentChanged="characterEquipmentChanged"
                       @characterEquipmentRemoved="characterEquipmentRemoved"
                       @characterNameChanged="characterNameChanged"
                       @characterRaceChanged="characterRaceChanged"
      />
      <PartyMemberView v-else-if="Player.Class.toUpperCase() === 'WARRIOR'"
                       ref="playerView"
                       :game-code="GameCode"
                       :player="Player"
                       @setSuggestions="setSuggestions"
                       @experienceUpdated="experienceUpdated"
                       @healthUpdated="healthUpdated"
                       @gameReloading="Game.loading = true"
                       @gameUpdated="updateGame"
                       @playerUpdated="updatePlayer"
                       @playerNameChanged="playerNameChanged"
                       @playerExperienceChanged="playerExperienceChanged"
                       @playerHealthChanged="playerHealthChanged"
                       @questChanged="questChanged"
                       @characterEquipmentAdded="characterEquipmentAdded"
                       @characterEquipmentChanged="characterEquipmentChanged"
                       @characterEquipmentRemoved="characterEquipmentRemoved"
                       @characterNameChanged="characterNameChanged"
                       @characterRaceChanged="characterRaceChanged"
      />
      <PartyMemberView v-else-if="Player.Class.toUpperCase() === 'WIZARD'"
                       ref="playerView"
                       :game-code="GameCode"
                       :player="Player"
                       @setSuggestions="setSuggestions"
                       @experienceUpdated="experienceUpdated"
                       @healthUpdated="healthUpdated"
                       @gameReloading="Game.loading = true"
                       @gameUpdated="updateGame"
                       @playerUpdated="updatePlayer"
                       @playerNameChanged="playerNameChanged"
                       @playerExperienceChanged="playerExperienceChanged"
                       @playerHealthChanged="playerHealthChanged"
                       @questChanged="questChanged"
                       @characterEquipmentAdded="characterEquipmentAdded"
                       @characterEquipmentChanged="characterEquipmentChanged"
                       @characterEquipmentRemoved="characterEquipmentRemoved"
                       @characterNameChanged="characterNameChanged"
                       @characterRaceChanged="characterRaceChanged"
      />
      <AudienceView v-else-if="Player.Class.toUpperCase() === 'AUDIENCE'"
                    ref="playerView"
                    :game-code="GameCode"
                    :player="Player"
                    :suggestions-enabled="Game.SuggestionsEnabled"
                    @setSuggestions="setSuggestions"
                    @experienceUpdated="experienceUpdated"
                    @healthUpdated="healthUpdated"
                    @gameReloading="Game.loading = true"
                    @gameUpdated="updateGame"
                    @playerUpdated="updatePlayer"
                    @playerNameChanged="playerNameChanged"
                    @playerExperienceChanged="playerExperienceChanged"
                    @playerHealthChanged="playerHealthChanged"
                    @questChanged="questChanged"
                    @characterEquipmentAdded="characterEquipmentAdded"
                    @characterEquipmentChanged="characterEquipmentChanged"
                    @characterEquipmentRemoved="characterEquipmentRemoved"
                    @characterNameChanged="characterNameChanged"
                    @characterRaceChanged="characterRaceChanged"
      />
    </div>
    <div v-else>
      <LoginView @joinedGame="joinedGame"
                 @createdGame="createdGame"
      />
    </div>
  </div>
</template>
<script>
import Axios from 'axios';
import WebSocket from '../components/WebSocket.vue';
import LoginView from '../components/Login.vue';
import LeaderView from './Leader.vue';
import DisplayView from './Display.vue';
import GodView from './God.vue';
import AudienceView from './Audience.vue';
import PartyMemberView from './PartyMember.vue';
import useLocalStore from "@/stores/local";

export default {
  name: 'GameView',
  components: {
    PartyMemberView,
    AudienceView,
    GodView,
    DisplayView,
    LeaderView,
    LoginView,
    WebSocket,
  },
  data() {
    return {
      GameCode: 0,
      Game: {
        Key: '',
        Options: [],
        Players: [],
        Prompt: '',
        Quest: '',
        Votes: [],
        SuggestionsEnabled: false,
        Volume: 0,
        loading: false,
      },
      Player: {
        CharacterName: '',
        CharacterRace: '',
        Class: '',
        Equipment: [],
        Experience: 0,
        Health: 0,
        Name: '',
        Option: -1,
        RoomKey: '',
        UUID: '',
        Vote: -1,
      },
      localStore: useLocalStore(),
    };
  },
  computed: {
    wsuri() {
      return `${this.wsUri}/${this.Player.RoomKey}/${this.Player.UUID}`;
    },
  },
  created() {
    this.emitter.on('sessionExpired', (token) => {
      if (token.Service === 'Trudge') {
        // Token expired so leave game
        this.quitGame();
      }
    });
    this.emitter.on("voteAdded", this.voteAdded);
    this.localStore.managedTokens.forEach((token) => {
      if (token.Service === "Trudge") {
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
      return this.localStore.managedTokens.find((token) => token.Service === 'Trudge');
    },
    clearOptions() {
      this.Game.Options = [];
      this.Game.Votes = [];
      this.Game.Prompt = '';
    },
    changeVolume() {
      // Sync Volume Change
      const token = this.getToken();
      if (token) {
        const data = {
          Volume: this.Game.Volume,
        };
        Axios({
          method: 'post',
          url: `${this.apiUri}/game/volume`,
          headers: {
            Authorization: `Bearer ${token.AccessToken}`,
          },
          data,
        });
      }
    },
    createdGame(game, player) {
      this.Game = game;
      this.Player = player;
    },
    joinedGame(player) {
      this.Player = player;
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
          method: 'get',
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
          method: 'get',
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
      this.emitter.emit('removeTokensByService', 'Trudge');
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
          this.emitter.off("tokenUpdated", tokenUpdateWatcher);
        };
        this.emitter.on("tokenUpdated", tokenUpdateWatcher);
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
    experienceUpdated(value) {
      this.Player.Experience = value;
    },
    healthUpdated(value) {
      this.Player.Health = value;
    },
    optionChanged(index, value) {
      this.Game.Options.splice(index, 1, value);
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
    playerExperienceChanged(playerIndex, value) {
      if (this.Game.Players.length > playerIndex) {
        let player = this.Game.Players[playerIndex];
        if (player == null) {
          player = {};
        }
        player.Experience = value;
        this.Game.Players.splice(playerIndex, 1, player);
      }
    },
    playerHealthChanged(playerIndex, value) {
      if (this.Game.Players.length > playerIndex) {
        let player = this.Game.Players[playerIndex];
        if (player == null) {
          player = {};
        }
        player.Health = value;
        this.Game.Players.splice(playerIndex, 1, player);
      }
    },
    characterEquipmentAdded(playerIndex, value, uuid) {
      if (this.Game.Players.length > playerIndex) {
        let player = this.Game.Players[playerIndex];
        if (player == null) {
          player = {};
        }
        if (player.Equipment === null) {
          player.Equipment = [];
        }
        player.Equipment.push(value);
        this.Game.Players.splice(playerIndex, 1, player);
      }
      if (this.Player.UUID === uuid) {
        if (this.Player.Equipment === null) {
          this.Player.Equipment = [];
        }
        this.Player.Equipment.push(value);
      }
    },
    characterEquipmentChanged(playerIndex, equipmentIndex, value, uuid) {
      if (this.Game.Players.length > playerIndex) {
        let player = this.Game.Players[playerIndex];
        if (player == null) {
          player = {};
        }
        if (player.Equipment === null) {
          player.Equipment = [];
        }
        player.Equipment.splice(equipmentIndex, 1, value);
        this.Game.Players.splice(playerIndex, 1, player);
      }
      if (this.Player.UUID === uuid) {
        if (this.Player.Equipment === null) {
          this.Player.Equipment = [];
        }
        this.Player.Equipment.splice(equipmentIndex, 1, value);
      }
    },
    characterEquipmentRemoved(playerIndex, equipmentIndex, uuid) {
      if (this.Game.Players.length > playerIndex) {
        let player = this.Game.Players[playerIndex];
        if (player == null) {
          player = {};
        }
        if (player.Equipment === null) {
          player.Equipment = [];
        }
        player.Equipment.splice(equipmentIndex, 1);
        this.Game.Players.splice(playerIndex, 1, player);
      }
      if (this.Player.UUID === uuid) {
        if (this.Player.Equipment === null) {
          this.Player.Equipment = [];
        }
        this.Player.Equipment.splice(equipmentIndex, 1);
      }
    },
    characterNameChanged(playerIndex, value, uuid) {
      if (this.Game.Players.length > playerIndex) {
        let player = this.Game.Players[playerIndex];
        if (player == null) {
          player = {};
        }
        player.CharacterName = value;
        this.Game.Players.splice(playerIndex, 1, player);
      }
      if (this.Player.UUID === uuid) {
        this.Player.CharacterName = value;
      }
    },
    characterRaceChanged(playerIndex, value, uuid) {
      if (this.Game.Players.length > playerIndex) {
        let player = this.Game.Players[playerIndex];
        if (player == null) {
          player = {};
        }
        player.CharacterRace = value;
        this.Game.Players.splice(playerIndex, 1, player);
      }
      if (this.Player.UUID === uuid) {
        this.Player.CharacterRace = value;
      }
    },
    questChanged(value) {
      if (this.Game) {
        this.Game.Quest = value;
      }
    },
    setVolume(value) {
      this.Game.Volume = value;
    },
    setSuggestions(value) {
      this.Game.SuggestionsEnabled = value;
    },
    updateToken() {
      this.emitter.emit("refreshTokensByService", "Trudge");
    },
    voteAdded(vote) {
      if (this.Game.Options[vote.index] === vote.option) {
        if (!this.Game.Votes) {
          this.Game.Votes = Array(this.Game.Options.length).fill(0);
        }
        const votes = this.Game.Votes[vote.index];
        if (!votes.Total) {
          votes.Total = 0;
        }
        votes.Total += 1;
        if (!votes[vote.player.Class]) {
          votes[vote.player.Class] = 0;
        }
        votes[vote.player.Class] += 1;
        this.Game.Votes.splice(vote.index, 1, votes);
      }
    },
  },
};

</script>
<style scoped lang="scss">
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
