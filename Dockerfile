FROM golang:1.20-alpine as builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy

COPY . .
RUN go build -o main ./cmd

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 50053

CMD ["./main"]
