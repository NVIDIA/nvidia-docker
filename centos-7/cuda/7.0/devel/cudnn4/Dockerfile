FROM cuda:7.0-devel

ENV CUDNN_VERSION 4
LABEL com.nvidia.cudnn.version="4"

ENV CUDNN_DOWNLOAD_SUM 4e64ef7716f20c87854b4421863328e17cce633330c319b5e13809b61a36f97d
ENV CUDNN_DOWNLOAD_FILENAME  cudnn-7.0-linux-x64-v4.0-prod.tgz

RUN curl -fsSL http://developer.download.nvidia.com/compute/redist/cudnn/v4/$CUDNN_DOWNLOAD_FILENAME -O && \
    echo "$CUDNN_DOWNLOAD_SUM $CUDNN_DOWNLOAD_FILENAME" | sha256sum -c --strict - && \
    tar -xzf $CUDNN_DOWNLOAD_FILENAME -C /usr/local && \
    rm $CUDNN_DOWNLOAD_FILENAME && \
    ldconfig

