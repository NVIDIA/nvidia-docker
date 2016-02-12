FROM cuda:7.5-devel
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

RUN echo "deb http://developer.download.nvidia.com/compute/machine-learning/repos/ubuntu1404/x86_64 /" > /etc/apt/sources.list.d/nvidia-ml.list

ENV CUDNN_VERSION 4
LABEL com.nvidia.cudnn.version="4"

RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
            libcudnn4-dev=4.0.7 && \
    rm -rf /var/lib/apt/lists/*
