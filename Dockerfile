# syntax=docker/dockerfile:1

FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/epicpaste
COPY . .
RUN go mod tidy
RUN mkdir /server
RUN cp .env /server/ 
# RUN cp -r sql /server/sql
RUN go build -o /server/epicpaste main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /rmsubekti
COPY --from=builder /server .
ENTRYPOINT [ "/rmsubekti/epicpaste" ]
