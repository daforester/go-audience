import { describe, expect, it, beforeEach } from 'vitest';
import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import TeamBars from '../display/TeamBars.vue';
import type { Team } from '@/types/game';

const team = (ID: string, Name: string): Team => ({ ID, Name, RoundWins: 0, TotalVotes: 0, Members: [] });

describe('TeamBars', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  for (const count of [2, 3]) {
    it(`lays out ${count} teams`, () => {
      const teams = [team('a', 'Alpha'), team('b', 'Bravo'), team('c', 'Charlie')].slice(0, count);
      const votes = { a: 3, b: 1, c: 0 };
      const wrapper = mount(TeamBars, {
        props: { teams, votes, winners: ['a'] },
        global: { stubs: { FontAwesomeIcon: true } },
      });

      expect(wrapper.findAll('.team')).toHaveLength(count);
      expect(wrapper.find('.team-bars').attributes('style')).toContain(`--teams: ${count}`);
      expect(wrapper.findAll('.winner')).toHaveLength(1);
      expect(wrapper.find('.bar').attributes('style')).toContain('height: 75%');
    });
  }

  it('shows empty bars with no votes', () => {
    const wrapper = mount(TeamBars, {
      props: { teams: [team('a', 'Alpha'), team('b', 'Bravo')], votes: {} },
      global: { stubs: { FontAwesomeIcon: true } },
    });
    expect(wrapper.find('.bar').attributes('style')).toContain('height: 0%');
  });
});
