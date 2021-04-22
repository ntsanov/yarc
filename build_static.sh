#/bin/bash

export VERSION=$(git rev-list --count HEAD)
export COMMIT=$(git describe --always --long --dirty)
docker run --rm -ti \
    -e VERSION=${VERSION} \
    -e COMMIT=${COMMIT} \
    -v ${PWD}/cli:/workspace \
    harmonyone/main make -C /workspace static
    