FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app/

RUN CGO_ENABLED=0 go build -o ./bin/service ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/service /app/
COPY --from=builder /app/config/config.yaml /app/config.yaml