<template>
  <div>
    <div>
      <h1>Submit your suggestions</h1>
      <form @submit.prevent="sendSuggestion">
        <fieldset ref="suggestionForm">
          <div class="form-group">
            <input id="suggestionInput"
                   v-model="Suggestion"
                   class="form-control"
                   type="text"
                   required="required"
                   placeholder="Suggestion"
                   aria-label="Suggestion"
            >
            <small class="form-text text-muted">
              Suggestions will appear on the main display to inspire the gods.
            </small>
            <small v-if="Errors.Suggestion"
                   class="form-text text-muted text-danger"
            >
              {{ Errors.Suggestion }}
            </small>
          </div>
          <button class="btn btn-primary">
            Submit
          </button>
        </fieldset>
      </form>
    </div>
  </div>
</template>

<script>
import Axios from 'axios';

export default {
  name: 'NewSuggestion',
  data() {
    return {
      Suggestion: '',
      Errors: {
        Suggestion: '',
      },
    };
  },
  methods: {
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'Trudge');
    },
    sendSuggestion() {
      // Obtain Game Data
      const token = this.getToken();

      const data = {
        Suggestion: this.Suggestion,
      };

      this.Errors.Suggestion = '';
      this.$refs.suggestionForm.setAttribute('disabled', 'disabled');

      Axios({
        method: 'put',
        url: `${this.apiUri}/game/suggestion/new`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data,
      }).then(() => {
        this.$refs.suggestionForm.removeAttribute('disabled');
        this.Errors.Suggestion = '';
        this.Suggestion = '';
      }).catch((e) => {
        this.$refs.suggestionForm.removeAttribute('disabled');
        if (typeof e.response === 'object') {
          this.Errors.Suggestion = e.response.data.ErrorMessage;
        }
      });
    },
  },
};

</script>

<style scoped lang="scss">
li {
  cursor: pointer;
  font-size: 2rem;
}
</style>
