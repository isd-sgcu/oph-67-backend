FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/main.go

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/server .
RUN apk --no-cache add ca-certificates
EXPOSE 4000
CMD ["./server"]