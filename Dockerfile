# Arguments
ARG TIPIMATE_VERSION

# Builder
FROM golang:1.23-alpine3.20 AS builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY internal/ internal/
COPY cmd/ cmd/

RUN go mod tidy

RUN go build -o tipimate -ldflags "-s -w -X tipimate/internal/constants.Version=${TIPIMATE_VERSION}"

# Runner
FROM alpine:3.20 AS runner

WORKDIR /tipimate

RUN mkdir /data

COPY --from=builder /build/tipimate /tipimate

ENV TIPIMATE_DATABASE_PATH=/data/tipimate.db

ENTRYPOINT ["/tipimate/tipimate", "server"]