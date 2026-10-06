import { afterEach, describe, expect, it } from 'vitest';
import { mount } from '@vue/test-utils';
import mitt from 'mitt';
import Login from '../Login.vue';

describe('Login', () => {
  afterEach(() => {
    window.history.replaceState({}, '', '/');
  });

  it('fills in the room code from a scanned join link', async () => {
    window.history.replaceState({}, '', '/?code=47VSS');
    const wrapper = mount(Login, {
      global: {
        stubs: { FontAwesomeIcon: true },
        config: { globalProperties: { emitter: mitt(), apiUri: 'http://api.test/api' } as never },
      },
    });
    await wrapper.vm.$nextTick();
    const roomCode = wrapper.find('input[aria-label="Room Code"]').element as HTMLInputElement;
    expect(roomCode.value).toBe('47VSS');
  });
});
