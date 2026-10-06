// Mirrors the views the backend sends (backend/apps/fantasy/models/game_view.go)

export type PlayerClass = 'Admin' | 'Audience' | 'Director' | 'Display' | 'Staff';

export type SubmissionKind = 'Question' | 'Complaint' | 'IncidentReport' | 'Other';

export type SubmissionStatus = 'Pending' | 'Approved' | 'Rejected' | 'Played';

export const GameCodes = {
  Lobby: 20,
  Intake: 30,
  TopicReveal: 100,
  Responses: 110,
  AudienceVote: 120,
  RoundResult: 130,
  FinalTally: 600,
} as const;

export interface Player {
  Class: PlayerClass | '';
  Name: string;
  RoomKey: string;
  TeamID: string;
  UUID: string;
}

export interface PlayerSummary {
  UUID: string;
  Name: string;
}

export interface Team {
  ID: string;
  Name: string;
  RoundWins: number;
  TotalVotes: number;
  Members: PlayerSummary[];
}

export interface Submission {
  ID: string;
  Kind: SubmissionKind;
  Text: string;
  AuthorName?: string;
  AuthorUUID?: string;
  Status: SubmissionStatus;
  Vetoed?: boolean;
  StaffVotes: number;
  MyStaffVote?: boolean;
  Ups: number;
  Downs: number;
  Net: number;
  Created: string;
}

export interface Round {
  Number: number;
  Topic?: Submission;
  TeamOrder: string[];
  PresentingTeam: string;
  Votes?: Record<string, number>;
  TotalVotes: number;
  WinnerTeamIDs: string[];
  Decided: boolean;
  MyVote?: string;
  MyReaction?: number;
}

export interface Final {
  WinnerTeamIDs: string[];
  WinningSubmission?: Submission;
}

export interface Game {
  Key: string;
  GameCode: number;
  SubmissionsEnabled: boolean;
  ShowSubmitterNames: boolean;
  Volume: number;
  MaxTeams?: number;
  SubmissionKinds: SubmissionKind[];
  Players?: Player[];
  Teams: Team[];
  Submissions?: Submission[];
  MySubmissions?: Submission[];
  NextTopicID?: string;
  PinnedTopicID?: string;
  TopPickIDs?: string[];
  Round?: Round;
  Final?: Final;
}

export interface Tokens {
  AccessToken: string;
  AccessExpires: number;
  RefreshToken: string;
  RefreshExpires: number;
  Endpoint?: string;
  Service?: string;
  Vendor?: string;
}

export const SubmissionKindLabels: Record<SubmissionKind, string> = {
  Question: 'Question',
  Complaint: 'Complaint',
  IncidentReport: 'Incident Report',
  Other: 'Other',
};

export function emptyGame(): Game {
  return {
    Key: '',
    GameCode: 0,
    SubmissionsEnabled: false,
    ShowSubmitterNames: false,
    Volume: 50,
    SubmissionKinds: ['Question', 'Complaint', 'IncidentReport', 'Other'],
    Teams: [],
  };
}

export function emptyPlayer(): Player {
  return {
    Class: '',
    Name: '',
    RoomKey: '',
    TeamID: '',
    UUID: '',
  };
}
