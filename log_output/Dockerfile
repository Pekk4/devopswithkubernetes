FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o log_app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/log_app .

EXPOSE 3000

CMD ["./log_app"]
