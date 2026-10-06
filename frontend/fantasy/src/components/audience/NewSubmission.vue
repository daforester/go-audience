<template>
  <div>
    <h1 class="h3 text-center">
      Tell the con directors
    </h1>
    <form @submit.prevent="sendSubmission">
      <fieldset :disabled="requesting">
        <div class="kinds mb-2" role="radiogroup" aria-label="Kind of submission">
          <template v-for="kind in game.SubmissionKinds" :key="kind">
            <input :id="`kind-${kind}`"
                   v-model="Kind"
                   type="radio"
                   class="btn-check"
                   name="kind"
                   :value="kind"
            >
            <label class="btn btn-outline-primary" :for="`kind-${kind}`">
              <FontAwesomeIcon :icon="['fas', KindIcons[kind]]" />
              {{ SubmissionKindLabels[kind] }}
            </label>
          </template>
        </div>
        <div class="mb-2">
          <textarea v-model="Text"
                    class="form-control"
                    rows="3"
                    :maxlength="maxLength"
                    required
                    :placeholder="KindPrompts[Kind]"
                    aria-label="Your submission"
          />
          <div class="form-text d-flex">
            <span class="me-auto">The staff team pick the best for the directors to deal with.</span>
            <span class="text-nowrap ms-2">{{ Text.length }} / {{ maxLength }}</span>
          </div>
          <div v-if="ErrorMessage" class="form-text text-danger">
            {{ ErrorMessage }}
          </div>
          <div v-if="Sent" class="form-text text-success">
            Sent! Send another whenever you like.
          </div>
        </div>
        <button class="btn btn-primary big-button">
          Submit
        </button>
      </fieldset>
    </form>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import PlayerView from '@/views/Player';
import { KindIcons, KindPrompts } from '@/libs/kinds';
import { SubmissionKindLabels, type Submission, type SubmissionKind } from '@/types/game';

export default defineComponent({
  name: 'NewSubmission',
  components: {
    FontAwesomeIcon,
  },
  mixins: [PlayerView],
  data() {
    return {
      Kind: 'Complaint' as SubmissionKind,
      Text: '',
      Sent: false,
      maxLength: 280,
      KindIcons,
      KindPrompts,
      SubmissionKindLabels,
    };
  },
  methods: {
    sendSubmission() {
      this.Sent = false;
      this.gameRequest<{ Submission: Submission }>('put', '/game/submission/new', {
        Kind: this.Kind,
        Text: this.Text,
      })
        .then((r) => {
          this.store.addMySubmission(r.Submission);
          this.Text = '';
          this.Sent = true;
        })
        .catch(() => {});
    },
  },
});
</script>

<style scoped lang="scss">
.kinds {
  display: grid;
  gap: 0.5rem;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  .btn {
    align-items: center;
    display: flex;
    gap: 0.35rem;
    justify-content: center;
    line-height: 1.2;
  }
}
</style>
