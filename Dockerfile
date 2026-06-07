FROM golang:1.23.2-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build \
    -ldflags="-s -w" \
    -o main \
    ./cmd/server

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8000

CMD ["/app/main"]