FROM cuda:7.5-cudnn5-runtime
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

ENV CAFFE_VERSION 0.15
LABEL com.nvidia.caffe.version="0.15"

ENV CAFFE_PKG_VERSION 0.15.9-1+cuda7.5
RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
        caffe-nv=$CAFFE_PKG_VERSION \
        caffe-nv-tools=$CAFFE_PKG_VERSION \
        python-caffe-nv=$CAFFE_PKG_VERSION \
        libcaffe-nv0=$CAFFE_PKG_VERSION && \
    rm -rf /var/lib/apt/lists/*
