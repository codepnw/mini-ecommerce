# Builder
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api/main.go

# Runner
FROM alpine:latest

# tzdata manage timezone
RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# copy binary from builder
COPY --from=builder /app/main .

USER appuser

EXPOSE 8080

CMD [ "./main" ]
