# builder
FROM registry.hoven.top/library/golang:1.24.3-alpine AS builder

WORKDIR /app
COPY ./ /app

ENV GO111MODULE=auto
ENV GOINSECURE="github.com/go-puzzles"
ENV GOPROXY="https://goproxy.cn,direct"

RUN cd /app && \
	go mod tidy && \
	go build -o ./server && \
	chmod +x server 

# runner
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/server /app/server
