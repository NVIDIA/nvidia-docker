FROM cuda:7.5-devel
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

ENV CUDNN_VERSION 4
LABEL com.nvidia.cudnn.version="4"

ENV CUDNN_DOWNLOAD_SUM 4e64ef7716f20c87854b4421863328e17cce633330c319b5e13809b61a36f97d

RUN curl -fsSL http://developer.download.nvidia.com/compute/redist/cudnn/v4/cudnn-7.0-linux-x64-v4.0-prod.tgz -O && \
    echo "$CUDNN_DOWNLOAD_SUM  cudnn-7.0-linux-x64-v4.0-prod.tgz" | sha256sum -c --strict - && \
    tar -xzf cudnn-7.0-linux-x64-v4.0-prod.tgz -C /usr/local && \
    rm cudnn-7.0-linux-x64-v4.0-prod.tgz && \
    ldconfig
