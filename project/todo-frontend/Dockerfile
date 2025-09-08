FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o todo_app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo_app .
COPY ./templates ./templates

EXPOSE 3000

CMD ["./todo_app"]
