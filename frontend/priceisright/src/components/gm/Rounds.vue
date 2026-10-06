<template>
  <div>
    <h1>Rounds</h1>
    <form ref="saveRounds"
          class="form-group"
          method="post"
          action="/api/game/rounds/save"
          enctype="multipart/form-data"
          @submit.prevent="saveRounds"
    >
      <ul v-if="game.Rounds && game.Rounds.length"
          class="list-unstyled"
      >
        <li v-for="(round, i) in sortedRounds"
            :key="round.UUID"
            class="row pb-1"
            draggable="true"
        >
          <div class="col-1">
            <FontAwesomeIcon :icon="['fas', 'angle-up']"
                             class="pointer-event text-danger"
                             @click="moveUp(round)"
            />
            <FontAwesomeIcon :icon="['fas', 'angle-down']"
                             class="pointer-event text-primary"
                             @click="moveDown(round)"
            />
            <input v-model="round.UUID" :name="`UUID[${i}]`" type="hidden">
          </div>
          <TextInput class="col-5"
                     input-class="form-control"
                     :value="round.Name"
                     :name="`Name[${i}]`"
                     placeholder="Name"
                     @change="saveRoundName($event, round.UUID)"
          />
          <TextInput class="col-5"
                     input-class="form-control"
                     :value="round.Value"
                     :name="`Value[${i}]`"
                     placeholder="Value"
                     @change="saveRoundValue($event, round.UUID)"
          />
          <div class="col-1">
            <FontAwesomeIcon :icon="['fas', 'xmark']"
                             class="pointer-event text-danger"
                             @click="removeRound(round)"
            />
          </div>
          <div class="col-1"></div>
          <div v-if="round.Image.length" class="col-5">
            <img ref="imagePreview"
                 :src="`${imageURL}/${round.Image}`"
                 class="roundImage"
                 alt=""
            >
          </div>
          <div v-else class="col-5">
            <img ref="imagePreview"
                 class="roundImage"
                 src=""
                 alt=""
            >
          </div>
          <div class="col-5">
            <input class="form-control-file"
                   :name="`imageFile[${i}]`"
                   type="file"
                   @change="updatePreview($event, i)"
            >
          </div>
          <div class="col-1"></div>
        </li>
      </ul>
      <button class="btn btn-primary">Save</button>
    </form>
    <button class="btn btn-danger" @click="addRound">
      Add Round
    </button>
  </div>
</template>

<script>
import Axios from 'axios';
import { v4 as UUIDv4 } from 'uuid';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import TextInput from '../inputs/TextInput.vue';
import useLocalStore from "@/stores/local";

export default {
  name: 'Rounds',
  components: {
    FontAwesomeIcon,
    TextInput,
  },
  props: {
    game: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      ErrorMessage: "",
      localStore: useLocalStore(),
      imageURL: this.imageUri,
    };
  },
  computed: {
    sortedRounds() {
      return [...this.game.Rounds].sort((a, b) => (a.Order < b.Order ? -1 : 1));
    },
  },
  methods: {
    addRound() {
      this.game.Rounds.push({
        Order: this.game.Rounds.length + 1,
        Image: '',
        Name: '',
        Value: 0,
        Guesses: [],
        UUID: UUIDv4(),
      });
    },
    saveRounds() {
      // Clear Options & Votes
      const token = this.getToken();

      const formData = new FormData(this.$refs.saveRounds);

      Axios({
        method: "post",
        url: `${this.apiUri}/game/rounds/save`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
        data: formData,
      }).then(() => {
        this.$emit('roundsUpdated');
      }).catch(() => {
        // Error with reset
      });
    },
    clearGuesses() {
      // Clear Options & Votes
      const token = this.getToken();

      Axios({
        method: "post",
        url: `${this.apiUri}/game/guesses/clear`,
        headers: {
          Authorization: `Bearer ${token.AccessToken}`,
        },
      }).then(() => {
        this.$emit('guessesCleared');
      }).catch(() => {
        // Error with reset
      });
    },
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'PriceIsRight');
    },
    moveDown(round) {
      if (round.Order < this.game.Rounds.length) {
        this.game.Rounds.forEach((r) => {
          if (r.Order === round.Order + 1) {
            r.Order -= 1;
          }
        });
        round.Order += 1;
      }
    },
    moveUp(round) {
      if (round.Order > 1) {
        this.game.Rounds.forEach((r) => {
          if (r.Order === round.Order - 1) {
            r.Order += 1;
          }
        });
        round.Order -= 1;
      }
    },
    removeRound(round) {
      this.game.Rounds = this.game.Rounds.filter((r) => r.UUID !== round.UUID);
    },
    saveRoundName(event, index) {
      const r = this.game.Rounds.find((e) => e.UUID === index);
      r.Name = event.target.value;
    },
    saveRoundValue(event, index) {
      const r = this.game.Rounds.find((e) => e.UUID === index);
      r.Value = event.target.value;
    },
    updatePreview(e, i) {
      this.$refs.imagePreview[i].src = URL.createObjectURL(e.target.files[0]);
    },
  },
};

</script>

<style scoped>
.pointer-event {
  cursor: pointer;
}
.roundImage {
  max-height: 200px;
  max-width: 200px;
}
</style>
