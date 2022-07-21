FROM golang:1.18-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/go-sample-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.* ./

RUN go mod tidy

COPY . .

# Unit tests
#RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go build -o auth-sv




# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates

COPY --from=build_base /tmp/go-sample-app/auth-sv /app/go-app

# Set the Current Working Directory inside the container
WORKDIR /app

# This container exposes port 8081 to the outside world
EXPOSE 50092

# Run the binary program produced by `go install`
CMD ["./go-app"]
