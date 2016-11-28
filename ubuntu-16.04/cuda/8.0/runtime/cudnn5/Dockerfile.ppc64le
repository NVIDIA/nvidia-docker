FROM cuda:8.0-runtime
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

RUN apt-get update && apt-get install -y \
        curl && \
    rm -rf /var/lib/apt/lists/*

ENV CUDNN_VERSION 5
LABEL com.nvidia.cudnn.version="5"

RUN CUDNN_DOWNLOAD_SUM=51f698d468401cef2e3e2ef9bb557bd57cbeb4dca895d1d1ae8a751d090bbe39 && \
    curl -fsSL http://developer.download.nvidia.com/compute/redist/cudnn/v5.1/cudnn-8.0-linux-ppc64le-v5.1.tgz -O && \
    echo "$CUDNN_DOWNLOAD_SUM  cudnn-8.0-linux-ppc64le-v5.1.tgz" | sha256sum -c --strict - && \
    tar -xzf cudnn-8.0-linux-ppc64le-v5.1.tgz -C /usr/local --wildcards 'cuda/targets/ppc64le-linux/lib/libcudnn.so.*' && \
    ln -s cuda/targets/ppc64le-linux/lib /usr/local/cuda/lib64 && rm cudnn-8.0-linux-ppc64le-v5.1.tgz && \
    ldconfig
