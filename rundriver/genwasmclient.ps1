
$BUILD_VER=$args[0]

$env:GOOS="js" 
$env:GOARCH="wasm" 
Write-Output "go build -o clientdata/wasmclient.wasm -ldflags `"-X main.Ver=${BUILD_VER}`" wasmclient.go"
go build -o clientdata/wasmclient.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclient.go
$env:GOOS=""
$env:GOARCH=""

