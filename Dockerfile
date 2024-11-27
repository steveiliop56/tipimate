# --- BUILDER ----
FROM golang:1.23-alpine3.20 AS builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY internal/ internal/
COPY cmd/ cmd/

RUN go mod tidy

RUN go build -ldflags "-s -w" -o tipicord main.go

# --- RUNNER ----
FROM alpine:3.20 AS runner

WORKDIR /tipicord

RUN mkdir /data

COPY --from=builder /build/tipicord /tipicord

ENV DATABASE_PATH=/data/tipicord.db

ENTRYPOINT ["/tipicord/tipicord"]