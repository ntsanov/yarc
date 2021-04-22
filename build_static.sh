#/bin/bash

if [ -n "$TRAVIS_BUILD_DIR" ]; then
    CLI_DIR=${TRAVIS_BUILD_DIR}/cli
else
    CLI_DIR=${PWD}/cli 
fi

export VERSION=$(git rev-list --count HEAD)
export COMMIT=$(git describe --always --long --dirty)
docker run --rm -ti \
    -e VERSION=${VERSION} \
    -e COMMIT=${COMMIT} \
    -v ${CLI_DIR}:/workspace \
    harmonyone/main make -C /workspace static
    