import { describe, expect, it } from 'vitest';
import { joinUrl } from '../join';

describe('joinUrl', () => {
  it('points at the site the display is on, with the room code filled in', () => {
    expect(joinUrl('47VSS')).toBe(`${window.location.origin}/?code=47VSS`);
  });
});
