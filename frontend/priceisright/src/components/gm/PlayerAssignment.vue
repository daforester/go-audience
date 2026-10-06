<template>
  <div>
    <table class="table table-striped w-100">
      <thead>
        <tr>
          <th>
            Player
          </th>
          <th>
            Badge
          </th>
          <th class="text-right">
            Assignment
          </th>
        </tr>
      </thead>
      <template v-for="playerClass in ['LEADER',
                                       'DISPLAY',
                                       'AUDIENCE']"
      >
        <tbody v-for="player in game.Players.filter((p) => p.Class.toUpperCase() === playerClass)"
               :key="`${player.UUID}`"
        >
          <tr>
            <td>
              <TextInput class="textInput"
                         :value="player.Name"
                         @change="savePlayerName($event, player)"
              />
            </td>
            <td>
              <TextInput class="textInput"
                         :value="player.BadgeNumber"
                         @change="saveBadgeNumber($event, player)"
              />
            </td>
            <td class="text-right user-select-none">
              <span class="playerClass" title="Audience" @click="assignPlayer(player, 'Audience')">
                <FontAwesomeIcon :icon="['fal', 'eye']"
                                 :class="`${player.Class === 'Audience' ? 'text-danger' : ''}`"
                />
              </span>
              <span class="playerClass" title="Display" @click="assignPlayer(player, 'Display')">
                <FontAwesomeIcon :icon="['fal', 'tv']"
                                 :class="`${player.Class === 'Display' ? 'text-danger' : ''}`"
                />
              </span>
              <span class="playerClass" title="DM" @click="assignPlayer(player, 'Leader')">
                <FontAwesomeIcon :icon="['fal', 'scroll-old']"
                                 :class="`${player.Class === 'Leader' ? 'text-danger' : ''}`"
                />
              </span>
              <span class="playerAction" title="Kick Player" @click="kickPlayer(player.UUID)">
                <FontAwesomeIcon :icon="['fas', 'unlink']" />
              </span>
            </td>
          </tr>
        </tbody>
      </template>
    </table>
  </div>
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import Axios from 'axios';
import TextInput from '../inputs/TextInput.vue';
import useLocalStore from "@/stores/local";

export default {
  name: 'PlayerAssignment',
  components: {
    TextInput,
    FontAwesomeIcon,
  },
  data() {
    return {
      localStore: useLocalStore(),
    };
  },
  props: {
    game: {
      type: Object,
      required: true,
    },
  },
  methods: {
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'PriceIsRight');
    },
    assignPlayer(player, role) {
      const token = this.getToken();
      const data = {
        UUID: player.UUID,
        PlayerClass: role,
      };
      Axios({
        method: 'post',
        url: `${this.apiUri}/game/role/assign`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data,
      });
    },
    kickPlayer(uuid) {
      const token = this.getToken();
      Axios({
        method: 'get',
        url: `${this.apiUri}/game/kick/${uuid}`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
      });
    },
    savePlayerName(event, player) {
      const token = this.getToken();
      const playerIndex = this.game.Players.indexOf(player);

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/player/change/${playerIndex}`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Name: event.target.value,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.$emit('playerNameChanged', playerIndex, event.target.value);
    },
    saveBadgeNumber(event, player) {
      const token = this.getToken();
      const playerIndex = this.game.Players.indexOf(player);

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/player/badge/change/${playerIndex}`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          BadgeNumber: event.target.value,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.$emit('badgeNumberChanged', playerIndex, event.target.value, player.UUID);
    },
  },
};

</script>

<style scoped>
.playerClass, .playerAction {
  cursor: pointer;
}
.textInput {
  height: 1.5rem;
  min-width: 2rem;
}
</style>
