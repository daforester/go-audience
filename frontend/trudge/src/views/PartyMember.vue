<template>
  <div>
    <Sunburst v-if="[30, 40 ].includes(gameCode)" />
    <IconDance v-else-if="[90, 100, 120, 130, 150, 160, 180, 190, 550, 560 ].includes(gameCode)" />
    <Travel v-else />
    <div class="d-flex vh-100 vw-100 main-display position-absolute overflow-hidden">
      <div class="m-auto">
        <CharacterDetails v-if="gameCode >= 30 && gameCode <= 40"
                          :player="player"
                          :game-code="gameCode"
        />
        <VoteWeapon v-else-if="gameCode >= 100 && gameCode <= 100"
                    :player="player"
                    :game-code="gameCode"
        />
        <VoteSong v-else-if="gameCode >= 130 && gameCode <= 130"
                  :player="player"
                  :game-code="gameCode"
        />
        <VoteSpell v-else-if="gameCode >= 160 && gameCode <= 160"
                   :player="player"
                   :game-code="gameCode"
        />
        <VoteQuest v-else-if="gameCode >= 190 && gameCode <= 190"
                   :player="player"
                   :game-code="gameCode"
        />
        <VoteOption v-else-if="gameCode >= 560 && gameCode <= 560"
                    :player="player"
                    :game-code="gameCode"
        />
        <DiceRoller v-else-if="gameCode >= 510 && gameCode <= 510"
                    :player="player"
                    :game-code="gameCode"
        />
        <div v-else>
          <PleaseWait :player="player"
                      :game-code="gameCode"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import PlayerView from './Player.vue';
import CharacterDetails from '../components/party/CharacterDetails.vue';
import DiceRoller from '../components/party/DiceRoller.vue';
import VoteWeapon from '../components/audience/VoteWeapon.vue';
import VoteSong from '../components/audience/VoteSong.vue';
import VoteSpell from '../components/audience/VoteSpell.vue';
import VoteQuest from '../components/audience/VoteQuest.vue';
import VoteOption from '../components/audience/VoteOption.vue';
import PleaseWait from '../components/player/PleaseWait.vue';
import Sunburst from '../components/backgrounds/Sunburst.vue';
import IconDance from '../components/backgrounds/IconDance.vue';
import Travel from '../components/backgrounds/Travel.vue';

export default {
  name: 'PartyMemberView',
  components: {
    PleaseWait,
    CharacterDetails,
    IconDance,
    DiceRoller,
    Sunburst,
    Travel,
    VoteOption,
    VoteQuest,
    VoteSpell,
    VoteSong,
    VoteWeapon,
  },
  extends: PlayerView,
};

</script>

<style scoped lang="scss">
.main-display {
  top: 0;
  left: 0;
}
</style>
