FROM cuda:8.0-devel
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

ENV CUDNN_VERSION 5
LABEL com.nvidia.cudnn.version="5"

ENV CUDNN_DOWNLOAD_SUM af80eb1ce0cb51e6a734b2bdc599e6d50b676eab3921e5bddfe5443485df86b6

RUN curl -fsSL http://developer.download.nvidia.com/compute/redist/cudnn/v5/cudnn-8.0-linux-x64-v5.0-ga.tgz -O && \
    echo "$CUDNN_DOWNLOAD_SUM  cudnn-8.0-linux-x64-v5.0-ga.tgz" | sha256sum -c --strict - && \
    tar -xzf cudnn-8.0-linux-x64-v5.0-ga.tgz -C /usr/local && \
    rm cudnn-8.0-linux-x64-v5.0-ga.tgz && \
    ldconfig
