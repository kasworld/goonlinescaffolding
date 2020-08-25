
################################################################################
Set-Location lib
Write-Output "genlog -leveldatafile ./goslog/goslog.data -packagename goslog "
genlog -leveldatafile ./goslog/goslog.data -packagename goslog 
Set-Location ..

################################################################################
$PROTOCOL_GOS_VERSION=makesha256sum protocol_gos/*.enum protocol_gos/gos_obj/protocol_*.go
Write-Output "Protocol GOS Version: ${PROTOCOL_GOS_VERSION}"
Write-Output "genprotocol -ver=${PROTOCOL_GOS_VERSION} -basedir=protocol_gos -prefix=gos -statstype=int"
genprotocol -ver="${PROTOCOL_GOS_VERSION}" -basedir=protocol_gos -prefix=gos -statstype=int
Set-Location protocol_gos
goimports -w .
Set-Location ..

################################################################################
# generate enum
Write-Output "generate enums"

Set-Location enum
goimports -w .
Set-Location ..

$Data_VERSION=makesha256sum config/gameconst/*.go config/gamedata/*.go enum/*.enum
Write-Output "Data Version: ${Data_VERSION}"
mkdir -ErrorAction SilentlyContinue config/dataversion
Write-Output "package dataversion
const DataVersion = `"${Data_VERSION}`" 
" > config/dataversion/dataversion_gen.go 


################################################################################
$DATESTR=Get-Date -UFormat '+%Y-%m-%dT%H:%M:%S%Z:00'
$GITSTR=git rev-parse HEAD
$BUILD_VER="${DATESTR}_${GITSTR}_release_windows"
Write-Output "Build Version: ${BUILD_VER}"

################################################################################
# build bin

$BIN_DIR="bin"
$SRC_DIR="rundriver"

mkdir -ErrorAction SilentlyContinue "${BIN_DIR}"
Write-Output ${BUILD_VER} > ${BIN_DIR}/BUILD_windows

# build bin here
go build -o "${BIN_DIR}\server.exe" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\serverwin.go"
go build -o "${BIN_DIR}\multiclient.exe" -ldflags "-X main.Ver=${BUILD_VER}" "${SRC_DIR}\multiclient.go"

Set-Location rundriver
./genwasmclient.ps1 ${BUILD_VER}
Set-Location ..

Write-Output "cp -r rundriver/serverdata ${BIN_DIR}"
Copy-Item -Force -r rundriver/serverdata ${BIN_DIR}
Write-Output "cp -r rundriver/clientdata ${BIN_DIR}"
Copy-Item -Force -r rundriver/clientdata ${BIN_DIR}

