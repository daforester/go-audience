<template>
  <div>
    <div class="text-center">
      <template v-if="Round && Round.UUID">
        <img v-if="Round.Image !== ''"
             :src="`${imageURL}/${Round.Image}`"
             class="roundImage"
             alt=""
        />
        <h2>{{ Round.Name }}</h2>
      </template>
      <template v-if="Guess === null">
        <h3 class="text-center">
          Submit your Guess
        </h3>
        <form class="container" @submit.prevent="sendGuess">
          <fieldset ref="guessForm" class="form-group">
            <div class="row">
              <div class="col-1">
                <span class="form-control-plaintext text-white">
                  &pound;
                </span>
              </div>
              <div class="col-sm-9 col-10">
                <input id="guessInput"
                       v-model="GuessValue"
                       class="form-control"
                       type="number"
                       required="required"
                       placeholder="Guess"
                       aria-label="Guess"
                       step="1"
                       min="0"
                       aria-describedby="guessHelp"
                >
                <small id="guessHelp" class="form-text text-white text-left">
                  Try to guess the price without going over the actual value.
                </small>
              </div>
              <div class="col-sm-2 col-12 text-center text-sm-left mt-2 mt-sm-0">
                <button class="btn btn-primary" type="submit">
                  Submit
                </button>
              </div>
              <small v-if="Errors.Guess"
                     class="form-text text-muted text-danger"
              >
                {{ Errors.Guess }}
              </small>
            </div>
          </fieldset>
        </form>
      </template>
      <template v-else>
        <DisplayMessage effect="toaster"
                        message="Thanks for playing along!"
                        :whole-word="true"
                        class="d-block p-2"
        />
        <p class="text-center p-2">
          If you're close enough you might get called up on to the stage!
        </p>
        <p>
          Now back to the show...
        </p>
      </template>
    </div>
  </div>
</template>

<script>
import Axios from 'axios';
import DisplayMessage from '../display/DisplayMessage.vue';
import useLocalStore from "@/stores/local";

export default {
  name: 'NewGuess',
  components: {
    DisplayMessage,
  },
  data() {
    return {
      Guess: null,
      GuessValue: null,
      Errors: {
        Guess: '',
      },
      Round: {
        Name: '',
        Image: '',
        UUID: '',
      },
      localStore: useLocalStore(),
      imageURL: this.imageUri,
    };
  },
  mounted() {
    this.getRound();
    this.guessExists();
  },
  methods: {
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'PriceIsRight');
    },
    getRound() {
      const token = this.getToken();

      Axios({
        method: "get",
        url: `${this.apiUri}/game/round`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
      }).then((r) => {
        this.Round = r.data.Response.Round;
      }).catch(() => {
        //
      });
    },
    guessExists() {
      const token = this.getToken();

      Axios({
        method: 'get',
        url: `${this.apiUri}/game/guess/exists`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
      }).then((r) => {
        if (r.data.Response.Guess && r.data.Response.Guess.PlayerUUID) {
          this.Guess = r.data.Response.Guess;
        }
      }).catch(() => {
        //
      });
    },
    sendGuess() {
      // Obtain Game Data
      const token = this.getToken();

      const data = {
        Guess: Math.floor(this.GuessValue),
      };

      this.Errors.Guess = '';
      this.$refs.guessForm.setAttribute('disabled', 'disabled');

      Axios({
        method: 'put',
        url: `${this.apiUri}/game/guess/new`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data,
      }).then(() => {
        this.$refs.guessForm.removeAttribute('disabled');
        this.Errors.Guess = '';
        this.Guess = '';
      }).catch((e) => {
        this.$refs.guessForm.removeAttribute('disabled');
        if (typeof e.response === 'object') {
          this.Errors.Guess = e.response.data.ErrorMessage;
        }
      });
    },
  },
};

</script>

<style scoped>
h2, h3 {
  font-family: "LEMON MILK", sans-serif;
}
li {
  cursor: pointer;
  font-size: 2rem;
}
.roundImage {
  max-height: 200px;
  max-width: 200px;
}
</style>
