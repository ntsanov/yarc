#/bin/bash
docker run --rm -ti \
    -v ${PWD}/cli:/workspace \
    harmonyone/main make -C /workspace static
    