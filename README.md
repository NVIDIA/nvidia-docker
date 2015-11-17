# NVIDIA Docker
This repository includes utilities to build and run NVIDIA Docker images.  
Please be aware that this project is currently **experimental**.  

![docker](https://cloud.githubusercontent.com/assets/3028125/11199468/c0e09f50-8c82-11e5-846d-1f5e6a410598.png)  
*Example of how CUDA integrates with Docker*  

### Benefits of GPU containerization

* Reproducible builds
* Ease of deployment
* Isolation of individual devices
* Run across heterogeneous driver/toolkit environments *(e.g. CUDA 7.0 and 7.5 on a single host machine)*
* Requires only the NVIDIA driver

### Building images

Images can be built on any machine running Docker, it doesn't require a NVIDIA GPU nor any driver installation.

Building images is done through the Docker CLI or using the ```nvidia-docker``` wrapper similarly:
```sh
# With latest versions
$ docker build -t cuda ubuntu/cuda/latest
```

```sh
# With specific versions
$ docker build -t cuda:7.5 ubuntu-14.04/cuda/7.5
```

Alternatively, one can build an image directly from this repository:
```sh
# With latest versions
$ docker build -t cuda github.com/NVIDIA/nvidia-docker#:ubuntu/cuda/latest
```
```sh
# With specific versions
$ docker build -t cuda:7.5 github.com/NVIDIA/nvidia-docker#:ubuntu-14.04/cuda/7.5
```

### NVIDIA Docker wrapper

The ```nvidia-docker``` script is a drop-in replacement for ```docker``` CLI. In addition, it takes care of setting up the NVIDIA host driver environment inside Docker containers for proper execution.

GPUs are exported through a list of comma-separated IDs using the environment variable ```GPU```.
The numbering is the same as reported by ```nvidia-smi``` or when running CUDA code with ```CUDA_DEVICE_ORDER=PCI_BUS_ID```, it is however **different** from the default CUDA ordering.

```sh
$ GPU=0,1 ./nvidia-docker <docker-options> <docker-command> <docker-args>
```

### CUDA requirements

Running a CUDA container requires a machine with at least one CUDA-capable GPU and a driver compatible with the CUDA toolkit version you are using.  
The machine running the CUDA container **only requires the NVIDIA driver**, the CUDA toolkit doesn't have to be installed.

NVIDIA drivers are **backward-compatible** with CUDA toolkits versions

CUDA toolkit version   | Minimum driver version  | Minimum GPU architecture
:---------------------:|:-----------------------:|:-------------------------:
  6.5                  | >= 340.29               | >= 2.0 (Fermi)
  7.0                  | >= 346.46               | >= 2.0 (Fermi)
  7.5                  | >= 352.39               | >= 2.0 (Fermi)


### Samples

Once you have built the required images, a few examples are provided in the folder ```samples```.  
The following assumes that you have an image in your repository named ```cuda``` (see ```samples/deviceQuery/Dockerfile```):
```sh
# Run deviceQuery with one selected GPU
$ docker build -t device_query samples/deviceQuery
$ GPU=0 ./nvidia-docker run device_query

[ NVIDIA ] =INFO= Driver version: 352.39
[ NVIDIA ] =INFO= CUDA image version: 7.5

./deviceQuery Starting...

 CUDA Device Query (Runtime API) version (CUDART static linking)

Detected 1 CUDA Capable device(s)

Device 0: "GeForce GTX 980"
  [...]

deviceQuery, CUDA Driver = CUDART, CUDA Driver Version = 7.5, CUDA Runtime Version = 7.5, NumDevs = 1, Device0 = GeForce GTX 980
Result = PASS
```

# Issues and Contributing
* Please let us know by [filing a new issue](https://github.com/NVIDIA/nvidia-docker/issues/new)
* You can contribute by opening a [pull request](https://help.github.com/articles/using-pull-requests/)  
You will need to send a signed copy of the [Contributor License Agreement](CLA) to digits@nvidia.com before your change can be accepted.
