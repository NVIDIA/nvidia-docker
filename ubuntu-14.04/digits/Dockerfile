FROM caffe
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

ENV DIGITS_VERSION 3.0
LABEL com.nvidia.digits.version="3.0"

ENV DIGITS_PKG_VERSION 3.0.0-1
RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
            torch7-nv \
            graphviz \
            digits=$DIGITS_PKG_VERSION && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/log/digits
RUN mkdir -p /usr/share/digits/digits/jobs

EXPOSE 34448
WORKDIR /usr/share/digits
ENTRYPOINT ["./digits-server"]
