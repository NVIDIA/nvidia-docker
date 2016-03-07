FROM ubuntu:14.04
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

LABEL com.nvidia.volumes.needed="nvidia_driver"

RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
        libxau6 \
        libxdmcp6 \
        libxcb1 \
        libxext6 \
        libx11-6 \
        xauth && \
    rm -rf /var/lib/apt/lists/*

RUN echo "/usr/local/nvidia/lib" >> /etc/ld.so.conf.d/nvidia.conf && \
    echo "/usr/local/nvidia/lib64" >> /etc/ld.so.conf.d/nvidia.conf

RUN mkdir -p /opt/nvidia/lib && \
    ln -s /usr/local/nvidia/lib /opt/nvidia/lib/i386-linux-gnu && \
    ln -s /usr/local/nvidia/lib64 /opt/nvidia/lib/x86_64-linux-gnu

ENV PATH /usr/local/nvidia/bin:${PATH}
ENV LD_LIBRARY_PATH /usr/local/nvidia/lib:/usr/local/nvidia/lib64:${LD_LIBRARY_PATH}
ENV LD_PRELOAD /opt/nvidia/\$LIB/libGL.so.1:${LD_PRELOAD}
