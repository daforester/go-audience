<template>
  <div>
    <form v-if="!Submitted" @submit.prevent="sendBasicCharacterDetails">
      <fieldset v-if="Page === 1">
        <div class="form-group">
          <label for="characterNameInput">
            Name
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
            Domain: <small>(What are you the god of?)</small>
          </label>
          <input id="characterRaceInput"
                 v-model="CharacterRace"
                 class="form-control"
                 type="text"
                 placeholder="Domain"
                 aria-label="Domain"
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
import CharacterDetails from '../party/CharacterDetails.vue';
import PleaseWait from '../player/PleaseWait.vue';

export default {
  name: 'GodDetails',
  components: {
    PleaseWait,
  },
  extends: CharacterDetails,
};

</script>

<style scoped>

</style>
