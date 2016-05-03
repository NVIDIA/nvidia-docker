FROM cuda:7.5-runtime
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

RUN yum install -y \
        cuda-core-$CUDA_PKG_VERSION \
        cuda-misc-headers-$CUDA_PKG_VERSION \
        cuda-command-line-tools-$CUDA_PKG_VERSION \
        cuda-license-$CUDA_PKG_VERSION \
        cuda-nvrtc-dev-$CUDA_PKG_VERSION \
        cuda-cusolver-dev-$CUDA_PKG_VERSION \
        cuda-cublas-dev-$CUDA_PKG_VERSION \
        cuda-cufft-dev-$CUDA_PKG_VERSION \
        cuda-curand-dev-$CUDA_PKG_VERSION \
        cuda-cusparse-dev-$CUDA_PKG_VERSION \
        cuda-npp-dev-$CUDA_PKG_VERSION \
        cuda-cudart-dev-$CUDA_PKG_VERSION \
        cuda-driver-dev-$CUDA_PKG_VERSION \
        gcc-c++ && \
    rm -rf /var/cache/yum/*

RUN mkdir /tmp/gpu-deployment-kit && cd /tmp/gpu-deployment-kit && \
    rpm2cpio $(repoquery --location  gpu-deployment-kit) | cpio -id && \
    mv usr/include/nvidia/gdk/* /usr/local/cuda/include && \
    mv usr/src/gdk/nvml/lib/* /usr/local/cuda/lib64/stubs && \
    rm -rf /tmp/gpu-deployment-kit* && \
    rm -rf /var/cache/yum/*

ENV LIBRARY_PATH /usr/local/cuda/lib64/stubs:${LIBRARY_PATH}
