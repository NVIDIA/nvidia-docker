# NVIDIA Container Runtime for Docker

[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg?style=flat-square)](https://raw.githubusercontent.com/NVIDIA/nvidia-docker/master/LICENSE)
[![Documentation](https://img.shields.io/badge/documentation-wiki-blue.svg?style=flat-square)](https://github.com/NVIDIA/nvidia-docker/wiki)
[![Package repository](https://img.shields.io/badge/packages-repository-b956e8.svg?style=flat-square)](https://nvidia.github.io/nvidia-docker)

![nvidia-gpu-docker](https://cloud.githubusercontent.com/assets/3028125/12213714/5b208976-b632-11e5-8406-38d379ec46aa.png)

# Documentation

The full documentation and frequently asked questions are available on the [repository wiki](https://github.com/NVIDIA/nvidia-docker/wiki).  

An introduction to the NVIDIA Container Runtime is also covered in our [blog post](https://devblogs.nvidia.com/gpu-containers-runtime/).

## Quickstart

**Make sure you have installed the [NVIDIA driver](https://github.com/NVIDIA/nvidia-docker/wiki/Frequently-Asked-Questions#how-do-i-install-the-nvidia-driver) and a [supported version](https://github.com/NVIDIA/nvidia-docker/wiki/Frequently-Asked-Questions#which-docker-packages-are-supported) of [Docker](https://docs.docker.com/engine/installation/) for your distribution (see [prerequisites](https://github.com/NVIDIA/nvidia-docker/wiki/Installation-(version-2.0)#prerequisites)).**

**If you have a custom `/etc/docker/daemon.json`, the `nvidia-docker2` package might override it.**  

#### Ubuntu 14.04/16.04/18.04, Debian Jessie/Stretch

Ubuntu will install `docker.io` by default which isn't the latest version of Docker Engine. This implies that you will need to pin the version of nvidia-docker. [See more information here](https://github.com/NVIDIA/nvidia-docker/wiki/Frequently-Asked-Questions#how-do-i-install-20-if-im-not-using-the-latest-docker-version).

```sh
# If you have nvidia-docker 1.0 installed: we need to remove it and all existing GPU containers
docker volume ls -q -f driver=nvidia-docker | xargs -r -I{} -n1 docker ps -q -a -f volume={} | xargs -r docker rm -f
sudo apt-get purge -y nvidia-docker

# Add the package repositories
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | \
  sudo apt-key add -
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | \
  sudo tee /etc/apt/sources.list.d/nvidia-docker.list
sudo apt-get update

# Install nvidia-docker2 and reload the Docker daemon configuration
sudo apt-get install -y nvidia-docker2
sudo pkill -SIGHUP dockerd

# Test nvidia-smi with the latest official CUDA image
docker run --runtime=nvidia --rm nvidia/cuda:9.0-base nvidia-smi
```

#### CentOS 7 (docker-ce), RHEL 7.4/7.5 (docker-ce), Amazon Linux 1/2

If you are **not** using the official `docker-ce` package on CentOS/RHEL, use the next section.

```sh
# If you have nvidia-docker 1.0 installed: we need to remove it and all existing GPU containers
docker volume ls -q -f driver=nvidia-docker | xargs -r -I{} -n1 docker ps -q -a -f volume={} | xargs -r docker rm -f
sudo yum remove nvidia-docker

# Add the package repositories
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.repo | \
  sudo tee /etc/yum.repos.d/nvidia-docker.repo

# Install nvidia-docker2 and reload the Docker daemon configuration
sudo yum install -y nvidia-docker2
sudo pkill -SIGHUP dockerd

# Test nvidia-smi with the latest official CUDA image
docker run --runtime=nvidia --rm nvidia/cuda:9.0-base nvidia-smi
```
If `yum` reports a conflict on `/etc/docker/daemon.json` with the
`docker` package, you need to use the next section instead.

For docker-ce on `ppc64le`, look at the [FAQ](https://github.com/nvidia/nvidia-docker/wiki/Frequently-Asked-Questions#do-you-support-powerpc64-ppc64le).

#### CentOS 7 (docker), RHEL 7.4/7.5 (docker)
```sh
# If you have nvidia-docker 1.0 installed: we need to remove it and all existing GPU containers
docker volume ls -q -f driver=nvidia-docker | xargs -r -I{} -n1 docker ps -q -a -f volume={} | xargs -r docker rm -f
sudo yum remove nvidia-docker

# Add the package repositories
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-container-runtime/$distribution/nvidia-container-runtime.repo | \
  sudo tee /etc/yum.repos.d/nvidia-container-runtime.repo

# Install the nvidia runtime hook
sudo yum install -y nvidia-container-runtime-hook

# Test nvidia-smi with the latest official CUDA image
# You can't use `--runtime=nvidia` with this setup.
docker run --rm nvidia/cuda:9.0-base nvidia-smi
```

#### Other distributions and architectures

Look at the [Installation section](https://github.com/nvidia/nvidia-docker/wiki/Installation-(version-2.0)) of the wiki.

## Changelog

New nvidia-docker packages have been released for docker 18.09.2 and 18.06.2 addressing [CVE-2019-5736](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2019-5736)

Update nvidia-docker to adress runc critical vulnerability that allows specially-crafted containers to gain administrative privileges on the host.
```sh
# On Ubuntu/Debian
sudo apt upgrade

# On Centos/RHEL/Amazon Linux
sudo yum upgrade
```

## Issues and Contributing

A signed copy of the [Contributor License Agreement](https://raw.githubusercontent.com/NVIDIA/nvidia-docker/master/CLA) needs to be provided to <a href="mailto:digits@nvidia.com">digits@nvidia.com</a> before any change can be accepted.

* Please let us know by [filing a new issue](https://github.com/NVIDIA/nvidia-docker/issues/new)
* You can contribute by opening a [pull request](https://help.github.com/articles/using-pull-requests/)
