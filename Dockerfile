ARG GO_VERSION=1.16
ARG GOLANGCI_LINT_VERSION=v1.40.1

FROM golang:${GO_VERSION} AS base
WORKDIR /app
RUN apt-get install \
    make

FROM golangci/golangci-lint:${GOLANGCI_LINT_VERSION} AS lint-base

FROM base AS lint
COPY --from=lint-base /usr/bin/golangci-lint /usr/bin/golangci-lint
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    make -f rules/builder.mk check


FROM base as make-ctl

COPY --from=lint-base /usr/bin/golangci-lint /usr/bin/golangci-lint

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make APP_EXECUTABLE_OUT=/out -f rules/builder.mk

FROM scratch as ctl
COPY --from=make-ctl /out/* .

FROM base AS make-go-mod-tidy
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod tidy

FROM scratch AS go-mod-tidy
COPY --from=make-go-mod-tidy /app/go.mod .
COPY --from=make-go-mod-tidy /app/go.sum .
