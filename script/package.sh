#!/bin/bash
APP=${APP:-atlas}
OS=${OS:-windows freebsd linux}

DIFF=$(git diff --no-ext-diff)
if [ ! -z "$DIFF" ]; then
    echo "$DIFF"
fi

mkdir -p build

for v in $OS; do
    root=$(pwd)
    rm -rf ./bin/*
    echo " -> building ${v}"
    make GOOS=${v} && cd bin && zip -qD -r ${root}/build/${APP}-${v}-latest.zip .
    cd ${root}
done
