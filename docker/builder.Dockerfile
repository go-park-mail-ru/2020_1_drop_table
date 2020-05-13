FROM golang:1.13

RUN apt update && apt upgrade && \
    apt --update add git make

WORKDIR /app

COPY . .

RUN make build
