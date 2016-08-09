FROM ubuntu:16.04
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

LABEL com.nvidia.volumes.needed="nvidia_driver"

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl && \
    rm -rf /var/lib/apt/lists/*

ENV CUDA_VERSION 8.0
LABEL com.nvidia.cuda.version="8.0"

ENV CUDA_DOWNLOAD_SUM 24278d78afed380b4328c1e2f917b31d70c3f4c8f297b642200e003311944c22

ENV CUDA_PKG_VERSION 8-0=8.0.27-1
RUN curl -o cuda-repo.deb -fsSL http://developer.download.nvidia.com/compute/cuda/8.0/direct/cuda-repo-ubuntu1604-8-0-rc_8.0.27-1_amd64.deb && \
    echo "$CUDA_DOWNLOAD_SUM  cuda-repo.deb" | sha256sum -c --strict - && \
    dpkg -i cuda-repo.deb && \
    rm cuda-repo.deb && \
    apt-get update && apt-get install -y --no-install-recommends \
        cuda-nvrtc-$CUDA_PKG_VERSION \
        cuda-nvgraph-$CUDA_PKG_VERSION \
        cuda-cusolver-$CUDA_PKG_VERSION \
        cuda-cublas-$CUDA_PKG_VERSION \
        cuda-cufft-$CUDA_PKG_VERSION \
        cuda-curand-$CUDA_PKG_VERSION \
        cuda-cusparse-$CUDA_PKG_VERSION \
        cuda-npp-$CUDA_PKG_VERSION \
        cuda-cudart-$CUDA_PKG_VERSION && \
    ln -s cuda-$CUDA_VERSION /usr/local/cuda && \
    apt-get remove --purge -y cuda-repo-ubuntu1604-8-0-rc && \
    rm -rf /var/lib/apt/lists/*

RUN echo "/usr/local/cuda/lib" >> /etc/ld.so.conf.d/cuda.conf && \
    echo "/usr/local/cuda/lib64" >> /etc/ld.so.conf.d/cuda.conf && \
    ldconfig

RUN echo "/usr/local/nvidia/lib" >> /etc/ld.so.conf.d/nvidia.conf && \
    echo "/usr/local/nvidia/lib64" >> /etc/ld.so.conf.d/nvidia.conf

ENV PATH /usr/local/nvidia/bin:/usr/local/cuda/bin:${PATH}
ENV LD_LIBRARY_PATH /usr/local/nvidia/lib:/usr/local/nvidia/lib64:${LD_LIBRARY_PATH}
