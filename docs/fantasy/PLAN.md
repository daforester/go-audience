# Fantasy Furry Convention: Plan

Project codename and folder name: `fantasy`.
Ancestor: Trudge (`backend/apps/trudge`, `frontend/trudge`).

**Hard rule: no file under `backend/apps/trudge`, `frontend/trudge`, `backend/services`,
`backend/libs`, `backend/vendor`, `backend/go.mod` or the trudge Docker/config files is
modified.** Fantasy copies what it needs and evolves the copy.

---

## 1. Game summary

| Trudge role            | Fantasy role        | Count     | Purpose                                                                 |
|------------------------|---------------------|-----------|-------------------------------------------------------------------------|
| Leader                 | **Admin**           | 1+        | Runs the show: phases, teams, role assignment, moderation overrides     |
| Display                | **Display**         | 1+        | Projector / main screen                                                 |
| God                    | **Staff**           | many      | Triage the submission queue: approve/reject and upvote                  |
| Warrior / Wizard / Bard| **Director**        | many      | Belongs to one **Director Team**. Teams answer each topic               |
| Audience               | **Audience**        | many      | Submit problems, react ±1 to the live topic, vote for best team answer  |

### Round loop

```
Lobby ─► Intake (submissions open) ─┐
            ▲                       ▼
            │                 Topic Reveal       highest staff-voted approved, unused submission
            │                       ▼
            │                 Responses          teams present in turn; audience can ±1 the topic
            │                       ▼
            │                 Audience Vote      one vote per audience member, for a team
            │                       ▼
            └──────────────── Round Result       winning team revealed, scores updated
                                    ▼ (admin ends game)
                              Final Tally        winning Director Team + winning audience submission
```

Submission intake and staff triage run **continuously** in the background of every phase
(unless the admin closes submissions), so the queue keeps filling between rounds.

### Scoring (default proposal, see open questions)

- Team round score: the team with most audience votes wins the round: **+1 round win**.
  Also keep the raw vote totals per team as the tie-breaker for the final tally.
- Submission score: `upvotes - downvotes` from the audience while it was the live topic.
  The winning submission is the highest net score among submissions that were played.
- Staff votes decide *which* item gets played; they do not count towards the winning submission.

---

## 2. Technology: "latest everything"

Trudge shares `backend/go.mod` (Go 1.19, gin 1.8, `dgrijalva/jwt-go`, `mgo.v2`, vendored).
Upgrading that module would change Trudge, so **Fantasy gets its own Go module** at
`backend/apps/fantasy/go.mod`. A nested `go.mod` automatically excludes the directory from
the parent module, so Trudge's `-mod=vendor` build is unaffected.

Fantasy must not import `bitbucket.org/daforester/go-audience/services/...` (they pull the old
dependency set). Copy the small pieces needed (`controller`, `router`, `ginhelper/bearerauth`)
into `backend/apps/fantasy/internal/...` and modernise them.

Pin exact versions when scaffolding (`go get -u`, `npm create vue@latest`), not in this plan.

### Backend

| Concern        | Trudge                                   | Fantasy                                                       |
|----------------|------------------------------------------|---------------------------------------------------------------|
| Go             | 1.19 (Docker image 1.18)                 | Current toolchain (1.27 installed locally), matching Docker image |
| HTTP           | gin 1.8                                  | gin latest (keep gin for familiarity)                         |
| JWT            | `dgrijalva/jwt-go` (archived)            | `github.com/golang-jwt/jwt/v5`                                |
| WebSocket      | gorilla/websocket 1.5.0                  | gorilla/websocket latest, or `github.com/coder/websocket` (context-aware) |
| Config         | viper 1.15                               | viper latest                                                  |
| Logging        | logrus                                   | stdlib `log/slog`                                             |
| DI             | `daforester/go-di-container` (reflection)| Explicit constructor wiring in `cmd/fantasy/main.go` (see open questions) |
| Sessions       | gin-contrib/sessions + mgo.v2            | Drop: JWT bearer tokens are all the game uses                 |
| IDs            | google/uuid                              | google/uuid latest                                            |
| Persistence    | in-memory map + JSON dump on SIGINT      | Same model, but mutex-correct, atomic write (temp + rename), periodic snapshots, SIGTERM handled |
| Tests          | none for trudge                          | `go test -race` on domain + processors; httptest for API; ws integration test |

