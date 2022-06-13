# syntax=docker/dockerfile:1

FROM golang:1.18

WORKDIR /build
COPY . /build

RUN go build -o /docker-todo

EXPOSE 8080

CMD [ "/docker-todo" ]