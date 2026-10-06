<template>
  <div>
    <form v-if="!Submitted" @submit.prevent="sendOption">
      <fieldset v-if="Page === 1">
        <div class="form-group">
          <label for="characterNameInput">
            {{ RequestLabel }}
          </label>
          <input id="characterNameInput"
                 v-model="Option"
                 class="form-control"
                 type="text"
                 required="required"
                 placeholder="Name"
                 aria-label="Name"
          >
          <small v-if="Errors.Option"
                 class="form-text text-muted"
          >
            {{ Errors.Option }}
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
  name: 'RequestOption',
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
      Option: '',
      Errors: {
        Option: '',
      },
      Page: 1,
      RequestLabel: 'Suggest a choice',
      Submitted: false,
    };
  },
  methods: {
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'Trudge');
    },
    sendOption() {
      this.Submitted = true;

      // Send New Option
      const token = this.getToken();

      this.Errors.Option = '';

      const data = {
        Option: this.Option,
      };

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/option/add`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data,
      }).catch((e) => {
        this.Submitted = false;
        if (typeof e.response === 'object') {
          this.Errors.Option = e.response.data.ErrorMessage;
        }
      });
    },
  },
};

</script>

<style scoped>

</style>
