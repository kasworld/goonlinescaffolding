#!/usr/bin/env bash


rm wasmclient.wasm

echo "GOOS=js GOARCH=wasm go build -o wasmclient.wasm"
GOOS=js GOARCH=wasm go build -o wasmclient.wasm wasmclient.go

echo "mv wasmclient.wasm ./clientdata/"
mv wasmclient.wasm ./clientdata/
