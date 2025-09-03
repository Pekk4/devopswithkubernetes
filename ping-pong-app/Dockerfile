FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o pingpong_app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/pingpong_app .

EXPOSE 3000

CMD ["./pingpong_app"]
