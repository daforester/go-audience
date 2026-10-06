<template>
  <div class="py-4">
    <h1 class="text-center mb-4 title">
      <span class="d-block small text-uppercase">The</span>
      Fantasy Furry Convention
      <FontAwesomeIcon :icon="['fas', 'paw']" class="d-block mx-auto mt-2 dance" />
    </h1>
    <form class="mx-auto mb-4" @submit.prevent="connectGame">
      <fieldset class="border border-dark bg-light rounded-3 p-3">
        <legend class="visually-hidden">
          Connect to Game
        </legend>
        <div class="mb-2">
          <input v-model="RoomCode"
                 class="form-control text-uppercase"
                 placeholder="Room Code"
                 aria-label="Room Code"
                 maxlength="5"
                 required
          />
        </div>
        <div class="mb-2">
          <input v-model="Name"
                 class="form-control"
                 placeholder="Your Name"
                 aria-label="Your Name"
                 maxlength="32"
                 required
          >
        </div>
        <div class="text-center">
          <button class="btn btn-primary">
            Join the Convention
          </button>
        </div>
        <div v-if="Errors.Connect" class="mt-2 form-text text-danger text-center">
          {{ Errors.Connect }}
        </div>
      </fieldset>
    </form>
    <form class="mx-auto mb-4" @submit.prevent="createGame">
      <fieldset class="border border-dark bg-light rounded-3 p-3">
        <legend class="visually-hidden">
          Create Game
        </legend>
        <div class="mb-2">
          <input v-model="Name"
                 class="form-control"
                 placeholder="Your Name"
                 aria-label="Your Name"
                 maxlength="32"
                 required
          >
        </div>
        <div class="text-center">
          <button class="btn btn-secondary">
            Run a Convention
          </button>
        </div>
        <div v-if="Errors.Create" class="mt-2 form-text text-danger text-center">
          {{ Errors.Create }}
        </div>
      </fieldset>
    </form>
  </div>
</template>
<script lang="ts">
import { defineComponent } from 'vue';
import Axios from 'axios';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { SERVICE } from '@/mixins/api/GameRequest';
import type { Tokens } from '@/types/game';

export default defineComponent({
  name: 'LoginView',
  components: {
    FontAwesomeIcon,
  },
  emits: ['joinedGame', 'createdGame'],
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
    if (urlParams.has('code')) {
      this.RoomCode = urlParams.get('code') || '';
    }
    if (urlParams.has('action')) {
      if (!urlParams.has('name')) {
        return;
      }
      this.Name = urlParams.get('name') || '';
      switch (urlParams.get('action')) {
        case 'join':
          if (!this.RoomCode) {
            return;
          }
          this.connectGame();
          break;
        case 'create':
          this.createGame();
          break;
      }
    }
  },
  methods: {
    storeToken(token: Tokens) {
      token.Endpoint = `${this.apiUri}/token/refresh`;
      token.Service = SERVICE;
      token.Vendor = this.RoomCode;

      this.emitter.emit('newToken', token);
    },
    errorMessage(e: unknown): string {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const response = (e as any)?.response;
      if (typeof response === 'object' && response?.data?.ErrorMessage) {
        return response.data.ErrorMessage;
      }
      return 'Could not reach the convention, try again';
    },
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
      })
        .then((r) => {
          this.storeToken(r.data.Response.Tokens);
          this.$emit('joinedGame', r.data.Response.Player);
        })
        .catch((e) => {
          this.Errors.Connect = this.errorMessage(e);
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
      })
        .then((r) => {
          this.storeToken(r.data.Response.Tokens);
          this.$emit('createdGame', r.data.Response.Game, r.data.Response.Player);
        })
        .catch((e) => {
          this.Errors.Create = this.errorMessage(e);
        });
    },
  },
});
</script>
<style scoped lang="scss">
@use '@/libs/scss/animations/dance.scss';

form {
  max-width: 280px;
}
.title {
  color: #7b3fbf;
  font-weight: 800;
}
</style>
