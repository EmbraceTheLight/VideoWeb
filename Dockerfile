FROM golang:alpine AS builder
LABEL authors="zey"
WORKDIR /VideoWeb

ENV GOPROXY https://goproxy.cn,direct
ADD go.mod .
ADD go.sum .
COPY . .
RUN  chmod 777 ./wait-for-it.sh && go mod download && go mod tidy && go build -o main  main.go

FROM ubuntu:24.04
WORKDIR /VideoWeb

# 设置非交互式环境
ENV DEBIAN_FRONTEND=noninteractive
RUN yes|apt-get update && yes|apt-get upgrade && yes|apt-get install ffmpeg

COPY --from=builder /VideoWeb/main /VideoWeb/main
COPY --from=builder /VideoWeb/resources /VideoWeb/resources
COPY --from=builder /VideoWeb/config /VideoWeb/config
COPY --from=builder /VideoWeb/wait-for-it.sh /VideoWeb/wait-for-it.sh

