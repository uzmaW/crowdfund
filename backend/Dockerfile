FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o backend .


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/backend .

# Copy migrations folder for database migrations
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./backend"]
