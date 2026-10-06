import { describe, expect, it } from 'vitest';
import { mount } from '@vue/test-utils';
import TopicCard from '../display/TopicCard.vue';
import type { Submission } from '@/types/game';

const topic = (Text: string): Submission => ({
  ID: 't1',
  Kind: 'Complaint',
  Text,
  Status: 'Played',
  StaffVotes: 0,
  Ups: 0,
  Downs: 0,
  Net: 0,
  Created: '',
});

const mountCard = (text: string, props = {}) =>
  mount(TopicCard, { props: { topic: topic(text), ...props }, global: { stubs: { FontAwesomeIcon: true } } });

describe('TopicCard', () => {
  it('steps the text size down as topics get longer', () => {
    expect(mountCard('Short').classes()).not.toContain('medium');
    expect(mountCard('x'.repeat(120)).classes()).toContain('medium');
    expect(mountCard('x'.repeat(280)).classes()).toContain('long');
  });

  it('collapses long topics only when asked, and expands on request', async () => {
    expect(mountCard('x'.repeat(200)).classes()).not.toContain('collapsed');
    expect(mountCard('Short', { collapsible: true }).find('button').exists()).toBe(false);

    const card = mountCard('x'.repeat(200), { collapsible: true });
    expect(card.classes()).toContain('collapsed');
    await card.find('button').trigger('click');
    expect(card.classes()).not.toContain('collapsed');
    expect(card.find('button').text()).toBe('Show less');
  });
});
