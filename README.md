# NVIDIA Container Toolkit

[![GitHub license](https://img.shields.io/github/license/NVIDIA/nvidia-docker?style=flat-square)](https://raw.githubusercontent.com/NVIDIA/nvidia-docker/master/LICENSE)
[![Documentation](https://img.shields.io/badge/documentation-wiki-blue.svg?style=flat-square)](https://github.com/NVIDIA/nvidia-docker/wiki)
[![Package repository](https://img.shields.io/badge/packages-repository-b956e8.svg?style=flat-square)](https://nvidia.github.io/nvidia-docker)

![nvidia-gpu-docker](https://cloud.githubusercontent.com/assets/3028125/12213714/5b208976-b632-11e5-8406-38d379ec46aa.png)

## Introduction
The NVIDIA Container Toolkit allows users to build and run GPU accelerated Docker containers. The toolkit includes a container runtime [library](https://github.com/NVIDIA/libnvidia-container) and utilities to automatically configure containers to leverage NVIDIA GPUs. Full documentation and frequently asked questions are available on the [repository wiki](https://github.com/NVIDIA/nvidia-docker/wiki).

## Quickstart

**Make sure you have installed the [NVIDIA driver](https://github.com/NVIDIA/nvidia-docker/wiki/Frequently-Asked-Questions#how-do-i-install-the-nvidia-driver) and Docker 19.03 for your Linux distribution**

Note that with the release of Docker 19.03, usage of nvidia-docker2 packages are deprecated since NVIDIA GPUs are now natively supported as devices in the Docker runtime.
If you are an existing user of the nvidia-docker2 packages, review the instructions in the [“Upgrading with nvidia-docker2” section](https://github.com/NVIDIA/nvidia-docker/tree/master#upgrading-with-nvidia-docker2-deprecated).

For first-time users of Docker 19.03 and GPUs, continue with the instructions for getting started below.

### Ubuntu 16.04/18.04, Debian Jessie/Stretch
```sh
# Add the package repositories
$ distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
$ curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
$ curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list

$ sudo apt-get update && sudo apt-get install -y nvidia-container-toolkit
$ sudo systemctl restart docker
```

#### CentOS 7 (docker-ce), RHEL 7.4/7.5 (docker-ce), Amazon Linux 1/2
```
$ distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
$ curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.repo | sudo tee /etc/yum.repos.d/nvidia-docker.repo

$ sudo yum install -y nvidia-container-toolkit
$ sudo systemctl restart docker
```

## Usage
```
#### Test nvidia-smi with the latest official CUDA image
$ docker run --gpus all nvidia/cuda:9.0-base nvidia-smi

# Start a GPU enabled container on two GPUs
$ docker run --gpus 2 nvidia/cuda:9.0-base nvidia-smi

# Starting a GPU enabled container on specific GPUs
$ docker run --gpus '"device=1,2"' nvidia/cuda:9.0-base nvidia-smi
$ docker run --gpus '"device=UUID-ABCDEF,1"' nvidia/cuda:9.0-base nvidia-smi

# Specifying a capability (graphics, compute, ...) for my container
# Note this is rarely if ever used this way
$ docker run --gpus all,capabilities=utility nvidia/cuda:9.0-base nvidia-smi
```

## Upgrading with nvidia-docker2 (Deprecated)

If you are running an old version of docker (< 19.03) check the instructions on installing the [`nvidia-docker2`](https://github.com/NVIDIA/nvidia-docker/wiki/Installation-(version-2.0)) package which supports Docker >= 1.12.
If you already have the old package installed (nvidia-docker2), updating to the latest Docker version (>= 19.03) will still work and will  give you access to the new CLI options for supporting GPUs:

```
# On debian based distributions: Ubuntu / Debian
$ sudo apt-get update
$ sudo apt-get --only-upgrade install docker-ce nvidia-docker2
$ sudo systemctl restart docker

# On RPM based distributions: Centos / RHEL / Amazon Linux
$ sudo yum upgrade -y nvidia-docker2
$ sudo systemctl restart docker

# All of the following options will continue working
$ docker run --gpus all nvidia/cuda:9.0-base nvidia-smi
$ docker run --runtime nvidia nvidia/cuda:9.0-base nvidia-smi
$ nvidia-docker run nvidia/cuda:9.0-base nvidia-smi
```

Note that in the future, nvidia-docker2 packages will no longer be supported.

## Issues and Contributing

[Checkout the Contributing document!](CONTRIBUTING.md)

* Please let us know by [filing a new issue](https://github.com/NVIDIA/nvidia-docker/issues/new)
* You can contribute by opening a [pull request](https://help.github.com/articles/using-pull-requests/)
