import { beforeEach, describe, expect, it } from 'vitest';
import { createPinia, setActivePinia } from 'pinia';
import useGameStore from '../game';
import type { Game } from '@/types/game';

const msg = (Command: string, Data: unknown) => JSON.stringify({ Command, Data });

function game(overrides: Partial<Game> = {}): Game {
  return {
    Key: 'ABCDE',
    GameCode: 110,
    SubmissionsEnabled: true,
    ShowSubmitterNames: false,
    Volume: 50,
    SubmissionKinds: ['Question', 'Complaint', 'IncidentReport', 'Other'],
    Teams: [
      { ID: 't1', Name: 'Team One', RoundWins: 0, TotalVotes: 0, Members: [] },
      { ID: 't2', Name: 'Team Two', RoundWins: 0, TotalVotes: 0, Members: [] },
      { ID: 't3', Name: 'Team Three', RoundWins: 0, TotalVotes: 0, Members: [] },
    ],
    Round: {
      Number: 2,
      TeamOrder: ['t2', 't3', 't1'],
      PresentingTeam: 't2',
      TotalVotes: 0,
      WinnerTeamIDs: [],
      Decided: false,
      Topic: {
        ID: 's1',
        Kind: 'Complaint',
        Text: 'The lift smells',
        Status: 'Played',
        StaffVotes: 1,
        Ups: 0,
        Downs: 0,
        Net: 0,
        Created: '',
      },
    },
    ...overrides,
  };
}

describe('game store socket messages', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('follows game codes and quits on QUIT', () => {
    const store = useGameStore();
    expect(store.sockMsg('CODE 120')).toBe('handled');
    expect(store.GameCode).toBe(120);
    expect(store.sockMsg('QUIT')).toBe('quit');
    expect(store.sockMsg('nonsense')).toBe('ignored');
  });

  it('replaces the game on GAME_UPDATE', () => {
    const store = useGameStore();
    store.sockMsg(msg('GAME_UPDATE', game({ GameCode: 130 })));
    expect(store.GameCode).toBe(130);
    expect(store.teamName('t3')).toBe('Team Three');
  });

  it('applies tallies only to the round and topic they belong to', () => {
    const store = useGameStore();
    store.setGame(game());

    store.sockMsg(msg('REACTION_TALLY', { SubmissionID: 'other', Ups: 9, Downs: 9, Net: 0 }));
    expect(store.Game.Round?.Topic?.Ups).toBe(0);
    store.sockMsg(msg('REACTION_TALLY', { SubmissionID: 's1', Ups: 3, Downs: 1, Net: 2 }));
    expect(store.Game.Round?.Topic?.Net).toBe(2);

    store.sockMsg(msg('VOTE_TALLY', { Round: 1, Votes: { t1: 5 }, TotalVotes: 5 }));
    expect(store.Game.Round?.TotalVotes).toBe(0);
    store.sockMsg(msg('VOTE_TALLY', { Round: 2, Votes: { t1: 1, t2: 2, t3: 0 }, TotalVotes: 3 }));
    expect(store.Game.Round?.Votes).toEqual({ t1: 1, t2: 2, t3: 0 });

    store.sockMsg(msg('PRESENTING_TEAM', { TeamID: 't3' }));
    expect(store.Game.Round?.PresentingTeam).toBe('t3');
  });

  it('keeps the staff queue and the player list current', () => {
    const store = useGameStore();
    store.setGame(game({ Players: [] }));
    store.setPlayer({ Class: 'Audience', Name: 'Me', RoomKey: 'ABCDE', TeamID: '', UUID: 'me' });

    const submission = { ID: 'n1', Kind: 'Question', Text: 'Why?', Status: 'Pending', StaffVotes: 0, Ups: 0, Downs: 0, Net: 0, Created: '' };
    store.sockMsg(msg('NEW_SUBMISSION', { Submission: submission, NextTopicID: '' }));
    store.sockMsg(msg('SUBMISSION_UPDATED', { Submission: { ...submission, Status: 'Approved', StaffVotes: 1 }, NextTopicID: 'n1' }));
    expect(store.Game.Submissions).toHaveLength(1);
    expect(store.Game.Submissions?.[0]?.Status).toBe('Approved');
    expect(store.Game.NextTopicID).toBe('n1');

    store.sockMsg(msg('NEW_PLAYER', { Class: 'Audience', Name: 'Other', RoomKey: 'ABCDE', TeamID: '', UUID: 'other' }));
    store.sockMsg(msg('PLAYER_UPDATE', { Class: 'Director', Name: 'Me', RoomKey: 'ABCDE', TeamID: 't1', UUID: 'me' }));
    expect(store.Player.Class).toBe('Director');
    expect(store.myTeam?.Name).toBe('Team One');
    store.sockMsg(msg('PLAYER_QUIT', { UUID: 'other' }));
    expect(store.Game.Players?.map((p) => p.UUID)).toEqual(['me']);
  });

  it('gives every team its own colour', () => {
    const store = useGameStore();
    store.setGame(game());
    const colours = new Set(['t1', 't2', 't3'].map((id) => store.teamColour(id)));
    expect(colours.size).toBe(3);
  });
});
