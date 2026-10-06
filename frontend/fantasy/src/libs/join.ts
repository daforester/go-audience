// Where players join: the page the display is served from, with the room code filled in
export function joinUrl(roomKey: string): string {
  const url = new URL('/', window.location.origin);
  url.searchParams.set('code', roomKey);
  return url.href;
}

export function joinHost(): string {
  return window.location.host;
}
