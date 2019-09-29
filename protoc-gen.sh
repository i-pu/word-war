#!/bin/bash

set -ex

protoc --proto_path proto proto/sample.proto \
--go_out=plugins=grpc:server/pb \
--js_out=import_style=commonjs:client/src \
--grpc-web_out=import_style=commonjs,mode=grpcwebtext:client/src
