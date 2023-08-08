FROM golang:alpine
WORKDIR /rmsubekti
COPY . .
RUN go mod tidy
RUN go build -o epicpaste
RUN rm -rf app client proto system main.go Dockerfile .gi* go.* .dock* .devcontainer .vscode TODO.md 

ENTRYPOINT [ "/rmsubekti/epicpaste" ]