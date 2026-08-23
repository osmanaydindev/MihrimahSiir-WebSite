# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

MihrimahSiir — a poetry/literature social platform. Go + Fiber backend (`backend/`), Vue 3 + Vuetify SPA (`frontend/`), PostgreSQL. UI text, comments, and commit messages are in Turkish; keep that convention when editing existing files.

## Commands

```bash
# Backend (from backend/)
go run main.go                    # dev server, :8080 — reads backend/.env via godotenv
go build -o bin/main main.go
go vet ./...

# Frontend (from frontend/)
npm install
npm run dev                       # Vite dev server, :3000
npm run build                     # vue-tsc --noEmit && vite build (type errors fail the build)
npm run lint                      # eslint . --fix

# Docker (from repo root)
docker-compose -f docker-compose.yml -f docker-compose.local.yml up -d --build
docker-compose logs -f backend
./scripts/deploy-vps.sh           # VPS deploy (docker-compose.vps.yml)
./release.sh {major|minor|patch} {backend|frontend}   # git tag only, e.g. backend-0.1.0
```

There is no frontend test setup and only two Go test files (`security/isbn_test.go`, `util/jwt_roundtrip_test.go`, run with `go test ./...`). Most verification is done by running the servers, not by tests — don't claim coverage you don't have.

`make` targets are broken — the `Makefile` does `include infra/Makefile` and `infra/` does not exist. Use the docker-compose commands directly.

## Environment files

None are committed (`.gitignore` blocks `*.env` except `*.env.example`). Copy from examples before running:

- `backend/.env` — `DB_HOST/DB_USER/DB_PASSWORD/DB_NAME/DB_PORT/DB_SSLMODE`, `APP_PORT`, `CORS_ORIGINS` (comma-separated, must include the frontend origin or every request fails), `COOKIE_SECURE`, `ADMIN_*`, `JWT_SECRET` (required, ≥32 chars, `openssl rand -hex 32`).
- `backend/env.local` / `backend/env.vps` — used by `docker-compose.local.yml` / `docker-compose.vps.yml` respectively via `env_file`. Neither is tracked, so secrets live in them directly.
- Mail/external: `MAIL_ENABLED`, `RESEND_API_KEY`, `MAIL_FROM`, `MAIL_REPLY_TO`, `ADMIN_NOTIFY_EMAIL`, `MAIL_DAILY_LIMIT`, `APP_PUBLIC_URL`, `OPENLIBRARY_BASE_URL`. With `MAIL_ENABLED=false` or an empty API key the mailer renders to stdout instead of sending — that's the local-dev default.
- `frontend/.env` — `VITE_API_BASE_URL`. In Docker this is a **build arg** baked into the bundle at image build time; changing it requires a frontend rebuild, not a restart.

## Backend architecture

`main.go` wires, in strict order: `SecurityHeaders()` → `GlobalRateLimiter()` (100 req/min per IP) → CORS → `routes.Setup(app)`. Adding middleware before `SecurityHeaders` breaks the header guarantees.

`routes/routes.go` is the auth boundary. Everything registered **before** `app.Use(middlewares.IsAuthenticated)` is public (`/register`, `/login` — both behind `AuthRateLimiter`, 5 per 15 min — plus `/ws` and `/uploads`). Everything after it requires a valid `token` cookie. New public endpoints must be registered above that line; new protected ones below. Per-route admin checks use `middlewares.IsAdmin`.

**Auth is cookie-based, not header-based.** `Login` sets an HTTPOnly `token` cookie (5h, `SameSite=None`, `Secure` from `COOKIE_SECURE`). The frontend sends it via `withCredentials: true`; there is no `Authorization` header path (the commented-out variant in `middlewares/IsAuthenticated.go` is dead code). The JWT subject (`Issuer`) is the numeric admin ID as a string.

**JWT signing key comes from `JWT_SECRET`** (HS256, `util/jwt.go`). `util.InitJWT()` runs in `main.go` before anything else and **panics if the var is missing or shorter than 32 chars** — a misconfigured deploy fails loudly instead of silently minting sessions nobody can validate. Changing `JWT_SECRET` logs every user out, so treat it as a fixed per-environment value. `util/jwt_roundtrip_test.go` covers the restart-survival guarantee (`go test ./util/`).

**Roles** live on `Admin.RoleID`: `1` = admin, `2` = member, `3` = registered user. Authorization inside controllers is done by calling `helpers.GetUserRole(c)` and comparing, not by a middleware — follow that pattern for new admin-only handlers.

