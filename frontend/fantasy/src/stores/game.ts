import { computed, ref } from 'vue';
import { defineStore } from 'pinia';
import {
  emptyGame,
  emptyPlayer,
  type Game,
  type Player,
  type Submission,
  type SubmissionStatus,
  type Team,
} from '@/types/game';

export type SocketResult = 'quit' | 'handled' | 'ignored';

interface SocketMessage {
  Command: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  Data: any;
}

// Team colours, in order. Every team-facing screen uses these so a team looks the same
// everywhere. Three are needed, a fourth is spare.
export const TeamColours = ['#e8467c', '#2f9fd8', '#f2a516', '#5bbf6a'];

export const useGameStore = defineStore('gameStore', () => {
  const Game = ref<Game>(emptyGame());
  const GameCode = ref(0);
  const Player = ref<Player>(emptyPlayer());
  const loading = ref(false);

  const teamById = computed(() => {
    const teams: Record<string, Team> = {};
    Game.value.Teams.forEach((t) => {
      teams[t.ID] = t;
    });
    return teams;
  });

  const myTeam = computed(() => teamById.value[Player.value.TeamID]);

  const teamColour = (teamID: string) => {
    const index = Game.value.Teams.findIndex((t) => t.ID === teamID);
    return TeamColours[Math.max(index, 0) % TeamColours.length];
  };

  const teamName = (teamID: string) => teamById.value[teamID]?.Name ?? '';

  const setGame = (game: Game) => {
    Game.value = game;
    GameCode.value = game.GameCode;
    loading.value = false;
  };

  const setPlayer = (player: Player) => {
    Player.value = player;
  };

  const reset = () => {
    Game.value = emptyGame();
    Player.value = emptyPlayer();
    GameCode.value = 0;
    loading.value = false;
  };

  const upsertPlayer = (player: Player) => {
    if (player.UUID === Player.value.UUID) {
      Player.value = player;
    }
    if (!Game.value.Players) {
      return;
    }
    const index = Game.value.Players.findIndex((p) => p.UUID === player.UUID);
    if (index >= 0) {
      Game.value.Players.splice(index, 1, player);
    } else {
      Game.value.Players.push(player);
    }
  };

  const removePlayer = (player: Player) => {
    if (!Game.value.Players) {
      return;
    }
    Game.value.Players = Game.value.Players.filter((p) => p.UUID !== player.UUID);
  };

  const upsertSubmission = (submission: Submission) => {
    if (!Game.value.Submissions) {
      Game.value.Submissions = [];
    }
    const index = Game.value.Submissions.findIndex((s) => s.ID === submission.ID);
    if (index >= 0) {
      Game.value.Submissions.splice(index, 1, submission);
    } else {
      Game.value.Submissions.push(submission);
    }
  };

  const addMySubmission = (submission: Submission) => {
    if (!Game.value.MySubmissions) {
      Game.value.MySubmissions = [];
    }
    Game.value.MySubmissions.push(submission);
  };

  const setMySubmissionStatus = (id: string, status: SubmissionStatus) => {
    const mine = Game.value.MySubmissions?.find((s) => s.ID === id);
    if (mine) {
      mine.Status = status;
    }
  };

  const setMyVote = (teamID: string) => {
    if (Game.value.Round) {
      Game.value.Round.MyVote = teamID;
    }
  };

  const setMyReaction = (reaction: number) => {
    if (Game.value.Round) {
      Game.value.Round.MyReaction = reaction;
    }
  };

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const setQueueOrder = (data: any) => {
    Game.value.NextTopicID = data.NextTopicID || '';
    Game.value.PinnedTopicID = data.PinnedTopicID || '';
    Game.value.TopPickIDs = data.TopPickIDs || [];
  };

  const sockMessages: Record<string, (data: SocketMessage['Data']) => void> = {
    GAME_UPDATE: (data) => setGame(data),
    NEW_PLAYER: (data) => upsertPlayer(data),
    PLAYER_UPDATE: (data) => upsertPlayer(data),
    PLAYER_QUIT: (data) => removePlayer(data),
    NEW_SUBMISSION: (data) => {
      upsertSubmission(data.Submission);
      setQueueOrder(data);
    },
    SUBMISSION_UPDATED: (data) => {
      upsertSubmission(data.Submission);
      setQueueOrder(data);
    },
    QUEUE_UPDATED: (data) => setQueueOrder(data),
    SUBMISSION_STATUS: (data) => setMySubmissionStatus(data.ID, data.Status),
    REACTION_TALLY: (data) => {
      const topic = Game.value.Round?.Topic;
      if (topic && topic.ID === data.SubmissionID) {
        topic.Ups = data.Ups;
        topic.Downs = data.Downs;
        topic.Net = data.Net;
      }
    },
    VOTE_TALLY: (data) => {
      const round = Game.value.Round;
      if (round && round.Number === data.Round) {
        round.Votes = data.Votes;
        round.TotalVotes = data.TotalVotes;
      }
    },
    VOTE_COUNT: (data) => {
      const round = Game.value.Round;
      if (round && round.Number === data.Round) {
        round.TotalVotes = data.TotalVotes;
      }
    },
    PRESENTING_TEAM: (data) => {
      if (Game.value.Round) {
        Game.value.Round.PresentingTeam = data.TeamID;
      }
    },
    SUBMISSIONS_ENABLED: (data) => {
      Game.value.SubmissionsEnabled = data.Enabled;
    },
    SET_VOLUME: (data) => {
      Game.value.Volume = data.Volume;
    },
  };

  const sockMsg = (msg: string): SocketResult => {
    if (msg.startsWith('CODE')) {
      GameCode.value = parseInt(msg.substring(5), 10);
      Game.value.GameCode = GameCode.value;
      return 'handled';
    }
    if (msg.startsWith('QUIT')) {
      return 'quit';
    }
    try {
      const data = JSON.parse(msg) as SocketMessage;
      const instruction = sockMessages[data.Command];
      if (instruction) {
        instruction(data.Data);
        return 'handled';
      }
    } catch {
      // Not a message for us
    }
    return 'ignored';
  };

  return {
    Game,
    GameCode,
    Player,
    loading,
    myTeam,
    teamById,
    addMySubmission,
    reset,
    setGame,
    setMyReaction,
    setMyVote,
    setPlayer,
    sockMsg,
    teamColour,
    teamName,
  };
});

export default useGameStore;