### Frontend

| Concern     | Trudge                                  | Fantasy                                                    |
|-------------|-----------------------------------------|------------------------------------------------------------|
| Scaffold    | create-vue (2022)                       | `npm create vue@latest` → Vue 3 latest, TypeScript, Router, Pinia, Vitest, ESLint flat config, Prettier |
| Components  | Options API, `.vue` JS                  | `<script setup lang="ts">`, Composition API, typed props/emits |
| State       | Pinia 2 + persistedstate 3, mitt bus    | Pinia latest + pinia-plugin-persistedstate latest; event bus replaced by stores/composables |
| HTTP        | axios                                   | axios latest (keeps TokenManager logic portable) or `fetch` wrapper |
| WebSocket   | `WebSocket.vue` component               | `useSocket()` composable with reconnect + typed message dispatch |
| CSS         | Bootstrap 4 + Sass                      | Bootstrap 5.3 + Sass (modern `@use`, no `@import`)          |
| Icons       | FontAwesome Pro 5.15.4 tarballs         | FontAwesome latest (free set unless a current Pro licence/npm token is available) |
| Tooling     | Vite 4, TS 4.7, Vitest 0.25             | Vite, TypeScript, Vitest, vue-tsc latest                   |
| Node        | —                                       | Node 24 LTS (installed locally); add `.nvmrc` / `engines`  |

---

## 3. Target layout (new files only)

```
backend/
  apps/fantasy/
    go.mod, go.sum
    cmd/fantasy/main.go            wiring, server, graceful shutdown, snapshots
    internal/
      config/                      viper setup (from trudge setup/config.go, keys.go)
      auth/                        JWT issue/verify (golang-jwt v5), bearer middleware
      domain/                      Game, Player, Team, Submission, Round, Phase; pure logic + tests
      repository/                  GameRepository + MemoryGameRepository (race-free)
      processors/                  one per command (same pattern as trudge)
      http/controllers/            gin handlers (same pattern as trudge)
      http/router/                 api + ws routes
      realtime/                    hub, rooms, role-tagged connections, typed messages
  config-fantasy.toml              NO secrets committed: keys from env vars
  config-testing-fantasy.toml
  Dockerfile-fantasy               builds from backend/apps/fantasy
  docker-compose-fantasy.yml
frontend/
  fantasy/                         fresh create-vue project
    src/
      api/                         typed client, token refresh (port of TokenManager)
      realtime/useSocket.ts        socket composable + message types
      stores/                      session.ts, game.ts, submissions.ts
      views/                       Login, Admin, Display, Staff, Director, Audience
      components/{admin,display,staff,director,audience,shared}/
      styles/
docs/fantasy/PLAN.md               this file
```

---

## 4. Domain model (backend `internal/domain`)

```go
type Role string // Admin, Display, Staff, Director, Audience

type Phase string // Lobby, Intake, TopicReveal, Responses, AudienceVote, RoundResult, FinalTally

type Player struct {
    UUID, Name, RoomKey string
    Role   Role
    TeamID string // Directors only
}

type Team struct {
    ID, Name  string
    RoundWins int
    TotalVotes int
}

type SubmissionKind string // Question, Complaint, IncidentReport, Other

type SubmissionStatus string // Pending, Approved, Rejected, Played

type Submission struct {
    ID, AuthorUUID, AuthorName string
    Kind      SubmissionKind
    Text      string
    Status    SubmissionStatus
    StaffVotes map[string]struct{} // staff UUID set: one upvote each, toggleable
    Reactions  map[string]int8     // audience UUID -> +1 / -1 (changeable, not stackable)
    CreatedAt time.Time
}

type Round struct {
    Number       int
    SubmissionID string
    TeamOrder    []string          // presentation order
    PresentingTeam string
    Votes        map[string]string // audience UUID -> team ID (one vote each, changeable until close)
    WinnerTeamIDs []string         // ties allowed
}

type Game struct {
    Key        string
    Phase      Phase
    SubmissionsOpen bool
    Players    []*Player
    Teams      []*Team
    Submissions []*Submission
    Rounds     []*Round
    Version    uint64 // bumps on every mutation; clients drop stale deltas
    Volume     int
    Created, Updated time.Time
}
```

