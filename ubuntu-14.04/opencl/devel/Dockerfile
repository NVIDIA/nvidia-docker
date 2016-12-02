FROM opencl:runtime
LABEL maintainer "NVIDIA CORPORATION <cudatools@nvidia.com>"

RUN apt-get update && apt-get install -y --no-install-recommends \
        clinfo \
        ocl-icd-opencl-dev && \
    rm -rf /var/lib/apt/lists/*
