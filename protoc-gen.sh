#!/bin/bash

set -ex

CLIENT_OUTDIR=client/src/pb
SERVER_OUTPUT_DIR=server/interface/rpc/pb

mkdir -p ${CLIENT_OUTDIR} ${SERVER_OUTPUT_DIR}

protoc \
--proto_path proto \
--go_out=plugins=grpc:${SERVER_OUTPUT_DIR} \
--js_out=import_style=commonjs:${CLIENT_OUTDIR} \
--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:${CLIENT_OUTDIR} \
proto/word_war.proto
