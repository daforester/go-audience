<script>

import VoteOption from './VoteOption.vue';

export default {
  name: 'VoteWeapon',
  extends: VoteOption,
  data() {
    return {
      RequestLabel: 'Choose a weapon',
    };
  },
};

</script>
