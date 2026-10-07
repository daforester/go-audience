# Fantasy Furry Convention: Plan

Project codename and folder name: `fantasy`.
Ancestor: Trudge (`backend/apps/trudge`, `frontend/trudge`).

**Hard rule: no file under `backend/apps/trudge`, `frontend/trudge`, `backend/services`,
`backend/libs`, `backend/vendor`, `backend/go.mod` or the trudge Docker/config files is
modified.** Fantasy copies what it needs and evolves the copy.

**Style rule: Fantasy should read like a natural next step from Trudge, written by the same
author.** Same layout, same layering (router → controller → processor → repository / messages),
same `go-di-container` wiring, same naming conventions (`XxxController`, `XxxProcessor`,
`XxxProcessorInterface`, `New(...)` constructors, receiver names `C` / `P` / `r`), same
frontend shape (Options API `.vue` files, `views/`, `components/<role>/`, `mixins/`,
`TokenManager.vue`, `WebSocket.vue`, mitt emitter). Library versions move to latest; the
code style does not change. Where an upgrade forces an API change (e.g. jwt v5), keep the
surrounding code shape identical.

The one deliberate framework change is **gin → Fiber v3** (section 2a). Routers, controllers,
processors and bindings keep their Trudge shape; only the request context type changes
(`*gin.Context` → `fiber.Ctx`).

---

## 1. Game summary

| Trudge role            | Fantasy role   | Constant                     | Count | Purpose                                                            |
|------------------------|----------------|------------------------------|-------|--------------------------------------------------------------------|
| Leader                 | **Admin**      | `PLAYER_CLASS_ADMIN`         | 1+    | Runs the show: phases, teams, role assignment, settings, moderation overrides |
| Display                | **Display**    | `PLAYER_CLASS_DISPLAY`       | 1+    | Projector / main screen                                            |
| God                    | **Staff**      | `PLAYER_CLASS_STAFF`         | many  | Triage the submission queue: approve/reject and upvote             |
| Warrior / Wizard / Bard| **Director**   | `PLAYER_CLASS_DIRECTOR`      | many  | Belongs to one **Director Team**; teams answer each topic          |
| Audience               | **Audience**   | `PLAYER_CLASS_AUDIENCE`      | many  | Submit problems, react ±1 while teams present, vote for best team  |

### Director teams

- Typical game: **2 teams**. Must work equally well with **3**.
- Teams are a slice on the game (`Game.Teams []*Team`), never hard-coded pairs. A
  `MAX_TEAMS` config value (default `3`) caps how many the admin can create; minimum 2 to
  start a round.
- Every team-facing layout (display vote bars, round reveal, final podium, admin tallies)
  is built for N teams and visually checked at 2 and 3.

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
            └──────────────── Round Result       winning team revealed, round win awarded
                                    ▼ (admin ends game)
                              Final Tally        winning Director Team + winning audience submission
