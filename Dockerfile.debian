ARG VERSION_ID
FROM debian:${VERSION_ID}

# packaging dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
        dh-make \
        fakeroot \
        build-essential \
        devscripts && \
    rm -rf /var/lib/apt/lists/*

# packaging
ARG PKG_VERS
ARG PKG_REV
ARG RUNTIME_VERSION
ARG DOCKER_VERSION

ENV DEBFULLNAME "NVIDIA CORPORATION"
ENV DEBEMAIL "cudatools@nvidia.com"
ENV REVISION "$PKG_VERS-$PKG_REV"
ENV RUNTIME_VERSION $RUNTIME_VERSION
ENV DOCKER_VERSION $DOCKER_VERSION
ENV SECTION ""

# output directory
ENV DIST_DIR=/tmp/nvidia-docker2-$PKG_VERS
RUN mkdir -p $DIST_DIR /dist

# nvidia-docker 2.0
COPY nvidia-docker $DIST_DIR/nvidia-docker
COPY daemon.json $DIST_DIR/daemon.json

WORKDIR $DIST_DIR
COPY debian ./debian

RUN sed -i "s;@VERSION@;${REVISION#*+};" debian/changelog && \
    if [ "$REVISION" != "$(dpkg-parsechangelog --show-field=Version)" ]; then exit 1; fi

CMD export DISTRIB="unstable" && \
    debuild --preserve-env --dpkg-buildpackage-hook='sh debian/prepare' -i -us -uc -b && \
    mv /tmp/*.deb /dist
