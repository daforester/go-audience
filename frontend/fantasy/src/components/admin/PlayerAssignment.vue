<template>
  <div class="card-panel mb-3">
    <h2 class="h5">
      Players
      <small class="text-muted">{{ players.length }}</small>
    </h2>
    <div class="table-responsive">
      <table class="table table-sm align-middle mb-0">
        <thead>
          <tr>
            <th>Name</th>
            <th>Role</th>
            <th>Team</th>
            <th />
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in players" :key="p.UUID">
            <td>
              {{ p.Name }}
              <small v-if="p.UUID === player.UUID" class="text-muted">(you)</small>
            </td>
            <td>
              <select class="form-select form-select-sm"
                      :value="p.Class"
                      :aria-label="`Role for ${p.Name}`"
                      @change="assign(p.UUID, ($event.target as HTMLSelectElement).value, p.TeamID)"
              >
                <option v-for="(label, value) in Roles" :key="value" :value="value">
                  {{ label }}
                </option>
              </select>
            </td>
            <td>
              <select v-if="p.Class !== 'Admin' && game.Teams.length > 0"
                      class="form-select form-select-sm"
                      :value="p.Class === 'Director' ? p.TeamID : ''"
                      :aria-label="`Team for ${p.Name}`"
                      @change="assign(p.UUID, 'Director', ($event.target as HTMLSelectElement).value)"
              >
                <option value="" disabled>
                  Make director of...
                </option>
                <option v-for="team in game.Teams" :key="team.ID" :value="team.ID">
                  {{ team.Name }}
                </option>
              </select>
            </td>
            <td class="text-end">
              <button v-if="p.UUID !== player.UUID"
                      class="btn btn-outline-danger btn-sm"
                      :title="`Remove ${p.Name}`"
                      @click="kick(p.UUID, p.Name)"
              >
                <FontAwesomeIcon :icon="['fas', 'xmark']" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-if="ErrorMessage" class="alert alert-danger mt-2 mb-0 py-1">
      {{ ErrorMessage }}
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import PlayerView from '@/views/Player';
import type { Player } from '@/types/game';

const ClassOrder = ['Admin', 'Display', 'Staff', 'Director', 'Audience'];

export default defineComponent({
  name: 'PlayerAssignment',
  components: {
    FontAwesomeIcon,
  },
  mixins: [PlayerView],
  data() {
    return {
      Roles: {
        Admin: 'Admin',
        Display: 'Main Display',
        Staff: 'Staff Team',
        Director: 'Con Director',
        Audience: 'Audience',
      },
    };
  },
  computed: {
    players(): Player[] {
      return [...(this.game.Players || [])].sort(
        (a, b) => ClassOrder.indexOf(a.Class) - ClassOrder.indexOf(b.Class) || a.Name.localeCompare(b.Name),
      );
    },
  },
  methods: {
    assign(uuid: string, playerClass: string, teamID: string) {
      if (playerClass === 'Director' && !teamID) {
        // Pick the team with the fewest directors so far
        const team = [...this.game.Teams].sort((a, b) => a.Members.length - b.Members.length)[0];
        if (!team) {
          this.ErrorMessage = 'Create a director team first';
          return;
        }
        teamID = team.ID;
      }
      this.gameRequest('post', '/game/role/assign', {
        UUID: uuid,
        PlayerClass: playerClass,
        TeamID: teamID,
      }).catch(() => {});
    },
    kick(uuid: string, name: string) {
      if (window.confirm(`Remove ${name} from the game?`)) {
        this.gameRequest('get', `/game/kick/${uuid}`).catch(() => {});
      }
    },
  },
});
</script>

<style scoped lang="scss">
select {
  min-width: 9.5rem;
}
</style>
