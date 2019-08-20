#!/bin/bash
APP=${APP:-atlas}
OS=${OS:-windows freebsd linux}

mkdir -p build

for v in $OS; do
    root=$(pwd)
    rm -rf ./bin/*
    echo " -> building ${v}"
    make GOOS=${v} && cd bin && zip -D -r ${root}/build/${APP}-${v}-latest.zip .
    cd ${root}
done