**Visibility (`Community`) differs between poems and books, deliberately.** Poems have `1` = private (roles 1 & 2) and `2` = public, filtered by `applyCommunityFilter` in `controllers/PoemController.go`; `security.ValidateCommunity` enforces 1|2 and must stay that way, because the poem filter only matches `= 2` and a `3` would make a poem permanently invisible. Books add `3` = "only selected users", backed by the `book_visibilities` join table and validated by the separate `ValidateBookCommunity`. Every book query must go through `helpers.ApplyBookCommunityFilter(db, roleID, userID)` — it needs the **user** id, not just the role, and its SQL qualifies columns as `books.*` so it works in JOINed queries too. Writes go through `helpers.SetBookVisibility`, which clears the list whenever the level is not 3.

There is one `Admin` model for all users — regular registrants are `Admin` rows with `RoleID = 3`. Schema changes go through `database.ConnectDb()`'s `AutoMigrate` list; a new model is invisible until it is added there. There are no migration files.

List endpoints paginate through `helpers.GetPaginationParams` (`?offset=&limit=`, limit capped at 100, default 20) and return `helpers.CreatePaginationResponse` (`{data, total, offset, limit, has_more}`).

WebSocket (`/ws`, `websocket/`): `GlobalHub` runs as a goroutine; the client authenticates with `?token=` as a **query param** (the cookie is not readable there). Server-side pushes go through `GlobalHub.SendToUser(userID, msgType, payload)` — fire-and-forget, dropped if the user is offline, so never the only delivery path for something that matters.

## External integrations (`services/`)

Two outbound integrations live under `backend/services/`, the only place in the backend that makes HTTP calls. Both share the same discipline: one package-level `http.Client` singleton (never per-request), `context` timeouts, `io.LimitReader` on the body, and draining the body before close so connections are reused.

- **`services/openlibrary`** — `Default().FetchByISBN(ctx, isbn13)`. One primary call to `/api/books?jscmd=data` (title, authors, pages, cover, publisher, date) plus two best-effort calls for the work-level description, which may arrive as a plain string *or* `{"value": …}`. Total budget 10s, process-wide concurrency capped at 4. A missing cover is left empty on purpose — the `/b/isbn/…` fallback returns a blank 1×1 GIF rather than a 404. Note Open Library often returns a *plausible wrong book* for a mistyped ISBN, so the admin UI shows ISBN/publisher/date for sanity-checking and there is no one-click approve.
- **`services/mail`** — Resend HTTP API (not SMTP). `mail.StartWorker()` in `main.go` runs a single consumer goroutine over a buffered channel; handlers call `mail.Enqueue`/`mail.Notify*` with **plain value structs — never `*fiber.Ctx`**, which is pooled and reused after the handler returns. A failed or skipped email never fails the user's request. `models.MailLog` backs a restart-proof daily quota (`MAIL_DAILY_LIMIT`, default 80, under Resend's free ~100/day) and doubles as an audit trail. Templates are `html/template`; never wrap user- or API-supplied values in `template.HTML`.

## Book requests (ISBN → Open Library → admin approval)

Users submit an ISBN, the metadata snapshot is fetched and stored **at request time** (not at approval — the admin list would otherwise cost N external calls, and Open Library can return different editions on different calls). `controllers/BookRequestController.go` + `routes/bookRequestRoutes.go`; `BookRequest.Status` is `pending|approved|rejected`, following the `Friendship` string-status idiom.

Because each request costs an external call plus an email, the abuse defense is layered and ordered cheapest-first — only a request that survives all of it opens a socket or queues mail:

1. `middlewares.UserRateLimiter` (5/hour, keyed by user id via `helpers.GetUserIDFromCtx`, falling back to IP). `SkipFailedRequests` is on, so typos don't burn the quota — only successful requests count.
2. ISBN checksum + ISBN-10→13 normalization (`security.NormalizeISBN`). Dedup always happens on the normalized form, or the same book arrives twice under two ISBNs.
3. DB guards: already in library / already requested / rejected within 30 days (the cooldown is what stops a reject→re-request loop from draining the mail quota).
4. DB caps: 3 concurrent pending, 10 per 24h.

The rate limiter's store is in-memory and resets on restart; the DB-backed checks carry the real weight. `uniq_book_requests_pending_isbn` is a **partial** unique index (`WHERE status = 'pending'`) created in `database/indexes.go` — GORM tags can't express it, and a plain `(isbn, status)` unique would wrongly block a second rejected row.

Approval runs in one transaction (create book → set visibility → mark request), with mail and WebSocket pushes deliberately **outside** it. The requester is always auto-added to a `community = 3` visibility list, otherwise they couldn't see the book they asked for.

Input handling: `security/sanitizer.go`, `security/validator.go`, `security/file_validator.go`.

