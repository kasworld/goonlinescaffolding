#!/usr/bin/env bash

$BUILD_VER=$args[0]

# build gl client

Remove-Item -ErrorAction SilentlyContinue wasmclient.wasm

Write-Output "GOOS=js GOARCH=wasm go build -o wasmclient.wasm -ldflags -X main.Ver=${BUILD_VER}"
$env:GOOS="js" 
$env:GOARCH="wasm" 
go build -o wasmclient.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclient.go
$env:GOOS=""
$env:GOARCH=""

Write-Output "move files"
Move-Item -Force wasmclient.wasm clientdata/
