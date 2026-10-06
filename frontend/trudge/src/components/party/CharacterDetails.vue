<template>
  <div>
    <form v-if="!Submitted" @submit.prevent="sendBasicCharacterDetails">
      <fieldset v-if="Page === 1">
        <div class="form-group">
          <label for="characterNameInput">
            {{ player.Class }} Name
          </label>
          <input id="characterNameInput"
                 v-model="CharacterName"
                 class="form-control"
                 type="text"
                 required="required"
                 placeholder="Name"
                 aria-label="Name"
          >
          <small v-if="Errors.CharacterName"
                 class="form-text text-muted"
          >
            {{ Errors.CharacterName }}
          </small>
        </div>
        <div class="form-group">
          <label for="characterRaceInput">
            {{ player.Class }} Race
          </label>
          <input id="characterRaceInput"
                 v-model="CharacterRace"
                 class="form-control"
                 type="text"
                 placeholder="Race"
                 aria-label="Race"
                 required
          >
          <small v-if="Errors.CharacterRace"
                 class="form-text text-muted"
          >
            {{ Errors.CharacterRace }}
          </small>
        </div>
        <button class="btn btn-primary">
          Submit
        </button>
      </fieldset>
    </form>
    <div v-else>
      <PleaseWait :player="player" :game-code="gameCode" />
    </div>
  </div>
</template>

<script>
import Axios from 'axios';
import PleaseWait from '../player/PleaseWait.vue';

export default {
  name: 'CharacterDetails',
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
      CharacterName: '',
      CharacterRace: '',
      Errors: {
        CharacterName: '',
        CharacterRace: '',
      },
      Page: 1,
      Submitted: false,
    };
  },
  methods: {
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'Trudge');
    },
    sendBasicCharacterDetails() {
      this.Submitted = true;

      // Obtain Game Data
      const token = this.getToken();

      this.Errors.CharacterName = '';
      this.Errors.CharacterRace = '';

      let data = {
        CharacterName: this.CharacterName,
      };

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/name/assign`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data,
      }).catch((e) => {
        this.Submitted = false;
        if (typeof e.response === 'object') {
          this.Errors.CharacterName = e.response.data.ErrorMessage;
        }
      });

      data = {
        CharacterRace: this.CharacterRace,
      };

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/race/assign`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data,
      }).catch((e) => {
        this.Submitted = false;
        if (typeof e.response === 'object') {
          this.Errors.CharacterRace = e.response.data.ErrorMessage;
        }
      });
    },
  },
};

</script>

<style scoped>

</style>
