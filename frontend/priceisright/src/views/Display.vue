<template>
  <div class="vh-100 d-flex align-items-center justify-content-center">
    <div v-if='gameCode === 40 && calling.UUID !== ""'>
      <DisplayMessage :message="calling.Name" effect="comedown" :whole-word="true" />
      <DisplayMessage :message="calling.BadgeNumber"
                      effect="comedown"
                      :whole-word="true"
                      class="smaller"
      />
    </div>
  </div>
</template>

<script>
import Axios from 'axios';
import PlayerView from './Player.vue';
import DisplayMessage from '../components/display/DisplayMessage.vue';

export default {
  name: 'DisplayView',
  components: {
    DisplayMessage,
  },
  extends: PlayerView,
  watch: {
    gameCode(v) {
      this.codeChanged(v);
    },
  },
  created() {
    this.codeChanged(this.gameCode);
  },
  methods: {
    codeChanged(v) {
      if (v === 40) {
        // Update game because of results
        this.refreshGame();
      }
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
        this.top = r.data.Guesses;
      }).catch(() => {
        // Error with reset
      });
    },
  },
};

</script>
<style>
  body {
    background-color: rgba(0,0,0,0) !important;
  }
</style>
