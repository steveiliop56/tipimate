# --- BUILDER ----
FROM golang:1.23-alpine3.20 AS builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY internal/ internal/
COPY cmd/ cmd/

RUN go mod tidy

RUN CGO_ENABLED=0 go build -o tipimate -ldflags "-s -w"

# --- RUNNER ----
FROM alpine:3.20 AS runner

WORKDIR /tipimate

RUN mkdir /data

COPY --from=builder /build/tipimate /tipimate

ENV DB_PATH=/data/tipimate.db

ENTRYPOINT ["/tipimate/tipimate", "server"]