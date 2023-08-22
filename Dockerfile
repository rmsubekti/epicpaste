# syntax=docker/dockerfile:1

FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/epicpaste
COPY . .
RUN go mod tidy
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init --parseDependency  --parseInternal --parseDepth 1  -g main.go
RUN mkdir /server
RUN cp .env /server/ 
RUN go build -o /server/epicpaste main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /rmsubekti
COPY --from=builder /server .
ENTRYPOINT [ "/rmsubekti/epicpaste" ]
