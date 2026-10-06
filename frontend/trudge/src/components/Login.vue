<template>
  <div>
    <h1 class="text-center h4 text-uppercase">
      <img class="titleBanner"
           src="@/assets/images/title-banner.svg"
           alt="Trudge"
           title="Trudge"
      />
    </h1>
    <form class="ml-auto mr-auto mb-4" @submit.prevent="connectGame">
      <fieldset class="border-dark border bg-light rounded-lg p-2">
        <legend class="text-center">
          Connect to Game
        </legend>
        <div class="form-row">
          <div class="form-group col">
            <input v-model="RoomCode"
                   class="form-control text-uppercase"
                   placeholder="Room Code"
                   aria-label="Room Code"
                   required
            />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group col">
            <input v-model="Name"
                   class="form-control"
                   placeholder="Your Name"
                   aria-label="Your Name"
                   required
            >
          </div>
        </div>
        <div class="form-row">
          <div class="form-group col text-center mb-0">
            <button class="btn btn-primary">
              Connect to Game
            </button>
          </div>
        </div>
        <div v-if="Errors.Connect" class="mt-2">
          <span class="form-text text-danger text-center text-capitalize">
            {{ Errors.Connect }}
          </span>
        </div>
      </fieldset>
    </form>
    <form class="ml-auto mr-auto mb-4" @submit.prevent="createGame">
      <fieldset class="border-dark border bg-light rounded-lg p-2">
        <legend class="text-center">
          Create Game
        </legend>
        <div class="form-row">
          <div class="form-group col">
            <input v-model="Name"
                   class="form-control"
                   placeholder="Your Name"
                   aria-label="Your Name"
                   required
            >
          </div>
        </div>
        <div class="form-row">
          <div class="form-group col text-center mb-0">
            <button class="btn btn-secondary">
              Create Game
            </button>
          </div>
        </div>
        <div v-if="Errors.Create" class="mt-2">
          <span class="form-text text-danger text-center text-capitalize">
            {{ Errors.Create }}
          </span>
        </div>
      </fieldset>
    </form>
  </div>
</template>
<script>
import Axios from 'axios';

export default {
  name: 'LoginView',
  data() {
    return {
      Errors: {
        Connect: '',
        Create: '',
      },
      Name: '',
      RoomCode: '',
    };
  },
  mounted() {
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.has('action')) {
      if (!urlParams.has('name')) {
        return;
      }
      this.Name = urlParams.get('name');
      switch (urlParams.get('action')) {
        case 'join':
          if (!urlParams.has('code')) {
            return;
          }
          this.RoomCode = urlParams.get('code');
          this.connectGame();
          break;
        case 'create':
          this.createGame();
          break;
      }
    }
  },
  methods: {
    connectGame() {
      this.Errors.Connect = '';
      const data = {
        RoomCode: this.RoomCode,
        Name: this.Name,
      };
      Axios({
        method: 'post',
        url: `${this.apiUri}/game/join`,
        data,
      }).then((r) => {
        const player = r.data.Response.Player;
        const token = r.data.Response.Tokens;

        token.Endpoint = `${this.apiUri}/token/refresh`;
        token.Service = 'Trudge';
        token.Vendor = this.RoomCode;

        this.emitter.emit('newToken', token);
        this.$emit('joinedGame', player);
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.Errors.Connect = e.response.data.ErrorMessage;
        }
      });
    },
    createGame() {
      this.Errors.Create = '';
      const data = {
        Name: this.Name,
      };
      Axios({
        method: 'post',
        url: `${this.apiUri}/game/create`,
        data,
      }).then((r) => {
        const game = r.data.Response.Game;
        const player = r.data.Response.Player;
        const token = r.data.Response.Tokens;

        token.Endpoint = `${this.apiUri}/token/refresh`;
        token.Service = 'Trudge';
        token.Vendor = this.RoomCode;

        this.emitter.emit('newToken', token);
        this.$emit('createdGame', game, player);
      }).catch((e) => {
        if (typeof e.response === 'object') {
          this.Errors.Create = e.response.data.ErrorMessage;
        }
      });
    },
  },
};

</script>
<style scoped>
  form {
    max-width: 240px;
  }
  legend {
    position: absolute;
    left: -10000px;
  }
  .titleBanner {
    max-height: calc(25vh);
    max-width: 90%;
    width: 90%;
  }
</style>
