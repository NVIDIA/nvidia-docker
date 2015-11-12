FROM ubuntu:14.04
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

# setup keys
RUN apt-key adv --keyserver ha.pool.sks-keyservers.net --recv-keys 889BEE522DA690103C4B085ED88C3D385C37D3BE

# setup sources.list
RUN echo "deb http://developer.download.nvidia.com/compute/cuda/repos/ubuntu1404/x86_64 /" > /etc/apt/sources.list.d/cuda.list

# set cuda version
ENV CUDA_VERSION 7.5
LABEL com.nvidia.cuda.version="7.5"

# install cuda packages
RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
    cuda-toolkit-7-5=7.5-18 \
    && rm -rf /var/lib/apt/lists/*

# configure cuda toolkit
RUN echo "/usr/local/cuda/lib" >> /etc/ld.so.conf.d/cuda.conf && \
    echo "/usr/local/cuda/lib64" >> /etc/ld.so.conf.d/cuda.conf && \
    ldconfig

# configure nvidia libraries
RUN echo "/usr/local/nvidia/lib" >> /etc/ld.so.conf.d/nvidia.conf && \
    echo "/usr/local/nvidia/lib64" >> /etc/ld.so.conf.d/nvidia.conf

# setup environment variables
ENV PATH /usr/local/cuda/bin:${PATH}
ENV LD_LIBRARY_PATH /usr/local/nvidia/lib:/usr/local/nvidia/lib64:${LD_LIBRARY_PATH}
