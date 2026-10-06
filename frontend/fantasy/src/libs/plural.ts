// plural("vote", 1) -> "1 vote", plural("vote", 3) -> "3 votes"
export function plural(word: string, count: number): string {
  return `${count} ${count === 1 ? word : `${word}s`}`;
}
