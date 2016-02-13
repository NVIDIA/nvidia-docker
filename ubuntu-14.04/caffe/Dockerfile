FROM cuda:7.0-cudnn4-runtime
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

ENV CAFFE_VERSION 0.14
LABEL com.nvidia.caffe.version="0.14"

ENV CAFFE_PKG_VERSION 0.14.2-1
RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
            caffe-nv=$CAFFE_PKG_VERSION \
            caffe-nv-tools=$CAFFE_PKG_VERSION \
            python-caffe-nv=$CAFFE_PKG_VERSION && \
    rm -rf /var/lib/apt/lists/*