Key rules (all unit-tested, all in `domain`, no I/O):

- `NextTopic()`: highest `len(StaffVotes)` among `Approved`; tie-break oldest `CreatedAt`.
  Admin can override with a specific submission ID.
- Votes and reactions are **maps keyed by voter UUID**, so a voter can never vote twice.
  (Trudge increments counters with no per-voter record, so repeat votes double-count, and its
  `Votes` padding path appends nil maps that would panic if written.)
- Phase transitions are validated by a small state machine (`CanTransition(from, to)`);
  the admin UI only shows legal next steps. Replaces Trudge's magic `GameCode` integers.
- Directors and staff cannot vote in the audience vote; audience cannot see the staff queue.

---

## 5. API (REST commands + WebSocket events)

Same split as Trudge: REST for commands, WebSocket for broadcasts. All under `/api`.

| Method | Path                                   | Role          | Purpose                               |
|--------|----------------------------------------|---------------|---------------------------------------|
| POST   | /game                                  | anyone        | Create game, creator becomes Admin    |
| POST   | /game/join                             | anyone        | Join as Audience                      |
| GET    | /game                                  | any player    | Snapshot filtered by role             |
| POST   | /token/refresh                         | any player    | Refresh JWT                           |
| POST   | /game/quit                             | any player    |                                       |
| POST   | /game/phase                            | Admin         | Transition phase                      |
| POST   | /game/submissions/toggle               | Admin         | Open/close intake                     |
| POST   | /teams, PATCH/DELETE /teams/:id        | Admin         | Manage director teams                 |
| POST   | /players/:uuid/role                    | Admin         | Assign role (+ team for Directors)    |
| DELETE | /players/:uuid                         | Admin         | Kick                                  |
| POST   | /rounds/topic                          | Admin         | Start round (auto top item or by ID)  |
| POST   | /rounds/current/presenting             | Admin         | Set presenting team                   |
| POST   | /submissions                           | Audience      | New submission (kind + text, length-limited, rate-limited) |
| POST   | /submissions/:id/status                | Staff, Admin  | Approve / reject                      |
| POST   | /submissions/:id/staff-vote            | Staff         | Toggle upvote                         |
| POST   | /rounds/current/reaction               | Audience      | +1 / -1 / 0 on live topic             |
| POST   | /rounds/current/vote                   | Audience      | Vote for team                         |
| POST   | /game/volume                           | Admin, Display| Music volume (as Trudge)              |

WebSocket: `GET /ws/:room?token=<jwt>`. Token is verified and the player resolved server-side
(Trudge trusts a UUID in the path). `CheckOrigin` uses the CORS allowlist (Trudge accepts all).

Message envelope: `{"type": "...", "version": N, "data": {...}}`. On connect the server sends
`SNAPSHOT` (role-filtered); afterwards deltas:

`PHASE_CHANGED`, `PLAYER_JOINED`, `PLAYER_LEFT`, `ROLE_ASSIGNED`, `TEAM_UPDATED`,
`SUBMISSION_NEW` (staff/admin only), `SUBMISSION_UPDATED` (staff/admin), `TOPIC_REVEALED`,
`PRESENTING_TEAM`, `REACTION_TALLY` (throttled ~4/s), `VOTE_TALLY` (admin/display, throttled),
`ROUND_RESULT`, `FINAL_TALLY`, `VOLUME`.
Keep the `PING`/`PONG` keepalive so the client logic stays familiar.

Role-filtered broadcast replaces Trudge's tag-based `HasTag(...)` calls: the hub exposes
`Broadcast(room, msg, roles...)`.

---

## 6. Frontend views