**Uploads are written with paths relative to the process working directory** (`./uploads/profiles`, `app.Static("/uploads", "./uploads")`). The backend Dockerfile therefore sets `WORKDIR /app` to match the `/app/uploads` volume mount in every compose file — changing either side detaches uploaded files from the volume and they vanish on the next `docker-compose up --build`. `main.go` creates `./uploads/profiles` at startup so a freshly mounted empty volume still works. If a path is ever made absolute, change all three places together.

## Frontend architecture

Vite is configured with heavy auto-import plugins (`vite.config.mts`): `unplugin-auto-import` (Vue APIs, `useRoute`/`useRouter`), `unplugin-vue-components` over an explicit `dirs` list, and Vuetify auto-import. Components in a directory not listed in that `dirs` array will not resolve in templates — add the directory when creating a new component folder. The generated `auto-imports.d.ts`, `components.d.ts`, and `typed-router.d.ts` are build artifacts; don't hand-edit them.

**Routing is the manual `src/router/index.ts`**, not file-based, despite `unplugin-vue-router` being installed. Add new pages there. The file is `// @ts-nocheck`. Its global `beforeEach` guard hits `/auth-check` on **every** `requiresAuth` navigation, pushes the returned user into the Pinia store, and enforces admin-only pages via a hardcoded `forbiddenRoutes` name list — a new admin page must be added to that array or any logged-in user can reach it.

**Two store directories exist and both are live**: `src/store/app.ts` is the real `useAppStore` (user, liked/bookmarked poem ID sets, book state) imported as `@/store/app`; `src/stores/index.ts` is the Pinia instance registered in `plugins/index.ts`, and `src/stores/app.ts` is an empty stub. Use `@/store/app` for state; don't add anything to `src/stores/app.ts`.

Liked/bookmarked/read state is held as `Set<number>` of IDs in the store so cards can render their state without per-item requests — after a like/bookmark/read mutation, update the corresponding set (`addLikedPoemId` / `removeLikedPoemId`, etc.) rather than refetching.

**Two API styles coexist.** The typed service layer (`src/services/api/*.ts` on a shared `apiClient` axios instance with interceptors) is used only by the composables in `src/composables/`; most `.vue` files still call the global `axios` directly, configured in `plugins/index.ts` (`baseURL` from `VITE_API_BASE_URL`, `withCredentials = true`). Prefer the service layer + composables for new code. Gotcha: `import … from '../services/api'` resolves to the legacy `src/services/api.ts`, **not** `src/services/api/index.ts` — a new service must be re-exported from both barrels or the build fails at bundle time (type-checking passes).

Composables carry the reusable behavior — `usePoemActions`, `useBookActions`, `useFriendActions`/`useFriendsList`, `useLazyLoad` (pairs with the backend's offset/limit pagination), `useNotification`, `useConfirmDialog`, `useErrorHandler`, `useSanitizer` (DOMPurify; use it for any user-supplied HTML, poem/comment content is Quill-authored rich text).

Two layouts: `layouts/default.vue` (user site) and `layouts/adminPanel.vue` (`/panel` and `/*-management` routes).

`components/shared/Table/Table.vue` is **not** a generic table: it switches on `route.fullPath` and hardcodes `<td>`s per management page. New admin pages are usually better off with their own `v-data-table` — `pages/Management/BookRequestManagement.vue` does that. Its `imagePath()` now routes through `utils/imageHelper.getImageUrl`, so absolute URLs (e.g. Open Library covers) survive instead of being prefixed with the local `/uploads` host.

## Known rough edges

- `scripts/deploy.sh` and `scripts/deploy-vps.sh` health-check `http://localhost:8080/health`, but no `/health` route exists — that check always reports failure.
- `models.Author` is missing from the `AutoMigrate` list. The `authors` table exists anyway (created transitively via `Book.AuthorData`), and adding it is risky rather than free: `Author.Slug` is `unique;not null`, so one duplicate or empty slug in prod would fail the migration and `ConnectDb` **panics**, taking the deploy down. Fix it in its own commit after checking prod data.
- `UserReadBookController.GetReadsBooksPaginated` applies no visibility filter at all, so a private/selected-only book stays fully visible once it's on a user's read list. May be intentional (they already had it) — confirm before changing.
- `middlewares/IsAuthenticated.go` returns **HTTP 200** with `{"message":"unauthenticated"}` on failure instead of 401, so frontend code can't rely on the axios error path for auth failures.
- `Book.Slug` has no unique constraint and `GetBook` looks it up with `First()`. New writes go through `helpers.UniqueBookSlug`, but pre-existing duplicates are still unreachable by slug.
- `github.com/dgrijalva/jwt-go` is unmaintained; `github.com/golang-jwt/jwt/v5` is the drop-in successor if the dependency is ever refreshed.
- The JWT itself expires after 24h (`util.SetToken`) while the login cookie expires after 5h (`controllers.Login`) — the cookie is the effective session length.
