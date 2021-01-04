FROM golang:1.15.6-buster

WORKDIR /go/src/app

RUN apt-get update && apt-get install -y \
        git \
        vim \
        net-tools 

RUN go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate
ENV POSTGRESQL_URL 'postgres://user1:password@store:5432/app?sslmode=disable'