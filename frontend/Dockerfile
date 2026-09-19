# Multi-stage build: compile the SPA, then serve the static bundle with
# Nginx (docs/SECURITY_AND_PERFORMANCE_PLAYBOOK.md section 5's Docker
# checklist: multi-stage, non-root, health check on /healthz).
# Not wired into any CI/CD pipeline; build and run manually or via your own
# compose file. Could NOT be build-verified in this environment: the pinned
# packageManager (pnpm@11.12.0) fails corepack's signature check here — see
# the queued feedback report. Re-verify this Dockerfile once that's resolved.

FROM node:22-alpine AS build
WORKDIR /src
RUN corepack enable
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN corepack pnpm install --frozen-lockfile
COPY . .
RUN corepack pnpm build

FROM nginx:1.27-alpine
# nginx:alpine already ships an unprivileged "nginx" user; just hand it
# ownership of the paths it needs to write to as non-root.
RUN touch /var/run/nginx.pid && \
    chown -R nginx:nginx /var/cache/nginx /var/run/nginx.pid /usr/share/nginx/html
COPY --from=build /src/dist /usr/share/nginx/html
# deploy/nginx.conf.example (repo root) is the reference for the security
# headers, gzip and cache rules this image's Nginx config should apply —
# adapt it into a real nginx.conf here before deploying.
COPY nginx.conf /etc/nginx/conf.d/default.conf
USER nginx
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/healthz || exit 1
