FROM cuda:latest

WORKDIR /usr/local/cuda/samples/1_Utilities/bandwidthTest
RUN make

CMD ./bandwidthTest --mode=shmoo
