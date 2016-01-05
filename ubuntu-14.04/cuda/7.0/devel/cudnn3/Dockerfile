FROM cuda:7.0-devel
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

RUN apt-get update && apt-get install -y \
        curl && \
    rm -rf /var/lib/apt/lists/*

ENV CUDNN_VERSION 3
LABEL com.nvidia.cudnn.version="3"

ENV CUDNN_DOWNLOAD_SUM 98679d5ec039acfd4d81b8bfdc6a6352d6439e921523ff9909d364e706275c2b

RUN curl -fsSL http://developer.download.nvidia.com/compute/redist/cudnn/v3/cudnn-7.0-linux-x64-v3.0-prod.tgz -O && \
    echo "$CUDNN_DOWNLOAD_SUM  cudnn-7.0-linux-x64-v3.0-prod.tgz" | sha256sum -c --strict - && \
    tar -xzf cudnn-7.0-linux-x64-v3.0-prod.tgz -C /usr/local && \
    rm cudnn-7.0-linux-x64-v3.0-prod.tgz && \
    ldconfig
