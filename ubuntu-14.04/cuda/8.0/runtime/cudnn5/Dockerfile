FROM cuda:8.0-runtime
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

RUN apt-get update && apt-get install -y \
        curl && \
    rm -rf /var/lib/apt/lists/*

ENV CUDNN_VERSION 5
LABEL com.nvidia.cudnn.version="5"

RUN CUDNN_DOWNLOAD_SUM=a87cb2df2e5e7cc0a05e266734e679ee1a2fadad6f06af82a76ed81a23b102c8 && \
    curl -fsSL http://developer.download.nvidia.com/compute/redist/cudnn/v5.1/cudnn-8.0-linux-x64-v5.1.tgz -O && \
    echo "$CUDNN_DOWNLOAD_SUM  cudnn-8.0-linux-x64-v5.1.tgz" | sha256sum -c --strict - && \
    tar -xzf cudnn-8.0-linux-x64-v5.1.tgz -C /usr/local --wildcards 'cuda/lib64/libcudnn.so.*' && \
    rm cudnn-8.0-linux-x64-v5.1.tgz && \
    ldconfig
