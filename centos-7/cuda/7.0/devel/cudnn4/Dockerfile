FROM cuda:7.0-devel

# Install fpm and rpmrebuild
RUN yum install -y epel-release

RUN rpm --import http://li.nux.ro/download/nux/RPM-GPG-KEY-nux.ro

RUN rpm -Uvh http://li.nux.ro/download/nux/dextop/el7/x86_64/nux-dextop-release-0-5.el7.nux.noarch.rpm

RUN yum -y update && yum install -y \
	ruby-devel gcc automake make rpmrebuild && \
	rm -rf /var/cache/yum/*

RUN gem install fpm

# Download ML toolkit
ENV NV_ROOT /tmp/nvtemp
ENV NV_DL $NV_ROOT/download
ENV NV_RPM $NV_ROOT/rpm
ENV LIBCUDNN_DEB libcudnn4_4.0.7_amd64.deb
ENV LIBCUDNN_ARCH x86_64
ENV LIBCUDNN_RPM libcudnn4-4.0.7-1.x86_64.rpm

RUN mkdir -p $NV_DL $NV_RPM
RUN curl http://developer.download.nvidia.com/compute/machine-learning/repos/ubuntu1404/x86_64/$LIBCUDNN_DEB -o $NV_DL/$LIBCUDNN_DEB

# Transform to rpm
RUN fpm -s deb -t rpm -p $NV_ROOT $NV_DL/$LIBCUDNN_DEB
RUN rpmrebuild \
	--change-spec-whole='sed -e "s|libc6.*|glibc|g" -e "s|libgcc1.*|libgcc|g" -e "s|libstdc++6.*|libstdc++|g"' \
	-d $NV_RPM -p $NV_ROOT/$LIBCUDNN_RPM

ENV CUDNN_VERSION 4
LABEL com.nvidia.cudnn.version="4"

RUN rpm -ivh $NV_RPM/$LIBCUDNN_ARCH/$LIBCUDNN_RPM

# DEV part
ENV LIBCUDNN_DEV_DEB libcudnn4-dev_4.0.7_amd64.deb
ENV LIBCUDNN_DEV_RPM libcudnn4-dev-4.0.7-1.x86_64.rpm

RUN curl http://developer.download.nvidia.com/compute/machine-learning/repos/ubuntu1404/x86_64/$LIBCUDNN_DEV_DEB -o $NV_DL/$LIBCUDNN_DEV_DEB

RUN fpm -s deb -t rpm -p $NV_ROOT $NV_DL/$LIBCUDNN_DEV_DEB
RUN rpmrebuild \
        --change-spec-whole='sed -e "s|libc6.*|glibc|g" -e "s|libgcc1.*|libgcc|g" -e "s|libstdc++6.*|libstdc++|g" -e "s|alternatives *--force|alternatives |g"' \
        -d $NV_RPM -p $NV_ROOT/$LIBCUDNN_DEV_RPM

RUN rpm -ivh $NV_RPM/$LIBCUDNN_ARCH/$LIBCUDNN_DEV_RPM

