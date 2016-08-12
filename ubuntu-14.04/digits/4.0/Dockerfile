FROM caffe:0.15
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

ENV DIGITS_VERSION 4.0
LABEL com.nvidia.digits.version="4.0"

# workaround: gcc and libhdf5-dev are dependencies that are currently missing from the torch package
ENV DIGITS_PKG_VERSION 4.0.0-1
RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
            torch7-nv=0.9.99-1+cuda7.5 \
            graphviz \
            gcc \
            libhdf5-dev \
            digits=$DIGITS_PKG_VERSION && \
    rm -rf /var/lib/apt/lists/*

VOLUME /data
VOLUME /jobs

COPY digits.cfg digits/digits.cfg

EXPOSE 34448
WORKDIR /usr/share/digits
ENTRYPOINT ["./digits-server"]
