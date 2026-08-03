# Build stage
FROM golang:1.25.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o server \
    ./cmd/server


# Runtime stage
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/server .

COPY migrations ./migrations

EXPOSE 5040

CMD ["./server"]