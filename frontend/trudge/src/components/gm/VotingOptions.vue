<template>
  <div>
    <h1>Voting Options</h1>
    <ul v-if="game.Options && game.Options.length" class="list-unstyled">
      <li v-for="(option, i) in game.Options" :key="i" class="row mb-1">
        <TextInput class="col"
                   input-class="form-control"
                   :value="option"
                   @change="saveOption($event, i)"
        />
        <div class="col p-2">
          ({{ totalVotes(i) }})
        </div>
      </li>
    </ul>
    <button class="btn btn-danger" @click="clearOptions">
      Clear
    </button>
  </div>
</template>

<script>
import Axios from 'axios';
import TextInput from '../inputs/TextInput.vue';

export default {
  name: 'VoteOption',
  components: {
    TextInput,
  },
  props: {
    game: {
      type: Object,
      required: true,
    },
  },
  methods: {
    clearOptions() {
      // Clear Options & Votes
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/option/clear`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
      }).then(() => {
        this.$emit('optionsCleared');
      }).catch(() => {
        // Error with reset
      });
    },
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'Trudge');
    },
    saveOption(event, index) {
      // Change Option
      const token = this.getToken();

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/option/change/${index}`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: {
          Option: event.target.value,
        },
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.ErrorMessage = e.response.data.ErrorMessage;
        }
      });
      this.$emit('optionChanged', index, event.target.value);
    },
    totalVotes(i) {
      return this.game.Votes[i].Total ? this.game.Votes[i].Total : 0;
    },
  },
};

</script>

<style scoped>

</style>
