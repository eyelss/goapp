ARG SERVICE

# builder
FROM golang:1.26-alpine AS builder

ARG SERVICE

RUN apk add --no-cache curl protobuf
RUN curl -1sLf 'https://dl.cloudsmith.io/public/task/task/setup.alpine.sh' | /bin/sh
RUN apk add --no-cache task

WORKDIR /app

COPY go.mod go.sum Taskfile.yml ./
COPY proto/ ./proto/
RUN task install-tools
RUN task generate
RUN go mod download

COPY ${SERVICE}/ ./service/
COPY framework/ ./framework/
WORKDIR /app/service
RUN go build -o bin .

# runtime
FROM alpine:latest

ARG SERVICE
ENV SERVICE=${SERVICE}

WORKDIR /root/
COPY --from=builder /app/service/bin .

RUN if [ "$DEBUG"="true" ]; then \
      apk add --no-cache curl nano vim; \
      apk add --no-cache --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing grpcurl; \
    fi

EXPOSE ${PORT}
CMD ["./bin"]