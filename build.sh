#!/usr/bin/env bash


DATESTR=`date -Iseconds`
GITSTR=`git rev-parse HEAD`
BUILD_VER=${DATESTR}_${GITSTR}_release
echo "Build" ${BUILD_VER}


BuildBin() {
    local srcfile=${1}
    local dstdir=${2}
    local dstfile=${3}
    local args="-X main.Ver=${BUILD_VER}"

    echo "[BuildBin] go build -i -o ${dstdir}/${dstfile} -ldflags "${args}" ${srcfile}"

    mkdir -p ${dstdir}
    go build -i -o ${dstdir}/${dstfile} -ldflags "${args}" ${srcfile}

    if [ ! -f "${dstdir}/${dstfile}" ]; then
        echo "${dstdir}/${dstfile} build fail, build file: ${srcfile}"
        # exit 1
    fi
    strip "${dstdir}/${dstfile}"
}


cd lib
genlog -leveldatafile ./goslog/goslog.data -packagename goslog 
cd ..


ProtocolGOSFiles="protocol_gos/*.enum"

PROTOCOL_GOS_VERSION=`cat ${ProtocolGOSFiles}| sha256sum | awk '{print $1}'`
echo "Protocol GOS Version:" ${PROTOCOL_GOS_VERSION}

cd protocol_gos
genprotocol -ver=${PROTOCOL_GOS_VERSION} \
    -basedir=. \
    -prefix=gos -statstype=int

goimports -w .

cd ..

# generate enum
# genenum -typename=TeamType -packagename=teamtype -basedir=enum -vectortype=int

cd enum 
goimports -w .
cd .. 

GameDataFiles="
config/gameconst/gameconst.go \
config/gameconst/serviceconst.go \
config/gamedata/*.go \
enum/*.enum \
"
Data_VERSION=`cat ${GameDataFiles}| sha256sum | awk '{print $1}'`
echo "Data Version:" ${Data_VERSION}

echo "
package gameconst

const DataVersion = \"${Data_VERSION}\"
" > config/gameconst/dataversion_gen.go 

# build bin

BIN_DIR="bin"
SRC_DIR="rundriver"

mkdir -p ${BIN_DIR}
echo ${BUILD_VER} > ${BIN_DIR}/BUILD

BuildBin ${SRC_DIR}/server.go ${BIN_DIR} server
BuildBin ${SRC_DIR}/multiclient.go ${BIN_DIR} multiclient

cd rundriver
./genwasmclient.sh
cd ..

# copy data to bin folder
echo cp -r rundriver/serverdata ${BIN_DIR}
cp -r rundriver/serverdata ${BIN_DIR}
echo cp -r rundriver/clientdata ${BIN_DIR}
cp -r rundriver/clientdata ${BIN_DIR}

