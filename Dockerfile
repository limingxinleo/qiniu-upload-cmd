FROM golang:1.14 as builder

LABEL maintainer="limx <l@hyperf.io>"
ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE=on

WORKDIR /go/cache

ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /go/builder

ADD . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -o app main.go