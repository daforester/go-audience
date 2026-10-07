# Fantasy Furry Convention (frontend)

Vue 3 client for the Fantasy Furry Convention game. The backend lives in
`backend/apps/fantasy`; the overall plan is in `docs/fantasy/PLAN.md`.

## Roles

One page serves every role, chosen by the admin in the player list:

| Role         | Screen                                                                        |
|--------------|-------------------------------------------------------------------------------|
| Admin        | Game controls, teams, player roles, the submission queue                      |
| Main Display | The projector: lobby, topic reveal, reactions, live vote, results, awards    |
| Staff Team   | Approve / reject audience submissions and star favourites                     |
| Con Director | Team card, the topic, a cue when the team is presenting, the scoreboard       |
| Audience     | Send submissions, thumbs up / down the topic, vote for the best team         |

## Development

The quickest start is the dev compose file in the repository root, which runs the backend
(rebuilt on change by air) and this dev server (hot reload) from your checkout:

```sh
docker compose -f docker-compose-fantasy-dev.yml up
```

Then open http://localhost:5173 (API on http://localhost:18090). Stop with
`docker compose -f docker-compose-fantasy-dev.yml down` (add `-v` to clear caches and
saved games).

Without Docker, run the backend with the testing config (API on port 18090):

```sh
cd backend/apps/fantasy
go run ./cmd/fantasy.go -config ../../config-testing-fantasy.toml
```

Then the frontend:

```sh
npm install
npm run dev
```

Each browser profile holds one player, so use separate profiles (or a private window)
to play several roles at once.

## Scripts

| Command              | What it does                         |
|----------------------|--------------------------------------|
| `npm run dev`        | Vite dev server                      |
| `npm run build`      | Type-check and production build      |
| `npm run test:unit`  | Vitest unit tests                    |
| `npm run lint`       | oxlint and ESLint                    |

## Configuration

API endpoints are never hardcoded. Every build needs `VITE_API_URL` (http/https) and
`VITE_WS_URL` (ws/wss), as absolute URLs, and the build fails without them. Provide them
as environment variables or in `.env.[mode].local` (git ignored):

```sh
cp .env.production.example .env.production.local   # https://ffc-api.furrific.uk
npm run build
```

or, for CI:

```sh
VITE_API_URL=https://ffc-api.furrific.uk/api VITE_WS_URL=wss://ffc-api.furrific.uk/ws npm run build
```

`.env.development` holds the local dev backend address used by `npm run dev`.

The backend must list the site the frontend is served from in `CORS_ALLOW_ORIGINS`
(or `FANTASY_CORS_ALLOW_ORIGINS`).

## Deployment

`.github/workflows/fantasy.yml` tests and builds both halves on every pull request that
touches Fantasy. On a push to `main` (or a manual "Run workflow"), it also:

- pushes the backend image to `ghcr.io/<owner>/fantasy-backend`, tagged `latest`, the
  branch name and the short commit SHA;
- publishes this frontend's production build to the `deploy/fantasy-frontend` branch, one
  commit per deploy, with `DEPLOYED_FROM.txt` recording the source commit. Serve that
  branch's files from any static host; roll back by deploying an earlier commit of it.

The build reads the endpoints from repository variables `FANTASY_API_URL` and
`FANTASY_WS_URL` (Settings > Secrets and variables > Actions > Variables), and fails if
they are missing.

## Music

The main display can play music per phase like Trudge's did. New tracks suited to the
convention theme are still to be chosen (Trudge's fantasy tracks are deliberately not
reused). Add files to `src/assets/music` and list them in `src/libs/music.ts`.
