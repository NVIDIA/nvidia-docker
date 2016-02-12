FROM cuda:7.5-runtime
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

RUN echo "deb http://developer.download.nvidia.com/compute/machine-learning/repos/ubuntu1404/x86_64 /" > /etc/apt/sources.list.d/nvidia-ml.list

ENV CUDNN_VERSION 4
LABEL com.nvidia.cudnn.version="4"

RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
            libcudnn4=4.0.7 && \
    rm -rf /var/lib/apt/lists/*