| View      | Shows                                                                                         |
|-----------|-----------------------------------------------------------------------------------------------|
| Login     | Room code + name; "create game" for admin (port Trudge `Login.vue`)                           |
| Admin     | Phase controls (legal next steps only), player list with role/team assignment, team editor, submission override/hide, live vote and reaction tallies, intake toggle, volume |
| Display   | Lobby/QR + room code, topic card, presenting-team banner, live reaction meter, vote bars, round winner reveal, final podium (team + submission). Port Trudge backgrounds/animations where they fit |
| Staff     | Queue tabs (Pending / Approved / Rejected), approve/reject, upvote toggle, sorted by votes; preview of "next up" |
| Director  | Team name and teammates, current topic, "your team is presenting" cue, scoreboard            |
| Audience  | Submission form (kind picker + text) while intake open; ±1 reaction on the live topic; team vote during Audience Vote; otherwise PleaseWait |

Use router routes per role (`/admin`, `/display`, `/staff`, `/director`, `/play`) guarded by the
session store, instead of Trudge's single `Game.vue` with a `v-if` per class.

---

## 7. Delivery phases and branches

Integration branch: **`feature/fantasy`** (off `main`). Phase branches come off it and merge
back via PR. `main` only receives `feature/fantasy` once the game is playable end-to-end.

| # | Branch                              | Scope                                                                                     | Done when |
|---|-------------------------------------|-------------------------------------------------------------------------------------------|-----------|
| 0 | `feature/fantasy`                   | This plan                                                                                 | Plan agreed |
| 1 | `feature/fantasy-backend-scaffold`  | New Go module, copied+modernised config/auth/controller/router, health endpoint, create/join/rejoin/refresh, Dockerfile, CI `go build` + `go vet` + `go test -race` | `go test ./...` green; trudge `go build -mod=vendor ./apps/trudge/...` still green |
| 1 | `feature/fantasy-frontend-scaffold` | `npm create vue@latest`, Bootstrap 5, icons, lint/format, API client + token refresh port, socket composable, Login view | `npm run build`, `lint`, `test:unit` green |
| 2 | `feature/fantasy-domain`            | Domain model, state machine, NextTopic, vote/reaction/tally logic, repository, unit tests | Full round simulated in a Go test |
| 3 | `feature/fantasy-realtime`          | Hub, token-authenticated ws, role-filtered broadcast, snapshot + deltas, throttled tallies | Integration test: 2 clients see correct role-filtered events |
| 4 | `feature/fantasy-api`               | All REST commands from section 5 with role guards, validation, rate limits                | httptest coverage per endpoint/role |
| 5 | `feature/fantasy-views-core`        | Admin, Audience, Staff views wired to stores                                              | Manual loop: submit → triage → topic |
| 6 | `feature/fantasy-views-stage`       | Display and Director views, reveal animations, final tally, music/volume                  | Full playthrough on 5 browser windows |
| 7 | `feature/fantasy-hardening`         | Snapshot/restore on restart, reconnection, load test (~200 audience sockets), profanity/length guards, docs, deploy config | Rehearsal game survives a backend restart |

Phases 1a/1b run in parallel; 5/6 can run in parallel once 3 and 4 land.

---

## 8. Things deliberately not carried over from Trudge

- Committed session keys in `config-trudge.toml`: Fantasy reads keys from env vars and fails
  fast if missing in production.
- `CheckOrigin: return true` and UUID-in-path socket auth.
- Map access in `MemoryGameRepository` outside the lock, and `GetByKey` before `Lock`
  (data race): Fantasy locks first, then reads; `go test -race` enforces it.
- Character/race/equipment/dice/map-maker features: not part of this game.
- `vendor/` directory: not needed with a module proxy; add later only if offline Docker
  builds require it.

---

## 9. Open questions

1. **Scoring**: one point per round won (proposed), or points proportional to vote share?
2. **Teams**: how many director teams typically, and does the admin name them or do directors?
3. **Anonymity**: show submitter names on the Display and in the final "winning submission"?
4. **Audience reactions**: allowed during the whole round, or only while teams are presenting?
5. **FontAwesome**: is there a current Pro licence (needs an npm token for v6/v7 Pro), or use the free set?
6. **DI container**: OK to replace `go-di-container` with explicit wiring, or keep the same
   pattern for consistency with Trudge/PriceIsRight?
7. **Hosting**: domain / CORS origins for production (Trudge uses `trudge.party`).
