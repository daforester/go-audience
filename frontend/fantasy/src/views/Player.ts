import { defineComponent } from 'vue';
import GameRequest from '@/mixins/api/GameRequest';
import useGameStore from '@/stores/game';
import { GameCodes, type Game, type Player, type Round, type Team } from '@/types/game';

// PlayerView is the base of every role view, as Trudge's Player.vue was
export default defineComponent({
  mixins: [GameRequest],
  props: {
    gameCode: {
      type: Number,
      required: true,
    },
  },
  data() {
    return {
      GameCodes,
      store: useGameStore(),
    };
  },
  computed: {
    game(): Game {
      return this.store.Game;
    },
    player(): Player {
      return this.store.Player;
    },
    round(): Round | undefined {
      return this.store.Game.Round;
    },
    roundTeams(): Team[] {
      const round = this.store.Game.Round;
      if (!round) {
        return this.store.Game.Teams;
      }
      return round.TeamOrder.map((id) => this.store.teamById[id]).filter((t): t is Team => !!t);
    },
  },
});
