FROM caffe:0.14
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

ENV DIGITS_VERSION 3.0
LABEL com.nvidia.digits.version="3.0"

ENV DIGITS_PKG_VERSION 3.0.0-1
# workaround: gcc and libhdf5-dev are dependencies that are currently missing from the torch package
RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
            torch7-nv=0.9.98-1+cuda7.5 \
            graphviz \
            gcc \
            libhdf5-dev \
            digits=$DIGITS_PKG_VERSION && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/log/digits
RUN mkdir -p /usr/share/digits/digits/jobs

EXPOSE 34448
WORKDIR /usr/share/digits
ENTRYPOINT ["./digits-server"]
