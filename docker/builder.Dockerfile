FROM golang:1.13

WORKDIR /application

COPY . .

RUN make build