```

Submission intake and staff triage run **continuously** in the background of every phase
(unless the admin closes submissions, like Trudge's suggestions toggle), so the queue keeps
filling between rounds.

### Scoring

- **Teams: one point per round won.** The team with the most audience votes in a round gets
  `RoundWins += 1`. If teams tie for the most votes, each tied team gets the win.
  Raw vote totals are also accumulated per team (`TotalVotes`) and used only to break a tie
  in the final tally; if still tied, the final reveal shows joint winners.
- **Submissions:** audience members react +1 / -1 to the live topic **only during the
  Responses phase** (while teams are presenting). Net score = ups - downs. The winning
  submission is the highest net score among played submissions; tie-break by most total
  reactions, then earliest played.
- Staff votes decide *which* item gets played; they do not count towards the winning
  submission.

### Submitter names (admin toggle)

`Game.ShowSubmitterNames bool`, toggled from the admin panel (default `false`), broadcast
the same way Trudge broadcasts `SUGGESTIONS_ENABLED`.

- Admin and Staff always see the author (needed for moderation).
- Display, Directors and Audience receive `AuthorName` only when the toggle is on; otherwise
  the server blanks it before sending. Filtering happens server-side, not in the UI.
- Applies everywhere: topic reveal, round result and the final "winning submission" reveal.

---

## 2. Technology: latest libraries, same shape

Trudge shares `backend/go.mod` (Go 1.19, gin 1.8, `dgrijalva/jwt-go`, `mgo.v2`, vendored).
Upgrading it would change Trudge, so **Fantasy gets its own Go module** at
`backend/apps/fantasy/go.mod`, module path `bitbucket.org/daforester/go-audience/apps/fantasy`.
A nested `go.mod` excludes the directory from the parent module, so Trudge's `-mod=vendor`
build is unaffected.

`backend/services` is built on gin, `ginhelper` imports `dgrijalva/jwt-go`, and
`services/sessions` imports `mgo.v2`, so Fantasy can't reuse them. It gets its own copy, laid
out the same way and ported to Fiber: `backend/apps/fantasy/services/{controller,router,fiberhelper}`
(plus `libs/slice` if needed). Import paths change from `.../go-audience/services/controller`
to `.../go-audience/apps/fantasy/services/controller`; the code inside stays recognisably the same.

### Backend

| Concern      | Trudge                                  | Fantasy                                                        |
|--------------|-----------------------------------------|----------------------------------------------------------------|
| Go           | 1.19 (Docker image 1.18)                | Current toolchain (1.27 locally); Docker image to match        |
| DI           | `daforester/go-di-container` v1.0.2     | **Same library**, latest tag; `bindings/bindings.go` pattern kept |
| HTTP         | gin 1.8                                 | **Fiber v3** (`github.com/gofiber/fiber/v3`), latest v3.x tag  |
| CORS         | gin-contrib/cors 1.4 with hard-coded origins | Fiber `middleware/cors` for the allow-list, plus a small custom middleware for naughty mode (section 5) |
| JWT          | `dgrijalva/jwt-go` (archived)           | `github.com/golang-jwt/jwt/v5` (drop-in style: `jwt.Parse`, `jwt.MapClaims`) |
| WebSocket    | gorilla/websocket 1.5.0                 | `github.com/gofiber/contrib/v3/websocket` (fasthttp/websocket, a gorilla fork: same `ReadMessage`/`WriteMessage` and message types) |
| Config       | viper 1.15                              | viper latest                                                   |
| Logging      | logrus                                  | logrus latest (kept for continuity)                            |
| Location     | gin-contrib/location                    | Dropped: `HTTP_SCHEMA`/`HTTP_HOST`/`HTTP_BASE` config read directly where needed |
| IDs          | google/uuid                             | google/uuid latest                                             |
| Sessions     | gin-contrib/sessions + mgo.v2           | Dropped: Trudge registers sessions but no controller uses them; auth is JWT bearer only |
| Persistence  | memory map + JSON dump on SIGINT        | Same `MemoryGameRepository` + `LoadFromDisk`/`SaveToDisk`, with the lock fixes in section 8; SIGINT and SIGTERM, then `app.ShutdownWithTimeout` |
| Tests        | none for trudge                         | `go test -race` on models + processors; `app.Test(req)` for controllers |

### 2a. Fiber v3 migration: keeping websocket behaviour and staying safe

**What changes in the code (and what doesn't):**

```go
// Trudge (gin)
router.POST("/api/game/option/vote", func(c *gin.Context) {
    P.RunController(c, &http2.VoteOptionController{})
})

