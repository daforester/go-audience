import { GameCodes } from '@/types/game';

export interface Track {
  attribution: boolean;
  copyright: string;
  file: string;
  name: string;
}

// Tracks the main display can play, as Trudge's Display.vue did. None are bundled yet:
// add files to src/assets/music, import them here and list them, for example
//
//   import musicTavern from '@/assets/music/177_Tavern_Music.mp3';
//   { attribution: true, copyright: 'tabletopaudio.com (CC BY-NC-ND 4.0)', file: musicTavern, name: 'Tavern Music' },
export const Music: Track[] = [];

// Which track plays for each game code, by index into Music
export const TrackForCode: Partial<Record<number, number>> = {
  [GameCodes.Lobby]: 0,
  [GameCodes.Intake]: 0,
  [GameCodes.TopicReveal]: 1,
  [GameCodes.AudienceVote]: 2,
  [GameCodes.FinalTally]: 3,
};
