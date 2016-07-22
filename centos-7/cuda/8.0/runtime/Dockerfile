FROM centos:7
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

LABEL com.nvidia.volumes.needed="nvidia_driver"

ENV CUDA_VERSION 8.0
LABEL com.nvidia.cuda.version="8.0"

ENV CUDA_DOWNLOAD_SUM bbecae635acf9f3c136e6d84523e59e14984ad2bd2fcb876ab49e027c2c331ae

ENV CUDA_PKG_VERSION 8-0-8.0.27-1
RUN curl -o cuda-repo.rpm -fsSL http://developer.download.nvidia.com/compute/cuda/8.0/direct/cuda-repo-rhel7-8-0-rc-8.0.27-1.x86_64.rpm && \
    echo "$CUDA_DOWNLOAD_SUM  cuda-repo.rpm" | sha256sum -c --strict - && \
    rpm -i cuda-repo.rpm && \
    rm cuda-repo.rpm && \
    yum install -y \
        cuda-nvrtc-$CUDA_PKG_VERSION \
        cuda-nvgraph-$CUDA_PKG_VERSION \
        cuda-cusolver-$CUDA_PKG_VERSION \
        cuda-cublas-$CUDA_PKG_VERSION \
        cuda-cufft-$CUDA_PKG_VERSION \
        cuda-curand-$CUDA_PKG_VERSION \
        cuda-cusparse-$CUDA_PKG_VERSION \
        cuda-npp-$CUDA_PKG_VERSION \
        cuda-cudart-$CUDA_PKG_VERSION && \
    yum remove -y cuda-repo-rhel7-8-0-rc.x86_64 && \
    ln -s cuda-$CUDA_VERSION /usr/local/cuda && \
    rm -rf /var/cache/yum/*

RUN echo "/usr/local/cuda/lib" >> /etc/ld.so.conf.d/cuda.conf && \
    echo "/usr/local/cuda/lib64" >> /etc/ld.so.conf.d/cuda.conf && \
    ldconfig

RUN echo "/usr/local/nvidia/lib" >> /etc/ld.so.conf.d/nvidia.conf && \
    echo "/usr/local/nvidia/lib64" >> /etc/ld.so.conf.d/nvidia.conf

ENV PATH /usr/local/nvidia/bin:/usr/local/cuda/bin:${PATH}
ENV LD_LIBRARY_PATH /usr/local/nvidia/lib:/usr/local/nvidia/lib64:${LD_LIBRARY_PATH}
