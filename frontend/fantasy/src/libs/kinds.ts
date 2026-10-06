import type { SubmissionKind } from '@/types/game';

export const KindIcons: Record<SubmissionKind, string> = {
  Question: 'circle-question',
  Complaint: 'face-angry',
  IncidentReport: 'triangle-exclamation',
  Other: 'bullhorn',
};

export const KindPrompts: Record<SubmissionKind, string> = {
  Question: 'Ask the con directors anything...',
  Complaint: 'What has gone wrong at this convention?',
  IncidentReport: 'Report what you saw, in as much detail as you dare...',
  Other: 'Anything else the directors should deal with?',
};
