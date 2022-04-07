#!/bin/bash

#Generate protobuf file (if its not there already)
protoc --go_out=./proto --go_opt=paths=source_relative     --go-grpc_out=./proto --go-grpc_opt=paths=source_relative     --proto_path=../sntrippled/src/ripple/proto/org/xrpl/rpc/v1/ gossip_message.proto

