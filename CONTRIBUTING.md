# Contributing to OpenSchool

Thanks for considering a contribution. This guide covers how to propose
changes - for how to actually get the app running, see
[`docs/SETUP.md`](docs/SETUP.md) and [`docs/THUNDERID.md`](docs/THUNDERID.md)
(the one-time identity-provider setup).

By participating, you agree to follow our
[Code of Conduct](CODE_OF_CONDUCT.md).

## Before you start

- **Small fix or obvious bug?** Open a PR directly.
- **New feature or larger change?** Open an issue first to discuss the
  approach - saves everyone rework if the direction needs adjusting.
- **Security issue?** Don't open a public issue - see
  [`SECURITY.md`](SECURITY.md) for how to report it privately.

## Getting a dev environment running

Follow [`docs/SETUP.md`](docs/SETUP.md) end to end, starting with
[`docs/THUNDERID.md`](docs/THUNDERID.md). That covers prerequisites,
starting Postgres/ThunderID/backend/frontend, and signing in.

## Orienting yourself

- [`CLAUDE.md`](CLAUDE.md) - fast orientation: backend layering, frontend
  structure, data model.
- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) - the full picture.
- [`docs/adr/`](docs/adr/) - *why* behind non-obvious decisions; check here
  before "fixing" something that looks wrong but is deliberate.

## Making changes

**Backend (Go):**

- SQL lives in `backend/db/queries/*.sql`. After editing it, run `sqlc
  generate` from `backend/` - never hand-edit `backend/db/sqlc/`.
- Follow the existing layering for a feature module: `routes` →
  `handlers` → `services` → `repositories`.
- One clear doc-comment line per exported function/type - match the
  existing style rather than writing long comment blocks.
- Before submitting: `go build ./...` and `go vet ./...` must pass; run
  `staticcheck ./...` too if you have it installed (CI does).

**Frontend (React/TypeScript):**

- Admin CRUD pages follow one shared template - list, a modal form,
  confirm-delete - built from `src/components/common/`. Deviating from it
  should be a deliberate choice, not an accident.
- Before submitting: `pnpm build` and `pnpm lint` must pass.

## Commit and PR conventions

- Branch from `development`: `feature/your-feature-name`.
- All PRs target `development`, not `main`.
- Keep commits focused; explain *why* a change was made, not just what
  changed.
- Reference the issue a PR resolves, if there is one.

## Need help?

- [Open an issue](https://github.com/openschool-org/openschool/issues)
- [GitHub Discussions](https://github.com/openschool-org/openschool/discussions)
