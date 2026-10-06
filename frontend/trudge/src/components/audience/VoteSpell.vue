<script>

import VoteOption from './VoteOption.vue';

export default {
  name: 'VoteSpell',
  extends: VoteOption,
  data() {
    return {
      RequestLabel: 'Choose a spell',
    };
  },
};

</script>
