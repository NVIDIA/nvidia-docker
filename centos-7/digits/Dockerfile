FROM caffe

ENV DIGITS_VERSION 3.0
LABEL com.nvidia.digits.version="3.0"

# Install build essential
ENV DIGITS_PKG_VERSION 3.0.0~rc.3-1
RUN yum update -y && yum install -y \
	graphviz 

# Get source code
ENV DIGITS_GIT https://github.com/NVIDIA/DIGITS.git
ENV DIGITS_TAG v3.0.0-rc.1
ENV DIGITS_WS /tmp/digits
RUN git clone --single-branch -b $DIGITS_TAG $DIGITS_GIT $DIGITS_WS 

# Setup python packages
RUN cd $DIGITS_WS && \
	pip install -r requirements.txt

# Fix issue of dependencies
RUN pip install Flask-SocketIO==0.6.0
# Fix issue of dependencies
RUN pip install Flask-SocketIO==0.6.0

RUN mkdir -p /var/log/digits
RUN mkdir -p /usr/share/digits/digits/jobs

RUN echo -e "[DIGITS]\ncaffe_root = /tmp/caffe\n" > $DIGITS_WS/digits/digits.cfg
WORKDIR $DIGITS_WS

# EXPOSE 34448
# ENTRYPOINT ["./digits-devserver"]