// Fantasy (Fiber v3)
router.Post("/api/game/round/vote", func(c fiber.Ctx) error {
    return P.RunController(c, &http2.VoteTeamController{})
})
```

- `services/router.BaseRouter.RunController` still does `B.container.Make(C)`, `SetContext(c)`,
  `Run()`; it now returns the error so Fiber's error handler sends the 500.
- Controllers keep `C.GetContext().(fiber.Ctx)`; `c.BindJSON(x)` becomes `c.Bind().JSON(x)`,
  `c.Param("uuid")` becomes `c.Params("uuid")`, `c.JSON(code, e)` becomes `c.Status(code).JSON(e)`.
- `controllers/http` keeps its folder and package name, matching Trudge.

**Websocket behaviour is kept:**

| Trudge behaviour                                        | Fantasy on Fiber v3                                               |
|---------------------------------------------------------|-------------------------------------------------------------------|
| `GET /ws/:room/:player`, upgrade in `EstablishSocketController` | Same route; a pre-upgrade middleware validates token + origin (`websocket.IsWebSocketUpgrade(c)`), then `websocket.New(handler)` |
| `sockets.Collection` / `Connection` with tags, `Related`, write mutex | Unchanged API; `Connection` wraps `*websocket.Conn` from contrib. Mutex kept: still one concurrent writer per socket |
| `CODE <n>` + settings messages on connect               | Unchanged                                                         |
| `PING`/`PONG` text keepalive in `ReadSocket`            | Unchanged                                                         |
| Blocking read loop, then remove connection + empty collection | Unchanged; the handler blocks in `ReadSocket` exactly as now  |
| Room/player from path params                            | Read through the copying helpers (rule below) and stored as copies on the `Connection`, never as references into the request |

**Security risks Fiber introduces, and how Fantasy closes each one:**

1. **fasthttp buffer reuse.** Strings and byte slices returned by `c.Params()`, `c.Query()`,
   `c.Get()`, `c.Cookies()`, `c.FormValue()` and `c.Body()` point into a request buffer that
   fasthttp reuses once the handler returns. Fantasy keeps request data in long-lived models
   (player names, submission text) and hands it to socket goroutines, so a retained reference
   could silently turn one user's data into another's.

   **Rule: data from requests is copied, never pointed to.**
   - `services/fiberhelper` provides the only accessors controllers use:
     `fiberhelper.Param(c, "uuid")`, `Query`, `Header`, `Cookie`, `FormValue` (each returns
     `strings.Clone(...)`) and `Body(c)` (returns `bytes.Clone(c.Body())`). Bearer tokens are
     read the same way in `fiberhelper.GetBearerToken`.
   - JSON bodies are bound into a fresh request struct (`app/requests/*.go`, as in Trudge) with
     the default `encoding/json` decoder, which allocates new strings. The decoder is pinned in
     `fiber.Config`, so swapping in a zero-copy decoder that could alias the buffer has to be
     a deliberate change.
   - Processors build **new** model values from request fields
     (`player.Name = request.Name`, `submission.Text = request.Text`). They never store a
     pointer to a request struct, to the JWT-decoded `*models.Player`, or to anything reachable
     from `fiber.Ctx`. The stored player is always looked up from the repository by UUID.
   - `fiber.Config{Immutable: false}` (the default) on purpose, so tests exercise the real
     buffer reuse and catch a missed copy, rather than Immutable hiding it.
   - Enforcement: a unit test in `services/fiberhelper` greps `app/controllers` and `middleware`
     and fails on direct `c.Params(`/`c.Query(`/`c.Get(`/`c.Body(`/`c.Cookies(`/`c.FormValue(`
     calls. A `-race` integration test fires concurrent joins/submissions with distinct
     payloads and asserts every stored name and text is exactly what was sent.
2. **Ctx pooling.** `fiber.Ctx` is recycled after the handler. Controllers are created per
   request by `RunController`, so they never outlive it. Rule: never capture `c` in a
   goroutine or store it in a model, a socket `Related` map or a processor; pass copied values
   instead. The websocket handler copies `room`/`player` out of the upgrade request into the
   `Connection` before it starts reading.
3. **Websocket origin check defaults to open.** contrib `websocket.Config.Origins` empty means
   *any origin*. Fantasy never relies on it. The pre-upgrade middleware always runs the same
   `setup.OriginAllowed(c)` check as CORS (section 5), and rejects with 403 before upgrading.
4. **CORS origin callback has no request.** Fiber's `AllowOriginsFunc(origin string) bool` can't
   see the request Host, so naughty mode is a small custom middleware (section 5), not a
   permissive `AllowOriginsFunc`. Credentials are never combined with `*`; Fiber panics on
   that combination anyway.
5. **Proxy headers.** `c.Hostname()` / `c.IP()` trust `X-Forwarded-*` only when
   `TrustProxy: true` and the caller is in `TrustProxyConfig`. Fantasy maps `TRUST_PROXY` and
   `TRUSTED_PROXIES` config onto these, default off, so naughty mode can't be fooled by a
   spoofed `X-Forwarded-Host` from the internet.
6. **Request size and slow clients.** Set `BodyLimit` to 64 KB (default 4 MB), and set
   `ReadTimeout`, `WriteTimeout` and `IdleTimeout` (all default to unlimited). Websocket read limit
   via `conn.SetReadLimit` (e.g. 4 KB, since clients only send PINGs).
7. **Panics.** `middleware/recover` on HTTP routes; `websocket.Config.RecoverHandler` on the ws
   route (the generic recover middleware does not cover the socket goroutine).

**Why it's worth it:** Fiber v3 is actively maintained, its contrib websocket keeps the
gorilla programming model, and the routing reads almost the same as gin. The risks above are
configuration, settled once in `cmd/fantasy.go` and `setup/`, and covered by tests in phase 1.

### Frontend

Fresh `create-vue` project at `frontend/fantasy`, then reshaped to mirror `frontend/trudge`:
same `main.ts` (provide/globalProperties, mitt emitter), same `router/index.ts`,
`stores/local.ts` with persisted tokens, `WebApp.vue`, `TokenManager.vue`, `WebSocket.vue`,
`mixins/`, Options API components in `.vue` files.

| Concern     | Trudge                              | Fantasy                                                         |
|-------------|-------------------------------------|-----------------------------------------------------------------|
| Framework   | Vue 3.2, vue-router 4.1             | Vue latest, vue-router latest                                   |
| State       | Pinia 2 + persistedstate 3, mitt 3  | Pinia, pinia-plugin-persistedstate, mitt latest                 |
| HTTP        | axios 1.2                           | axios latest                                                    |
| Build       | Vite 4, TS 4.7, vue-tsc 1, Vitest 0.25 | Vite, TypeScript, vue-tsc, Vitest latest                     |
| Lint        | ESLint 8 `.eslintrc.cjs`, Prettier 2 | ESLint latest (flat `eslint.config.js`), Prettier latest       |
| CSS         | Bootstrap 4 + Sass `@import`        | Bootstrap 5.3 + Sass latest (`@use`; Sass has removed `@import` support) |
| Icons       | FontAwesome Pro 5.15.4 tarballs + free-solid | **FontAwesome Free only** (`@fortawesome/free-solid-svg-icons`, `free-regular-svg-icons`, `vue-fontawesome`), same `libs/fontawesome/minimal.js` registration pattern. No Pro tarballs. |
| Node        | —                                   | Node 24 LTS; add `engines` to `package.json`                    |

---

## 3. Target layout (new files only, mirrors Trudge)

```
backend/
  apps/fantasy/
    go.mod, go.sum
    app/
      controllers/http/            one file per endpoint, e.g. new_submission.go, staff_vote.go
      helpers/                     get_request_player.go
      messages/                    one file per broadcast, e.g. topic_revealed_message.go
      processors/                  one file per command, e.g. new_submission_processor.go
      requests/                    create_game.go, join_game.go, assign_player.go, ...
      routers/                     api_router.go, websocket_router.go
      sockets/                     collection.go, connection.go, reader.go
    bindings/bindings.go           app.Singleton / app.Bind registrations
    cmd/fantasy.go                 main, same shape as cmd/trudge.go
    middleware/                    mustbeplayer.go, mustberole.go, naughty_cors.go, websocket_upgrade.go
    models/                        game.go, player.go, team.go, submission.go, round.go
    repositories/                  game_repository.go, memory_game_repository.go
    services/                      controller/, router/, fiberhelper/ (copied, Fiber v3 + jwt v5)
    setup/                         config.go, keys.go, logger.go, cors.go
  config-fantasy.toml              no secrets committed: session keys from env vars
  config-testing-fantasy.toml
  Dockerfile-fantasy               builds from backend/apps/fantasy
  docker-compose-fantasy.yml
frontend/
  fantasy/
    src/
      components/
        admin/  audience/  backgrounds/  director/  display/  staff/  inputs/  player/
        Login.vue  TokenManager.vue  WebSocket.vue
      libs/{debug,fontawesome,scss}/
      mixins/api/CachedRequest.vue
      router/index.ts
      stores/local.ts
      views/  Game.vue  Admin.vue  Display.vue  Staff.vue  Director.vue  Audience.vue  Logout.vue
      main.ts  WebApp.vue
docs/fantasy/PLAN.md               this file
```

`views/Game.vue` keeps Trudge's approach: one route, a socket, and a `v-if` per
`Player.Class` choosing the role view.

---

## 4. Models (`apps/fantasy/models`)

Written in the style of `apps/trudge/models`: plain structs, typed string constants, helper
methods on `*Game`.

```go
type PlayerClass string

const (
    PLAYER_CLASS_ADMIN    PlayerClass = "Admin"
    PLAYER_CLASS_AUDIENCE PlayerClass = "Audience"
    PLAYER_CLASS_DIRECTOR PlayerClass = "Director"
    PLAYER_CLASS_DISPLAY  PlayerClass = "Display"
    PLAYER_CLASS_STAFF    PlayerClass = "Staff"
)

var FantasyPlayerClassesStrings = []string{ /* as TrudgePlayerClassesStrings */ }

