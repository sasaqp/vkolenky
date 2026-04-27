FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=1 go build -o server .

FROM alpine:latest

RUN apk add --no-cache ca-certificates gcc musl-dev

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/data ./data

RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./server"]