FROM ubuntu:14.04
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

ENV NVIDIA_GPGKEY_SUM bd841d59a27a406e513db7d405550894188a4c1cd96bf8aa4f82f1b39e0b5c1c
ENV NVIDIA_GPGKEY_FPR 889bee522da690103c4b085ed88c3d385c37d3be

RUN apt-key adv --fetch-keys http://developer.download.nvidia.com/compute/cuda/repos/GPGKEY && \
    apt-key adv --export --no-emit-version -a $NVIDIA_GPGKEY_FPR | tail -n +2 > cudasign.pub && \
    echo "$NVIDIA_GPGKEY_SUM cudasign.pub" | sha256sum -c --strict - && rm cudasign.pub && \
    echo "deb http://developer.download.nvidia.com/compute/cuda/repos/ubuntu1404/x86_64 /" > /etc/apt/sources.list.d/cuda.list

ENV CUDA_VERSION 7.0
LABEL com.nvidia.cuda.version="7.0"

ENV CUDA_PKG_VERSION 7-0=7.0-28
RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
        cuda-nvrtc-$CUDA_PKG_VERSION \
        cuda-cusolver-$CUDA_PKG_VERSION \
        cuda-cublas-$CUDA_PKG_VERSION \
        cuda-cufft-$CUDA_PKG_VERSION \
        cuda-curand-$CUDA_PKG_VERSION \
        cuda-cusparse-$CUDA_PKG_VERSION \
        cuda-npp-$CUDA_PKG_VERSION \
        cuda-cudart-$CUDA_PKG_VERSION && \
    ln -s cuda-$CUDA_VERSION /usr/local/cuda && \
    rm -rf /var/lib/apt/lists/*

RUN echo "/usr/local/cuda/lib" >> /etc/ld.so.conf.d/cuda.conf && \
    echo "/usr/local/cuda/lib64" >> /etc/ld.so.conf.d/cuda.conf && \
    ldconfig

RUN echo "/usr/local/nvidia/lib" >> /etc/ld.so.conf.d/nvidia.conf && \
    echo "/usr/local/nvidia/lib64" >> /etc/ld.so.conf.d/nvidia.conf

ENV PATH /usr/local/cuda/bin:${PATH}
ENV LD_LIBRARY_PATH /usr/local/nvidia/lib:/usr/local/nvidia/lib64:${LD_LIBRARY_PATH}