type Player struct {
    Class   PlayerClass
    Name    string
    RoomKey string
    TeamID  string // Directors only
    UUID    string
}

type Team struct {
    ID         string
    Name       string
    RoundWins  int
    TotalVotes int
}

type SubmissionKind string   // Question, Complaint, IncidentReport, Other
type SubmissionStatus string // Pending, Approved, Rejected, Played

type Submission struct {
    ID         string
    AuthorUUID string
    AuthorName string
    Kind       SubmissionKind
    Text       string
    Status     SubmissionStatus
    StaffVotes map[string]bool // staff UUID -> voted, toggleable
    Reactions  map[string]int  // audience UUID -> +1 / -1, changeable, not stackable
    Created    time.Time
    Played     *time.Time
}

type Round struct {
    Number         int
    SubmissionID   string
    TeamOrder      []string
    PresentingTeam string
    Votes          map[string]string // audience UUID -> team ID, changeable until vote closes
    WinnerTeamIDs  []string          // ties allowed
}

type Game struct {
    Key                string
    GameCode           int
    Players            []*Player
    Teams              []*Team
    Submissions        []*Submission
    Rounds             []*Round
    SubmissionsEnabled bool
    ShowSubmitterNames bool
    Volume             int
    Created            time.Time
    Updated            time.Time
    Deleted            *time.Time
}

