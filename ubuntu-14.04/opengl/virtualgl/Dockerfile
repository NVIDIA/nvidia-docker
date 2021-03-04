FROM opengl:runtime
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
        ca-certificates \
        curl \
        libxv1 && \
    rm -rf /var/lib/apt/lists/*

RUN curl -fsSL -o virtualgl.deb https://sourceforge.net/projects/virtualgl/files/2.5/virtualgl_2.5_amd64.deb/download && \
    dpkg -i virtualgl.deb && rm virtualgl.deb

ENV PATH /opt/VirtualGL/bin/:${PATH}

ENTRYPOINT ["/opt/VirtualGL/bin/vglrun"]
