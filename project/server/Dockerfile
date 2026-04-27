FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache linux-headers

COPY go.mod ./
RUN go get github.com/glebarez/sqlite@v1.10.0
RUN go get github.com/gin-gonic/gin@v1.9.1

COPY main.go ./

RUN CGO_ENABLED=0 go build -o server main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server .

RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./server"]