FROM cuda:7.0-cudnn4-devel

ENV CAFFE_VERSION 0.14
LABEL com.nvidia.caffe.version="0.14"

# Install fpm and rpmrebuild
RUN yum install -y epel-release
RUN rpm --import http://li.nux.ro/download/nux/RPM-GPG-KEY-nux.ro
RUN rpm -Uvh http://li.nux.ro/download/nux/dextop/el7/x86_64/nux-dextop-release-0-5.el7.nux.noarch.rpm

# Install build essential
RUN yum update -y && yum install -y \
	protobuf-devel leveldb-devel snappy-devel \
	opencv-devel boost-devel hdf5-devel \
	gflags-devel glog-devel lmdb-devel && \
	rm -rf /var/cache/yum/*

RUN yum update -y && yum install -y \
	wget git cmake make automake openblas-devel && \
	rm -rf /var/cache/yum/*

# Get source code
ENV CAFFE_GIT https://github.com/NVIDIA/caffe.git
ENV CAFFE_TAG v0.14.2
ENV CAFFE_WS /tmp/caffe
RUN git clone --single-branch -b $CAFFE_TAG $CAFFE_GIT $CAFFE_WS 

# Get Python dependencies
RUN yum update -y && yum install -y \
	python-pip python-devel gcc-gfortran\
	freetype-devel libpng-devel libjpeg-turbo-devel && \
	rm -rf /var/cache/yum/*

RUN pip install --upgrade pip

# scipy

RUN cd $CAFFE_WS/python && for req in $(cat requirements.txt); do pip install $req; done

RUN mkdir -p $CAFFE_WS/build
RUN cd $CAFFE_WS/build && \
	cmake ../ -DCUDNN_ROOT=/usr/lib/x86_64-linux-gnu -DBLAS=open
RUN cd $CAFFE_WS/build && \
	make all && make install

