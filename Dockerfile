# builder
FROM golang:1.24.0-alpine AS builder

WORKDIR /app
COPY ./ /app

ENV GO111MODULE=auto
ENV GOINSECURE="gitea.hoven.com"
ENV GOPROXY="http://10.15.25.3:8081"

RUN apk add --no-cache git
RUN cd /app && \
	go mod tidy && \
	go build -o ./server && \
	chmod +x server 

# runner
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/server /app/server
