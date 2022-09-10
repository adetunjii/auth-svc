FROM golang:1.19-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/go-sample-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.* ./

RUN go mod tidy

COPY . .

# Unit tests
#RUN CGO_ENABLED=0 go test -v

# install curl 
RUN apk add curl

# download the golang-migrate package and unzip
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

# Build the Go app
RUN go build -o auth-sv



# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates

COPY --from=build_base /tmp/go-sample-app/auth-sv /app/go-app
COPY --from=build_base /tmp/go-sample-app/migrate /app/migrate

# Set the Current Working Directory inside the container
WORKDIR /app

# COPY app.env.example app.env

# This container exposes port 8081 to the outside world
EXPOSE 50092

# Run the binary program produced by `go install`
CMD ["./go-app"]
