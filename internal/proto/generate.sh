#!/bin/bash

protoc internal/proto/ --go_grpc_out=plugins=grpc:.

protoc --go_out=. --go_opt=paths=source_relative  internal/proto/*.proto

protoc --go_out=. --go_opt=paths=source_relative \
     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
     internal/proto/*.proto

