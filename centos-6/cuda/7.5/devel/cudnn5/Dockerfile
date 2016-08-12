FROM cuda:7.5-devel
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

ENV CUDNN_VERSION 5
LABEL com.nvidia.cudnn.version="5"

ENV CUDNN_DOWNLOAD_SUM 40d506d0a8a00a3faccce1433346806b8cd2535683b6f08a63683ce6e474419f

RUN curl -fsSL http://developer.download.nvidia.com/compute/redist/cudnn/v5.1/cudnn-7.5-linux-x64-v5.1.tgz -O && \
    echo "$CUDNN_DOWNLOAD_SUM  cudnn-7.5-linux-x64-v5.1.tgz" | sha256sum -c - && \
    tar -xzf cudnn-7.5-linux-x64-v5.1.tgz -C /usr/local && \
    rm cudnn-7.5-linux-x64-v5.1.tgz && \
    ldconfig
