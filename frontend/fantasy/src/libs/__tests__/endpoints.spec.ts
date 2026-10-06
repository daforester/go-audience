import { describe, expect, it } from 'vitest';
import { endpointErrors, resolveEndpoints } from '../endpoints';

describe('endpoints', () => {
  it('accepts absolute http(s) API and ws(s) socket URLs', () => {
    expect(endpointErrors('https://ffc-api.furrific.uk/api', 'wss://ffc-api.furrific.uk/ws')).toEqual([]);
    expect(resolveEndpoints('http://localhost:18090/api/', 'ws://localhost:18090/ws/')).toEqual({
      apiUri: 'http://localhost:18090/api',
      wsUri: 'ws://localhost:18090/ws',
    });
  });

  it('refuses missing, relative or wrong-protocol URLs', () => {
    expect(endpointErrors(undefined, undefined)).toHaveLength(2);
    expect(endpointErrors('/api', 'wss://x/ws')[0]).toContain('absolute');
    expect(endpointErrors('ftp://x/api', 'https://x/ws')).toHaveLength(2);
    expect(() => resolveEndpoints('', '')).toThrow('misconfigured');
  });
});
