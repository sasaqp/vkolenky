FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev libc-dev

COPY go.mod ./
RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/gin-gonic/gin

COPY main.go ./

RUN CGO_ENABLED=1 go build -o server main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates musl-dev

WORKDIR /app

COPY --from=builder /app/server .

RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./server"]