FROM --platform=${BUILDPLATFORM} node:26.3.0@sha256:e3ffe0cbaeebdcddbfe1ee7bca9b564a92863a8386d5b99a3d72677b3667b61d AS web-builder

WORKDIR /src

COPY .npmrc package.json package-lock.json .

RUN --mount=type=cache,target=node_modules  \
  npm ci

COPY tsconfig.json vite.config.ts .
COPY web web

RUN --mount=type=cache,target=node_modules \
  npm run build

FROM --platform=${BUILDPLATFORM} golang:1.26.4@sha256:68cb6d68bed024785b69195b89af7ac7a444f27791435f98647edff595aa0479 AS builder

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
