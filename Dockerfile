FROM golang:1.20-buster

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY *.go ./

COPY vendor ./vendor

ADD css /app/css
ADD templates /app/templates

RUN go build -mod vendor

EXPOSE 3333
