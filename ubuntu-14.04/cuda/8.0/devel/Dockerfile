FROM cuda:8.0-runtime
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

RUN curl -o cuda-repo.deb -fsSL http://developer.download.nvidia.com/compute/cuda/8.0/direct/cuda-repo-ubuntu1404-8-0-rc_8.0.27-1_amd64.deb && \
    echo "$CUDA_DOWNLOAD_SUM  cuda-repo.deb" | sha256sum -c --strict - && \
    dpkg -i cuda-repo.deb && \
    rm cuda-repo.deb && \
    apt-get update && apt-get install -y --no-install-recommends \
        cuda-core-$CUDA_PKG_VERSION \
        cuda-misc-headers-$CUDA_PKG_VERSION \
        cuda-command-line-tools-$CUDA_PKG_VERSION \
        cuda-nvrtc-dev-$CUDA_PKG_VERSION \
        cuda-nvml-dev-$CUDA_PKG_VERSION \
        cuda-nvgraph-dev-$CUDA_PKG_VERSION \
        cuda-cusolver-dev-$CUDA_PKG_VERSION \
        cuda-cublas-dev-$CUDA_PKG_VERSION \
        cuda-cufft-dev-$CUDA_PKG_VERSION \
        cuda-curand-dev-$CUDA_PKG_VERSION \
        cuda-cusparse-dev-$CUDA_PKG_VERSION \
        cuda-npp-dev-$CUDA_PKG_VERSION \
        cuda-cudart-dev-$CUDA_PKG_VERSION \
        cuda-driver-dev-$CUDA_PKG_VERSION && \
    apt-get remove --purge -y cuda-repo-ubuntu1404-8-0-rc && \
    rm -rf /var/lib/apt/lists/*

ENV LIBRARY_PATH /usr/local/cuda/lib64/stubs:${LIBRARY_PATH}
