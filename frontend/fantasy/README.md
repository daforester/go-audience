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

Run the backend with the testing config (API on port 18090):

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

`VITE_API_URL` and `VITE_WS_URL` (see `.env.development` / `.env.production`) may be
absolute, or relative to the page's host. The production build uses `/api` and `/ws` on
the same host, which suits a reverse proxy in front of both.

## Music

The main display can play music per phase like Trudge's did, but no tracks are bundled.
Add files to `src/assets/music` and list them in `src/libs/music.ts`.