func (g *Game) GetStaff() []*Player
func (g *Game) GetDirectors() []*Player
func (g *Game) GetTeamDirectors(teamID string) []*Player
func (g *Game) GetCurrentRound() *Round
func (g *Game) NextTopic() *Submission // most StaffVotes among Approved, oldest first on ties
```

### Game codes

Trudge drives the UI with numeric `GameCode`s (20 lobby, 30 started, 600 end, ...). Fantasy
keeps that mechanism so the admin buttons, `CODE <n>` socket message and `gameCode` props look
the same, but as named constants in `models/game_code.go` with an allowed-transition table
checked in `SetGameCodeProcessor`:

| Code | Constant                       | Phase          |
|------|--------------------------------|----------------|
| 20   | `GAME_CODE_LOBBY`              | Lobby          |
| 30   | `GAME_CODE_INTAKE`             | Intake         |
| 100  | `GAME_CODE_TOPIC_REVEAL`       | Topic Reveal   |
| 110  | `GAME_CODE_RESPONSES`          | Responses      |
| 120  | `GAME_CODE_AUDIENCE_VOTE`      | Audience Vote  |
| 130  | `GAME_CODE_ROUND_RESULT`       | Round Result   |
| 600  | `GAME_CODE_FINAL_TALLY`        | Final Tally    |

Rules enforced in processors (and covered by tests):

- Votes and reactions are maps keyed by voter UUID, so no one can vote twice. (Trudge
  increments counters with no per-voter record, so repeat votes double-count, and its `Votes`
  padding path appends nil maps that would panic if written.)
- Reactions are accepted only at `GAME_CODE_RESPONSES`; team votes only at
  `GAME_CODE_AUDIENCE_VOTE`; only `Audience` may react or vote.
- Starting a round needs at least 2 teams, each with at least one Director.
- Teams cannot be added or removed mid-round.

---

## 5. Allowed origins (configurable, with "be naughty")

New settings in `config-fantasy.toml`, read in `setup/cors.go`. One function,
`setup.OriginAllowed(c fiber.Ctx) bool`, is used by **both** the HTTP CORS middleware and the
websocket pre-upgrade check, so they cannot drift apart:

```toml
# Explicit allow-list, used when CORS_BE_NAUGHTY is false
CORS_ALLOW_ORIGINS = ["http://localhost:5173", "https://fantasy.example"]

