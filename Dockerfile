# syntax=docker/dockerfile:1

FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/epicpaste
COPY . .
RUN go mod tidy
RUN mkdir /server
RUN go build -o /server/epicpaste main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
RUN mkdir /server
COPY --from=builder /server/epicpaste  /server/.env /server/
COPY --from=builder /server/sql/*  /server/sql/
CMD ["/server/epicpaste"]