FROM caffe:0.14
MAINTAINER NVIDIA CORPORATION <digits@nvidia.com>

# workaround: gcc and libhdf5-dev are dependencies that are currently missing from the torch package
RUN apt-get update && apt-get install -y --no-install-recommends --force-yes \
        gcc \
        git \
        graphviz \
        gunicorn \
        libhdf5-dev \
        nginx \
        python-caffe-nv \
        python-flask \
        python-flaskext.socketio \
        python-flaskext.wtf \
        python-gevent \
        python-lmdb \
        python-pil \
        python-pip \
        python-pydot \
        python-requests \
        python-six \
        python-skimage \
        python-wtforms \
        torch7-nv=0.9.98-1+cuda7.5 && \
    rm -rf /var/lib/apt/lists/*

ENV DIGITS_VERSION 4.0
LABEL com.nvidia.digits.version="4.0"

WORKDIR /opt/digits

ENV CLONE_TAG=v4.0.0
RUN git clone -b ${CLONE_TAG} --depth 1 https://github.com/NVIDIA/DIGITS.git . && \
    for req in $(cat requirements.txt); do pip install $req; done

VOLUME /data
VOLUME /jobs

COPY digits.cfg digits/digits.cfg

EXPOSE 34448
ENTRYPOINT ["/opt/digits/digits-server"]
