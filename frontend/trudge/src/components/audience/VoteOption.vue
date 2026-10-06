<template>
  <div>
    <div v-if="!Submitted">
      <h1>{{ RequestLabel }}</h1>
      <ul class="list-unstyled">
        <li v-for="(option, i) in Options"
            :key="i"
            class="border border-dark rounded rounded-lg m-2 p-2 text-center"
            @click.prevent="voteOption(i)"
        >
          {{ option }}
        </li>
      </ul>
    </div>
    <div v-else>
      <PleaseWait :player="player" :game-code="gameCode" />
    </div>
  </div>
</template>

<script>
import Axios from 'axios';
import PleaseWait from '../player/PleaseWait.vue';

export default {
  name: 'VoteOption',
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
      Options: [],
      RequestLabel: 'Choose a fate',
      Submitted: false,
    };
  },
  created() {
    this.getOptions();
  },
  methods: {
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'Trudge');
    },
    getOptions() {
      // Obtain Game Data
      const token = this.getToken();

      Axios({
        method: 'get',
        url: `${this.apiUri}/game/options`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
      }).then((r) => {
        this.Options = r.data.Response.Options;
      }).catch((e) => {
        console.log(e);
      });
    },
    voteOption(index) {
      // Obtain Game Data
      const token = this.getToken();

      const data = {
        Index: index,
        Option: this.Options[index],
      };

      Axios({
        method: 'post',
        url: `${this.apiUri}/game/option/vote`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data,
      }).then(() => {
        this.Submitted = true;
      }).catch((e) => {
        this.Submitted = false;
        if (typeof e.response === 'object') {
          this.Errors.CharacterName = e.response.data.ErrorMessage;
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
