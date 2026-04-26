ARG SERVICE

# builder
FROM golang:1.21-alpine AS builder

ARG SERVICE

RUN apk add --no-cache protoc curl
RUN curl -1sLf 'https://dl.cloudsmith.io/public/task/task/setup.alpine.sh' | sudo -E bash
RUN apk add task

WORKDIR /app

COPY go.mod go.sum Taskfile.yml ./
COPY proto/ ./proto/
RUN task generate
RUN go mod download

COPY ${SERVICE}/ ./service/
WORKDIR /app/${SERVICE}
RUN go build -o bin .

# runtime
FROM alpine:latest

ARG SERVICE

WORKDIR /root/
COPY --from=builder /app/service/bin .

EXPOSE 50051
CMD ["./bin"]