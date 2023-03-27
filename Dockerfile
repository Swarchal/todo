FROM golang:1.20-buster

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY *.go ./

ADD css /app/css
ADD templates /app/templates

RUN go build

EXPOSE 3333
