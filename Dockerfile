FROM golang:1.22-alpine3.18 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/main.go

FROM alpine:3.18
WORKDIR /app

RUN apk --no-cache add ca-certificates
COPY --from=builder /app/server .

EXPOSE 4000

CMD ["./server"]