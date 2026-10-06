// API endpoints always come from build configuration (VITE_API_URL / VITE_WS_URL),
// never from code. vite.config.ts refuses to build without them.

export interface Endpoints {
  apiUri: string;
  wsUri: string;
}

const API_PROTOCOLS = ['http:', 'https:'];
const WS_PROTOCOLS = ['ws:', 'wss:'];

// endpointErrors explains what is wrong with the configured URLs, if anything
export function endpointErrors(apiUrl?: string, wsUrl?: string): string[] {
  const errors: string[] = [];
  const check = (name: string, value: string | undefined, protocols: string[]) => {
    if (!value) {
      errors.push(`${name} is not set`);
      return;
    }
    let url: URL;
    try {
      url = new URL(value);
    } catch {
      errors.push(`${name} must be an absolute URL, got "${value}"`);
      return;
    }
    if (!protocols.includes(url.protocol)) {
      errors.push(`${name} must use ${protocols.join(' or ')}, got "${value}"`);
    }
  };
  check('VITE_API_URL', apiUrl, API_PROTOCOLS);
  check('VITE_WS_URL', wsUrl, WS_PROTOCOLS);
  return errors;
}

export function resolveEndpoints(apiUrl?: string, wsUrl?: string): Endpoints {
  const errors = endpointErrors(apiUrl, wsUrl);
  if (errors.length > 0) {
    throw new Error(`Fantasy endpoints misconfigured: ${errors.join('; ')}`);
  }
  return {
    apiUri: (apiUrl as string).replace(/\/+$/, ''),
    wsUri: (wsUrl as string).replace(/\/+$/, ''),
  };
}
