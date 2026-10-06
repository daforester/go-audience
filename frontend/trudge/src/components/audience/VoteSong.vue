<script>

import VoteOption from './VoteOption.vue';

export default {
  name: 'VoteSong',
  extends: VoteOption,
  data() {
    return {
      RequestLabel: 'Choose a song',
    };
  },
};

</script>
