#!/bin/bash
APP=${APP:-atlas}
OS=${OS:-windows freebsd linux}
BUCKET=${BUCKET:-}

mkdir -p build

for v in $OS; do
    dir=$(mktemp -d)
    root=$(pwd)
    rm -rf ./bin/*
    echo " -> building ${v}"
    make GOOS=${v} && cd ${dir} && tar czf ${root}/build/${APP}-${v}-latest.tar.gz ${root}/bin/
    cd ${root}
    rm -rf ${dir}
done

s3cmd sync ./build/* s3://${BUCKET}/
s3cmd setacl --recursive s3://${BUCKET}/ --acl-public
