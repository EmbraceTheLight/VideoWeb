FROM golang:alpine AS builder
LABEL authors="zey"
WORKDIR /VideoWeb

ENV GOPROXY https://goproxy.cn,direct
ADD go.mod .
ADD go.sum .
COPY . .
RUN  chmod 777 ./wait-for-it.sh && go mod download && go mod tidy && go build -o main  main.go

FROM centos
WORKDIR /VideoWeb

COPY --from=builder /VideoWeb/main /VideoWeb/main
COPY --from=builder /VideoWeb/wait-for-it.sh /VideoWeb/wait-for-it.sh
#CMD ["./main"]
