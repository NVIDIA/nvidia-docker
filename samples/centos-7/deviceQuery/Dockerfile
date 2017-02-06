FROM nvidia/cuda:8.0-devel-centos7

RUN yum install -y \
        cuda-samples-$CUDA_PKG_VERSION && \
    rm -rf /var/cache/yum/*

WORKDIR /usr/local/cuda/samples/1_Utilities/deviceQuery
RUN make

CMD ./deviceQuery
