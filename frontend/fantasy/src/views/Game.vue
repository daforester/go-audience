<template>
  <div>
    <div v-if="Player.UUID">
      <WebSocket v-if="Player.RoomKey && Player.UUID && Player.UUID !== '.'"
                 :wsuri="wsuri"
                 @sockmsg="sockMsg"
      />
      <nav class="navbar navbar-light bg-light px-3 gap-2">
        <div :class="isClass('Display') ? 'onHover' : ''">
          <button class="btn btn-dark btn-sm" @click="confirmQuit">
            <FontAwesomeIcon :icon="['fas', 'right-from-bracket']" />
            Disconnect
          </button>
        </div>
        <span v-if="isClass('Admin') || isClass('Display')"
              :class="isClass('Display') ? 'onHover' : ''"
        >
          <input v-model.number="Game.Volume"
                 type="range"
                 min="0"
                 max="100"
                 class="form-range volumeControl"
                 aria-label="Music Volume"
                 @change="changeVolume"
          >
        </span>
        <span v-if="Game.Key" class="text-uppercase">
          Room Code: <strong>{{ Game.Key }}</strong>
        </span>
        <span v-if="Player.Class">
          <FontAwesomeIcon :icon="['fas', roleIcon]" />
          {{ roleLabel }}
          <template v-if="isClass('Director') && store.myTeam">
            &middot;
            <span class="team-swatch" :style="{ background: store.teamColour(Player.TeamID) }" />
            {{ store.myTeam.Name }}
          </template>
        </span>
      </nav>
      <AdminView v-if="isClass('Admin')" :game-code="GameCode" />
      <DisplayView v-else-if="isClass('Display')" :game-code="GameCode" :volume="Game.Volume" />
      <StaffView v-else-if="isClass('Staff')" :game-code="GameCode" />
      <DirectorView v-else-if="isClass('Director')" :game-code="GameCode" />
      <AudienceView v-else-if="isClass('Audience')" :game-code="GameCode" />
    </div>
    <div v-else>
      <LoginView @joined-game="joinedGame"
                 @created-game="createdGame"
      />
    </div>
  </div>
</template>
<script lang="ts">
import { defineComponent } from 'vue';
import Axios from 'axios';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import WebSocket from '@/components/WebSocket.vue';
import LoginView from '@/components/Login.vue';
import AdminView from './Admin.vue';
import AudienceView from './Audience.vue';
import DirectorView from './Director.vue';
import DisplayView from './Display.vue';
import StaffView from './Staff.vue';
import GameRequest, { SERVICE } from '@/mixins/api/GameRequest';
import useGameStore from '@/stores/game';
import type { Game, Player, PlayerClass } from '@/types/game';
import type { ManagedToken } from '@/stores/local';

const RoleLabels: Record<PlayerClass, string> = {
  Admin: 'Admin',
  Audience: 'Audience',
  Director: 'Con Director',
  Display: 'Main Display',
  Staff: 'Staff Team',
};

const RoleIcons: Record<PlayerClass, string> = {
  Admin: 'clipboard-list',
  Audience: 'paw',
  Director: 'user-tie',
  Display: 'tv',
  Staff: 'id-badge',
};

export default defineComponent({
  name: 'GameView',
  components: {
    AdminView,
    AudienceView,
    DirectorView,
    DisplayView,
    FontAwesomeIcon,
    LoginView,
    StaffView,
    WebSocket,
  },
  mixins: [GameRequest],
  data() {
    return {
      store: useGameStore(),
    };
  },
  computed: {
    Game(): Game {
      return this.store.Game;
    },
    GameCode(): number {
      return this.store.GameCode;
    },
    Player(): Player {
      return this.store.Player;
    },
    roleLabel(): string {
      return this.Player.Class ? RoleLabels[this.Player.Class] : '';
    },
    roleIcon(): string {
      return this.Player.Class ? RoleIcons[this.Player.Class] : 'paw';
    },
  },
  created() {
    this.emitter.on('sessionExpired', (token: ManagedToken) => {
      if (token.Service === SERVICE) {
        // Token expired so leave game
        this.resetGame();
      }
    });

    const token = this.getToken();
    if (token) {
      // We have a token so probably...
      // Already in a game!
      // Update token first
      this.store.Player.UUID = '.';
      const rejoinGame = () => {
        this.emitter.off('tokenUpdated', rejoinGame);
        this.rejoinGame();
      };
      this.emitter.on('tokenUpdated', rejoinGame);
      this.emitter.emit('refreshTokensByService', SERVICE);
    }
  },
  beforeUnmount() {
    this.emitter.off('sessionExpired');
  },
  methods: {
    isClass(c: PlayerClass): boolean {
      return this.Player.Class === c;
    },
    wsuri(): string {
      const token = this.getToken();
      const accessToken = token ? encodeURIComponent(token.AccessToken) : '';
      return `${this.wsUri}/${this.Player.RoomKey}/${this.Player.UUID}?token=${accessToken}`;
    },
    rejoinGame() {
      this.gameRequest<{ Game: Game; Player: Player }>('get', '/game/rejoin')
        .then((r) => {
          this.store.setGame(r.Game);
          this.store.setPlayer(r.Player);
        })
        .catch(() => {
          this.resetGame();
        });
    },
    changeVolume() {
      // Sync Volume Change
      this.gameRequest('post', '/game/volume', {
        Volume: this.Game.Volume,
      }).catch(() => {});
    },
    createdGame(game: Game, player: Player) {
      this.store.setGame(game);
      this.store.setPlayer(player);
    },
    joinedGame(player: Player) {
      this.store.setPlayer(player);
    },
    sockMsg(msg: string) {
      const before = this.Player.Class;
      if (this.store.sockMsg(msg) === 'quit') {
        this.resetGame();
        return;
      }
      if (this.Player.Class !== before) {
        // Our role changed, keep the token's copy of us current
        this.emitter.emit('refreshTokensByService', SERVICE);
      }
    },
    confirmQuit() {
      const lastAdmin = this.isClass('Admin') && (this.Game.Players || []).filter((p) => p.Class === 'Admin').length <= 1;
      const message = lastAdmin
        ? 'Disconnect and end the game for everyone?'
        : 'Disconnect? You will need to join again as a new player.';
      if (window.confirm(message)) {
        this.quitGame();
      }
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
        }).catch(() => {});
      }
      this.resetGame();
    },
    resetGame() {
      // Return to base state
      this.emitter.emit('removeTokensByService', SERVICE);
      this.store.reset();
    },
  },
});
</script>
<style scoped lang="scss">
nav {
  position: relative;
  z-index: 100;
}
.volumeControl {
  width: 8rem;
}
.onHover .btn,
.onHover .volumeControl {
  visibility: hidden;
}
.onHover:hover .btn,
.onHover:hover .volumeControl {
  visibility: visible;
}
</style>
