# Multi-stage build, non-root runtime user, health check on /healthz
# (docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md section 5's Docker checklist).
# Not wired into any CI/CD pipeline; build and run manually or via your own compose file.

FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# CGO disabled: the binary only needs to talk to Postgres over the wire
# protocol (pgx is pure Go) and shell out to pg_dump, never link against
# libpq, so a static binary keeps the runtime image minimal.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/openschool-api ./cmd/api

FROM alpine:3.20
# postgresql17-client provides pg_dump for the nightly backup agent
# (internal/modules/automation/system_health.go) — bump the version here if
# docker-compose.yml's postgres image version changes.
RUN apk add --no-cache postgresql17-client ca-certificates wget && \
    addgroup -S openschool && adduser -S openschool -G openschool
WORKDIR /home/openschool
COPY --from=build /out/openschool-api /usr/local/bin/openschool-api
# The nightly backup agent writes here (JOB_BACKUP_DIR); mount a volume over
# it if running the container with --read-only.
RUN mkdir -p /home/openschool/backups && chown openschool:openschool /home/openschool/backups
USER openschool
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/healthz || exit 1
ENTRYPOINT ["/usr/local/bin/openschool-api"]
