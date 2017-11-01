# Docker Engine Utility for NVIDIA GPUs

[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg?style=flat-square)](https://raw.githubusercontent.com/NVIDIA/nvidia-docker/master/LICENSE)
[![Package repository](https://img.shields.io/badge/packages-repository-b956e8.svg?style=flat-square)](https://nvidia.github.io/nvidia-docker)

![nvidia-gpu-docker](https://cloud.githubusercontent.com/assets/3028125/12213714/5b208976-b632-11e5-8406-38d379ec46aa.png)

This branch contains version 2.0 of the nvidia-docker utility.

**Warning: Version 2.0 is in alpha state, it is not intended to be used in production systems.**

## Differences with 1.0
* Doesn't require wrapping the Docker CLI,
* Doesn't require starting a separate daemon,
* GPU isolation is now achieved with environment variable `NVIDIA_VISIBLE_DEVICES`,
* Can enable GPU support for any Docker image. Not just the ones based on our official CUDA images,
* Package repositories are available for Ubuntu and CentOS,
* Uses a new implementation based on [libnvidia-container](https://github.com/NVIDIA/libnvidia-container).

## Removing nvidia-docker 1.0

Version 1.0 of the nvidia-docker package must be cleanly removed before continuing.  
You must stop and remove **all** containers started with nvidia-docker 1.0.

#### Ubuntu distributions
```sh
docker volume ls -q -f driver=nvidia-docker | xargs -r -I{} -n1 docker ps -q -a -f volume={} | xargs -r docker rm -f
sudo apt-get purge nvidia-docker
```

#### CentOS distributions

```
docker volume ls -q -f driver=nvidia-docker | xargs -r -I{} -n1 docker ps -q -a -f volume={} | xargs -r docker rm -f
sudo yum remove nvidia-docker
```

## Installation

**If you have a custom `/etc/docker/daemon.json`, the `nvidia-docker2` package will override it.**  
**In this case, it is recommended to install [nvidia-container-runtime](https://github.com/nvidia/nvidia-container-runtime#installation) instead and register the new runtime manually.**

#### Ubuntu distributions

1. Install the repository for your distribution by following the instructions [here](http://nvidia.github.io/nvidia-docker/).
2. Install the `nvidia-docker2` package and restart the Docker daemon:
```
sudo apt-get install nvidia-docker2
sudo pkill -SIGHUP dockerd
```

#### CentOS distributions
1. Install the repository for your distribution by following the instructions [here](http://nvidia.github.io/nvidia-docker/).
2. Install the `nvidia-docker2` package and restart the Docker daemon:
```
sudo yum install nvidia-docker2
sudo pkill -SIGHUP dockerd
```

## Usage

#### NVIDIA runtime
nvidia-docker 2.0 registers a new container runtime to the Docker daemon.  
You must select the `nvidia` runtime when using `docker run`:
```
docker run --runtime=nvidia --rm nvidia/cuda nvidia-smi
```

#### GPU isolation
Set the environment variable `NVIDIA_VISIBLE_DEVICES` in the container:
```
docker run --runtime=nvidia -e NVIDIA_VISIBLE_DEVICES=0 --rm nvidia/cuda nvidia-smi
```

#### Non-CUDA image:
Setting `NVIDIA_VISIBLE_DEVICES` will enable GPU support for any container image:
```
docker run --runtime=nvidia -e NVIDIA_VISIBLE_DEVICES=all --rm debian:stretch nvidia-smi
```

## Advanced

#### Backward compatibility

To help transitioning code from 1.0 to 2.0, a bash script is provided in `/usr/bin/nvidia-docker` for backward compatibility.  
It will automatically inject the `--runtime=nvidia` argument and convert `NV_GPU` to `NVIDIA_VISIBLE_DEVICES`.

#### Default runtime
The default runtime used by the Docker® Engine is [runc](https://github.com/opencontainers/runc), our runtime can become the default one by configuring the docker daemon with `--default-runtime=nvidia`.
Doing so will remove the need to add the `--runtime=nvidia` argument to `docker run`.
It is also the only way to have GPU access during `docker build`.

#### Environment variables
The behavior of the runtime can be modified through environment variables (such as `NVIDIA_VISIBLE_DEVICES`).   
Those environment variables are consumed by [nvidia-container-runtime](https://github.com/nvidia/nvidia-container-runtime) and are documented [here](https://github.com/nvidia/nvidia-container-runtime#environment-variables-oci-spec).  
Our official CUDA images use default values for these variables.

## Issues and Contributing

A signed copy of the [Contributor License Agreement](https://raw.githubusercontent.com/NVIDIA/nvidia-docker/master/CLA) needs to be provided to <a href="mailto:digits@nvidia.com">digits@nvidia.com</a> before any change can be accepted.

* Please let us know by [filing a new issue](https://github.com/NVIDIA/nvidia-docker/issues/new)
* You can contribute by opening a [pull request](https://help.github.com/articles/using-pull-requests/)
