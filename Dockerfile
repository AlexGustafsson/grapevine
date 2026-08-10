FROM --platform=${BUILDPLATFORM} node:26.5.1@sha256:a9875b5ccb02aa527cf7f2297b16ae425a0ff3da2f7d87fce3df41f04ffa0524 AS web-builder

WORKDIR /src

COPY .npmrc package.json package-lock.json .

RUN --mount=type=cache,target=node_modules  \
  npm ci

COPY tsconfig.json vite.config.ts .
COPY web web

RUN --mount=type=cache,target=node_modules \
  npm run build

FROM --platform=${BUILDPLATFORM} golang:1.26.5@sha256:3aff6657219a4d9c14e27fb1d8976c49c29fddb70ba835014f477e1c70636647 AS builder

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
