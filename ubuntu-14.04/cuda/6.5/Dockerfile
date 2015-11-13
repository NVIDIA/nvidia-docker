FROM ubuntu:14.04
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

RUN apt-get update && apt-get install -y wget

ENV NVIDIA_GPGKEY_SUM 4e181b4141cf3daca6f68e34d5387499fd093b9632b522329b10ed4d1929d417

RUN wget -q http://developer.download.nvidia.com/compute/cuda/repos/GPGKEY && \
    echo "$NVIDIA_GPGKEY_SUM GPGKEY" | sha256sum -c - && \
    apt-key add GPGKEY && rm GPGKEY && \
    echo "deb http://developer.download.nvidia.com/compute/cuda/repos/ubuntu1404/x86_64 /" > /etc/apt/sources.list.d/cuda.list

ENV CUDA_VERSION 6.5
LABEL com.nvidia.cuda.version="6.5"

RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
    cuda-toolkit-6-5 \
    && rm -rf /var/lib/apt/lists/*

RUN echo "/usr/local/cuda/lib" >> /etc/ld.so.conf.d/cuda.conf && \
    echo "/usr/local/cuda/lib64" >> /etc/ld.so.conf.d/cuda.conf && \
    ldconfig

RUN echo "/usr/local/nvidia/lib" >> /etc/ld.so.conf.d/nvidia.conf && \
    echo "/usr/local/nvidia/lib64" >> /etc/ld.so.conf.d/nvidia.conf

ENV PATH /usr/local/cuda/bin:${PATH}
ENV LD_LIBRARY_PATH /usr/local/nvidia/lib:/usr/local/nvidia/lib64:${LD_LIBRARY_PATH}
