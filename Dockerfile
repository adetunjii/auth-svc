#FROM golang:1.18-alpine AS build_base

#RUN apk add --no-cache git

# Set the Current Working Directory inside the container
#WORKDIR /tmp/go-sample-app

# We want to populate the module cache based on the go.{mod,sum} files.
#COPY go.* ./

#RUN go mod tidy

#COPY . .

# Unit tests
#RUN CGO_ENABLED=0 go test -v

# Build the Go app
#RUN go build -o .




# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates

#COPY --from=build_base /tmp/go-sample-app/checklos /app/go-app

# Set the Current Working Directory inside the container
WORKDIR /app


COPY auth-sv /app/auth-sv

# Copy env File
# COPY .env.example .env

# This container exposes port 8081 to the outside world
EXPOSE 50092

# Run the binary program produced by `go install`
CMD ["./auth-sv"]
