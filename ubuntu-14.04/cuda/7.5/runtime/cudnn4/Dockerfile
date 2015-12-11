FROM cuda:7.5-runtime
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

RUN apt-get update && apt-get install -y \
        curl && \
    rm -rf /var/lib/apt/lists/*

ENV CUDNN_VERSION 4
LABEL com.nvidia.cudnn.version="4"

ENV CUDNN_DOWNLOAD_SUM 4bcc67921018543e9066adb82ffd731f089eac80413e28eab6ae22a8bf5f4aaf

RUN curl -fsSL http://developer.download.nvidia.com/compute/redist/cudnn/v4/cudnn-7.0-linux-x64-v4.0-rc.tgz -O && \
    echo "$CUDNN_DOWNLOAD_SUM cudnn-7.0-linux-x64-v4.0-rc.tgz" | sha256sum -c --strict - && \
    tar -xzf cudnn-7.0-linux-x64-v4.0-rc.tgz --wildcards 'cuda/lib64/libcudnn.so*' -C /usr/local && \
    rm cudnn-7.0-linux-x64-v4.0-rc.tgz && \
    ldconfig
