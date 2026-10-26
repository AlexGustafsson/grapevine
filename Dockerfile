FROM --platform=${BUILDPLATFORM} node:26.11.1@sha256:32fa97f3363975684b08bf4e8a68a47c7905175cc20275b50b974bbd02aba731 AS web-builder

WORKDIR /src

COPY .npmrc package.json package-lock.json .

RUN --mount=type=cache,target=node_modules  \
  npm ci

COPY tsconfig.json vite.config.ts .
COPY web web

RUN --mount=type=cache,target=node_modules \
  npm run build

FROM --platform=${BUILDPLATFORM} golang:1.27.2@sha256:820153096fda4de6ac1b7194d12e416b7dd8f4e9a33d17e0345b111b945162d8 AS builder

WORKDIR /src

# Use the toolchain specified in go.mod, or newer
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
  go mod download && go mod verify

COPY cmd cmd
COPY internal internal

COPY --from=web-builder /src/internal/web/public /src/internal/web/public

ARG TARGETARCH
ARG TARGETOS
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
  GOARCH=${TARGETARCH} GOOS=${TARGETOS} CGO_ENABLED=0 go build -a -ldflags="-s -w" -o grapevine cmd/grapevine/*.go

FROM scratch AS export

COPY --from=builder /src/grapevine grapevine

FROM export

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV PATH=/

ENTRYPOINT ["grapevine"]
