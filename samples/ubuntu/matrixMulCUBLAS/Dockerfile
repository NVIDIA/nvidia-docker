FROM cuda

RUN apt-get update && apt-get install -y --no-install-recommends \
        cuda-samples-$CUDA_PKG_VERSION && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /usr/local/cuda/samples/0_Simple/matrixMulCUBLAS
RUN make

CMD ./matrixMulCUBLAS -sizemult=10
