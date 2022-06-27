#!/bin/bash

protoc src/auth/auth.proto --go_grpc_out=plugins=grpc:.

protoc --go_out=. --go_opt=paths=source_relative  src/auth/auth.proto

protoc --go_out=. --go_opt=paths=source_relative \
     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
     src/auth/auth.proto
