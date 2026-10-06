<template>
  <div>
    <table class="table table-striped w-100">
      <thead>
        <tr>
          <th>
            Player
          </th>
          <th>
            Name
          </th>
          <th>
            Race
          </th>
          <th>
            XP
          </th>
          <th class="text-right">
            Assignment
          </th>
        </tr>
      </thead>
      <template v-for="playerClass in ['LEADER',
                                       'DISPLAY',
                                       'WARRIOR',
                                       'BARD',
                                       'WIZARD',
                                       'GOD',
                                       'AUDIENCE']"
      >
        <tbody v-for="player in game.Players.filter(
                 (p) => p.Class.toUpperCase() === playerClass)
                 .sort((a, b) => a.Name < b.Name ? -1 : a > b ? 1 : 0)"
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
                         :value="player.CharacterName"
                         @change="saveCharacterName($event, player)"
              />
            </td>
            <td>
              <TextInput class="textInput"
                         :value="player.CharacterRace"
                         @change="saveCharacterRace($event, player)"
              />
            </td>
            <td class="user-select-none">
              <FontAwesomeIcon :icon="['fal', 'minus-circle']"
                               @click="adjustExperience(player, -1)"
              />
              {{ player.Experience }}
              <FontAwesomeIcon :icon="['fal', 'plus-circle']"
                               @click="adjustExperience(player, 1)"
              />
            </td>
            <td class="text-right user-select-none">
              <span class="playerClass" title="Audience" @click="assignPlayer(player, 'Audience')">
                <FontAwesomeIcon :icon="['fal', 'eye']"
                                 :class="`${player.Class === 'Audience' ? 'text-danger' : ''}`"
                />
              </span>
              <span class="playerClass" title="God" @click="assignPlayer(player, 'God')">
                <FontAwesomeIcon :icon="['fal', 'bolt']"
                                 :class="`${player.Class === 'God' ? 'text-danger' : ''}`"
                />
              </span>
              <span class="playerClass" title="Warrior" @click="assignPlayer(player, 'Warrior')">
                <FontAwesomeIcon :icon="['fal', 'swords']"
                                 :class="`${player.Class === 'Warrior' ? 'text-danger' : ''}`"
                />
              </span>
              <span class="playerClass" title="Bard" @click="assignPlayer(player, 'Bard')">
                <FontAwesomeIcon :icon="['fal', 'mandolin']"
                                 :class="`${player.Class === 'Bard' ? 'text-danger' : ''}`"
                />
              </span>
              <span class="playerClass" title="Wizard" @click="assignPlayer(player, 'Wizard')">
                <FontAwesomeIcon :icon="['fal', 'wand-magic']"
                                 :class="`${player.Class === 'Wizard' ? 'text-danger' : ''}`"
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
          <tr v-if="['WARRIOR', 'BARD', 'WIZARD'].includes(playerClass)">
            <td colspan="5">
              <div class="d-flex">
                <div class="mr-auto">
                  <span>
                    Equipment:
                  </span>
                  <div v-if="player.Equipment && player.Equipment.length > 0"
                       class="d-inline-block"
                  >
                    <div v-for="(equipment, index) in player.Equipment"
                         :key="`equipment${index}`"
                         class="d-inline-block"
                    >
                      <TextInput class="d-inline-block pt-1 textInput"
                                 :value="equipment"
                                 @change="saveEquipment($event, player, index)"
                      />
                      <FontAwesomeIcon :icon="['fal', 'minus-circle']"
                                       class="m-1 text-danger"
                                       @click="removeEquipment($event, player, index)"
                      />
                    </div>
                  </div>
                  <div class="d-inline">
                    <FontAwesomeIcon :icon="['fal', 'plus-circle']"
                                     class="m-1 text-success"
                                     @click="addEquipment($event, player)"
                    />
                  </div>
                </div>
                <div>
                  <span>Health:</span>
                  <FontAwesomeIcon :icon="['fal', 'skull']"
                                   class="m-1"
                                   @click="setHealth(player, 0)"
                  />
                  <template v-for="h in player.Health" :key="`health${h}`">
                    <FontAwesomeIcon
                                     :icon="['fas', 'heart']"
                                     class="m-1 text-danger"
                                     @click="setHealth(player, h)"
                    />
                  </template>
                  <template v-for="h in 5 - player.Health" :key="`death${h}`">
                    <FontAwesomeIcon
                                     :icon="['fal', 'heart']"
                                     class="m-1 text-dark"
                                     @click="setHealth(player, player.Health + h)"
                    />
                  </template>
                </div>
              </div>
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

export default {
  name: 'PlayerAssignment',
  components: {
    TextInput,
    FontAwesomeIcon,
  },
  props: {
    game: {
      type: Object,
      required: true,
    },
  },
  methods: {
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'Trudge');
    },
    adjustExperience(player, experience) {
      const token = this.getToken();
      const index = this.game.Players.indexOf(player);

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/experience/adjust`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          PlayerUUID: player.UUID,
          Adjustment: experience,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.$emit('playerExperienceChanged', index, player.Experience + experience);
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
    saveCharacterName(event, player) {
      const token = this.getToken();
      const playerIndex = this.game.Players.indexOf(player);

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/name/change/${playerIndex}`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          CharacterName: event.target.value,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.$emit('characterNameChanged', playerIndex, event.target.value, player.UUID);
    },
    saveCharacterRace(event, player) {
      const token = this.getToken();
      const playerIndex = this.game.Players.indexOf(player);

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/race/change/${playerIndex}`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          CharacterRace: event.target.value,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.$emit('characterRaceChanged', playerIndex, event.target.value, player.UUID);
    },
    addEquipment(event, player) {
      const token = this.getToken();
      const playerIndex = this.game.Players.indexOf(player);
      let value = '';

      if (typeof event !== 'undefined'
          && typeof event.target !== 'undefined'
          && typeof event.target.value !== 'undefined') {
        value = event.target.value;
      }

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/equipment/add/${playerIndex}`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Equipment: value,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.$emit('characterEquipmentAdded', playerIndex, value, player.UUID);
    },
    removeEquipment(event, player, index) {
      const token = this.getToken();
      const playerIndex = this.game.Players.indexOf(player);

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/equipment/remove/${playerIndex}/${index}`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      // Don't do this as socket message will cause double removal
      // this.$emit('removeEquipmentRemoved', playerIndex, index, player.UUID);
    },
    saveEquipment(event, player, index) {
      const token = this.getToken();
      const playerIndex = this.game.Players.indexOf(player);

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/equipment/change/${playerIndex}/${index}`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Equipment: event.target.value,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.$emit('characterEquipmentChanged', playerIndex, index, event.target.value, player.UUID);
    },
    setHealth(player, health) {
      const token = this.getToken();
      const index = this.game.Players.indexOf(player);

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/health/adjust`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          PlayerUUID: player.UUID,
          Adjustment: health - player.Health,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.$emit('playerHealthChanged', index, health);
    },
  },
};

</script>

<style scoped lang="scss">
.playerClass, .playerAction {
  cursor: pointer;
}
.textInput {
  height: 1.5rem;
  min-width: 2rem;
}
</style>