# Be naughty: allow any Origin whose host matches the Host the request arrived on.
# Handy behind changing tunnels/LAN IPs at a convention. Not for public production.
CORS_BE_NAUGHTY = false
```

- `CORS_BE_NAUGHTY = false`: Fiber's `middleware/cors` with `AllowOrigins` set from
  `CORS_ALLOW_ORIGINS` (exact match), `AllowCredentials: true`.
- `CORS_BE_NAUGHTY = true`: a small `middleware/naughty_cors.go` runs first. It parses the
  `Origin` header, compares its host (and port) with `c.Hostname()`, and on a match echoes that
  exact origin in `Access-Control-Allow-Origin` with `Vary: Origin`, then answers preflight
  `OPTIONS` with 204. Never `*`. Requests that don't match fall through to the allow-list
  middleware. (Fiber's `AllowOriginsFunc` only receives the origin string, so it can't compare
  against the request host.)
- `c.Hostname()` only honours `X-Forwarded-Host` when `TRUST_PROXY = true` and the peer is in
  `TRUSTED_PROXIES`; otherwise it is the raw `Host` header.
- The websocket pre-upgrade middleware calls the same `OriginAllowed`, so naughty mode applies
  to sockets too. A missing Origin on a websocket upgrade is rejected (browsers always send one).
- Both can be overridden by env vars (`FANTASY_CORS_BE_NAUGHTY=true`) through viper's
  `AutomaticEnv`, so a deploy can flip it without editing the file.
- Logs a `logrus.Warn` at startup when naughty mode is on.

Trudge hard-codes its origin list in `cmd/trudge.go` and accepts any websocket origin;
Fantasy does neither. Unit tests cover: allow-listed origin, unlisted origin, naughty match,
naughty mismatch, spoofed `X-Forwarded-Host` with `TRUST_PROXY` off, and websocket upgrade in
each mode.

---

## 6. API and socket messages

Same split as Trudge: REST for commands (`router.Get/Post/Put` in `api_router.go` returning
`P.RunController(c, &http2.XxxController{})`), websocket for broadcasts.

| Method | Path                                     | Role           | Controller                        |
|--------|------------------------------------------|----------------|-----------------------------------|
| GET    | /api/token/refresh                       | any            | RefreshTokenController            |
| GET    | /api/game                                | any player     | GetGameController (role-filtered) |
| POST   | /api/game/create                         | anyone         | CreateGameController (creator = Admin) |
| POST   | /api/game/join                           | anyone         | JoinGameController                |
| GET    | /api/game/rejoin                         | any player     | RejoinGameController              |
| GET    | /api/game/quit                           | any player     | QuitGameController                |
| POST   | /api/game/code                           | Admin          | SetGameCodeController             |
| GET    | /api/game/kick/:uuid                     | Admin          | KickPlayerController              |
| POST   | /api/game/role/assign                    | Admin          | AssignRoleController (+ TeamID)   |
| POST   | /api/game/team/add                       | Admin          | AddTeamController                 |
| POST   | /api/game/team/change/:id                | Admin          | ChangeTeamController              |
| POST   | /api/game/team/remove/:id                | Admin          | RemoveTeamController              |
| POST   | /api/game/submissions                    | Admin          | ToggleSubmissionsController       |
| POST   | /api/game/submitter-names                | Admin          | ToggleSubmitterNamesController    |
| POST   | /api/game/round/start                    | Admin          | StartRoundController (top item or by ID) |
| POST   | /api/game/round/presenting               | Admin          | SetPresentingTeamController       |
| PUT    | /api/game/submission/new                 | Audience       | NewSubmissionController           |
| POST   | /api/game/submission/status/:id          | Staff, Admin   | ChangeSubmissionStatusController  |
| POST   | /api/game/submission/vote/:id            | Staff          | StaffVoteController               |
| POST   | /api/game/submission/pin                 | Admin          | PinTopicController (override next topic) |
| POST   | /api/game/round/react                    | Audience       | ReactTopicController              |
| POST   | /api/game/round/vote                     | Audience       | VoteTeamController                |
| POST   | /api/game/volume                         | Admin, Display | SyncVolumeController              |

Role checks happen in the processor/controller from `helpers.GetRequestPlayer(c)` against the
**stored** player, not the JWT copy, so a role change takes effect immediately.

Websocket: `GET /ws/:room/:player` as in Trudge, plus `?token=<jwt>`. The pre-upgrade middleware
verifies the token, checks it matches `:player` and checks the origin, all before
`websocket.New(...)` upgrades the connection. On connect the server sends `CODE <n>` and the settings messages, as Trudge does
with `SUGGESTIONS_ENABLED`.

Messages keep Trudge's envelope `{"Command": "...", "Data": {...}}`, one file per message in
`app/messages`, sent with `wsCollection.HasTag(string(models.PLAYER_CLASS_X))`:

`NEW_PLAYER`, `PLAYER_LEFT`, `ROLE_ASSIGNED`, `TEAMS_UPDATED`, `SUBMISSIONS_ENABLED`,
`SUBMITTER_NAMES`, `NEW_SUBMISSION` (Staff, Admin), `SUBMISSION_UPDATED` (Staff, Admin),
`TOPIC_REVEALED`, `PRESENTING_TEAM`, `REACTION_TALLY` (throttled), `VOTE_TALLY`
(Admin, Display, throttled), `ROUND_RESULT`, `FINAL_TALLY`, `SET_VOLUME`.

---

## 7. Frontend views

| View      | Built from (Trudge)            | Shows                                                                         |
|-----------|--------------------------------|-------------------------------------------------------------------------------|
| Login     | `Login.vue`                    | Room code + name; create game                                                 |
| Admin     | `Leader.vue`, `gm/*`           | Game-code buttons (legal next steps only), `PlayerAssignment` with role + team, team editor (2-3), submission override, submissions toggle, **submitter names toggle**, live vote/reaction tallies, volume |
| Display   | `Display.vue`, `display/*`     | Lobby + room code, topic card (author if enabled), presenting-team banner, reaction meter (Responses only), vote bars for N teams, round winner reveal, final podium (team + submission). Reuse `Sunburst`, `DisplayMessage`, `GameEnd` patterns and music |
| Staff     | `God.vue`, `gods/*`            | Queue tabs (Pending / Approved / Rejected), approve/reject, upvote toggle, sorted by votes, "next up" preview |
| Director  | `PartyMember.vue`, `party/*`   | Team name and teammates, current topic, "your team is presenting" cue, scoreboard |
| Audience  | `Audience.vue`, `audience/*`   | Submission form (kind + text) while intake open (from `NewSuggestion.vue`); ±1 during Responses; team vote during Audience Vote (from `VoteOption.vue`); otherwise `PleaseWait` |

---

## 8. Things deliberately not carried over from Trudge

- Committed session keys in `config-trudge.toml`: Fantasy reads keys from env vars and fails
  fast in production when they are missing (dev still auto-generates, as `setup/keys.go` does).
- Hard-coded CORS list and `CheckOrigin: return true`: replaced by section 5.
- Websocket auth by UUID in the path alone: token now required too.
- `MemoryGameRepository` map access outside the lock, and `GetByKey` before `Lock` in
  processors: Fantasy keeps the same API, but `GetByKey`/`Save`/`Delete` take the repository
  mutex, and processors lock before mutating. `go test -race` enforces this.
- Characters, races, equipment, health, dice, maps, map maker: not part of this game.
- `vendor/`: not needed with the module proxy; add later only if offline Docker builds need it.

---

## 9. Delivery phases and branches

Integration branch: **`feature/fantasy`** (off `main`). Phase branches come off it and merge
back. `main` only receives `feature/fantasy` once the game is playable end-to-end.

| # | Branch                              | Scope                                                                                     | Done when |
|---|-------------------------------------|-------------------------------------------------------------------------------------------|-----------|
| 0 | `feature/fantasy`                   | This plan                                                                                 | Plan agreed |
| 1 | `feature/fantasy-backend-scaffold`  | New module; copy trudge skeleton (cmd, setup, bindings, routers, sockets, middleware, repositories, services) onto Fiber v3 + latest libs; Fiber safety config, copying `fiberhelper` accessors and tests (section 2a); configurable CORS + naughty mode; websocket with PING/PONG; create/join/rejoin/quit/refresh/kick/role; Dockerfile + compose | `go vet` + `go test -race` green; trudge still builds with `go build -mod=vendor ./apps/trudge/...` |
| 1 | `feature/fantasy-frontend-scaffold` | create-vue on latest, reshaped to trudge layout; Bootstrap 5, FontAwesome Free; port TokenManager, WebSocket, Login, Game.vue role switch | `npm run build`, `lint`, `test:unit` green |
| 2 | `feature/fantasy-models`            | Models, game codes + transitions, NextTopic, tally/scoring with ties, unit tests          | Full round simulated in a Go test |
| 3 | `feature/fantasy-gameplay-api`      | Teams, submissions, staff triage, rounds, reactions, votes, toggles; processors + messages + controllers | httptest coverage per endpoint/role |
| 4 | `feature/fantasy-views-core`        | Admin, Audience, Staff views                                                              | Manual loop: submit → triage → topic |
| 5 | `feature/fantasy-views-stage`       | Display and Director views, reveals for 2 and 3 teams, final tally, music/volume          | Full playthrough on 5 browser windows, with 2 and with 3 teams |
| 6 | `feature/fantasy-hardening`         | Snapshot/restore, reconnection, load test (~200 audience sockets), text length + rate limits, README, deploy config | Rehearsal game survives a backend restart |

Phases 1 (backend) and 1 (frontend) run in parallel; 4 and 5 can run in parallel once 3 lands.

---

## 10. Decisions log

| # | Question                    | Decision                                                                   |
|---|-----------------------------|----------------------------------------------------------------------------|
| 1 | Scoring                     | One point per round won; ties share the win; total votes break final ties |
| 2 | Team count                  | 2 typical, 3 supported; N-team data model, `MAX_TEAMS` default 3           |
| 3 | Submitter names             | Admin toggle, default off; server-side filtering; admin/staff always see   |
| 4 | Reaction window             | Only while teams are presenting (`GAME_CODE_RESPONSES`)                    |
| 5 | Icons                       | Free sets only; FontAwesome Free, keeping Trudge's registration pattern    |
| 6 | DI and code style           | Keep `go-di-container`; mirror Trudge structure and style, backend and frontend |
| 7 | Origins                     | Config allow-list + `CORS_BE_NAUGHTY` host-matching mode, shared by CORS and websocket |
| 8 | HTTP framework              | Fiber v3 instead of gin; websocket behaviour kept; risks in section 2a mitigated by config + tests |
| 9 | Request data                | Always copied (`strings.Clone`/`bytes.Clone` via `fiberhelper`), never referenced; `Immutable` left off so tests catch misses |

---

## 11. Implementation notes (as built)

Status: phases 0-5 are implemented on `feature/fantasy` (backend scaffold, models, gameplay
API, frontend views). Phase 6 (hardening) is partly done: snapshots, reconnects and limits
are in; a load test and deployment config for the real domain are still to do.

Where the build differs from, or adds to, the plan:

- **Frontend state** lives in a Pinia store (`stores/game.ts`) that handles every socket
  message, instead of Trudge's chain of events from each view up to `Game.vue`. Components
  are still Options API, now typed (`<script lang="ts">` + `defineComponent`).
- **Team management** is one `TeamProcessor` (add / change / remove) rather than three.
- **GAME_UPDATE** sends each player their own view whenever the game moves on; frequent
  tallies (reactions, votes) are throttled deltas.
- **Socket rooms**: messages never register a room (`sockets.GetCollection`), and joining /
  leaving a room is atomic (`JoinCollection` / `LeaveCollection`). Trudge's version leaks
  an entry per room and can orphan a socket that joins as the last one leaves.
- **Token encryption** is AES-GCM (Trudge: AES-CFB, deprecated in Go).
- **Dev port** is 18090 (8090 is in a Windows reserved range on the dev machine).
- **Music**: the display supports per-phase tracks like Trudge. New, convention-themed
  tracks are still to be chosen; Trudge's are not reused (`frontend/fantasy/src/libs/music.ts`).
- **Endpoints** are build configuration only: `VITE_API_URL` / `VITE_WS_URL` must be set
  (production: `https://ffc-api.furrific.uk/api`, `wss://ffc-api.furrific.uk/ws`, see
  `.env.production.example`) or the build fails. The frontend origin goes in the backend's
  `CORS_ALLOW_ORIGINS`.
- **Dependencies** are on the latest releases, except TypeScript stays on 6.x: vue-tsc and
  typescript-eslint do not support TypeScript 7 yet. `npm audit` reports `braces` (via
  `@vue/eslint-config-typescript` and `fast-glob`), a lint-time-only dev dependency with no
  patched release; the app bundle does not include it. Node must be 24.15+ (or 22.22+).
- **Admin top picks, veto and override**: stars are staff-only. The admin sees the staff's
  top three approved submissions (play order) and can *Play now*, *Play next* (pin: the
  pinned submission is revealed next, ahead of the staff ranking, even if still pending) or
  *Veto*. An admin rejection is a veto: staff cannot restore it, only an admin can. Playing,
  vetoing or resetting clears the pin.
- **CI/CD**: `.github/workflows/fantasy.yml` tests both halves, builds the backend image to
  `ghcr.io/<owner>/fantasy-backend` and publishes the frontend build to the
  `deploy/fantasy-frontend` branch, on pushes to `main` or manual runs. Actions are pinned to
  commit SHAs. `backend/Dockerfile-fantasy.dockerignore` limits the image build context to
  Fantasy without affecting the Trudge or PriceIsRight images.
