# NVIDIA Container Toolkit

[![GitHub license](https://img.shields.io/github/license/NVIDIA/nvidia-docker?style=flat-square)](https://raw.githubusercontent.com/NVIDIA/nvidia-docker/master/LICENSE)
[![Documentation](https://img.shields.io/badge/documentation-wiki-blue.svg?style=flat-square)](https://github.com/NVIDIA/nvidia-docker/wiki)
[![Package repository](https://img.shields.io/badge/packages-repository-b956e8.svg?style=flat-square)](https://nvidia.github.io/nvidia-docker)

![nvidia-gpu-docker](https://cloud.githubusercontent.com/assets/3028125/12213714/5b208976-b632-11e5-8406-38d379ec46aa.png)

## Introduction
The NVIDIA Container Toolkit allows users to build and run GPU accelerated Docker containers. The toolkit includes a container runtime [library](https://github.com/NVIDIA/libnvidia-container) and utilities to automatically configure containers to leverage NVIDIA GPUs. Full documentation and frequently asked questions are available on the [repository wiki](https://github.com/NVIDIA/nvidia-docker/wiki).

## Quickstart

**Make sure you have installed the [NVIDIA driver](https://github.com/NVIDIA/nvidia-docker/wiki/Frequently-Asked-Questions#how-do-i-install-the-nvidia-driver) and Docker 19.03 for your Linux distribution**
**Note that you do not need to install the CUDA toolkit on the host, but the driver needs to be installed**

Note that with the release of Docker 19.03, usage of nvidia-docker2 packages are deprecated since NVIDIA GPUs are now natively supported as devices in the Docker runtime.

**Please note that this native GPU support has not landed in docker-compose yet. Refer to [this issue](https://github.com/docker/compose/issues/6691) for discussion.**

If you are an existing user of the nvidia-docker2 packages, review the instructions in the [“Upgrading with nvidia-docker2” section](https://github.com/NVIDIA/nvidia-docker/tree/master#upgrading-with-nvidia-docker2-deprecated).

For first-time users of Docker 19.03 and GPUs, continue with the instructions for getting started below.

### Ubuntu 16.04/18.04, Debian Jessie/Stretch/Buster
```sh
# Add the package repositories
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list

sudo apt-get update && sudo apt-get install -y nvidia-container-toolkit
sudo systemctl restart docker
```

#### CentOS 7 (docker-ce), RHEL 7.4/7.5 (docker-ce), Amazon Linux 1/2
```
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.repo | sudo tee /etc/yum.repos.d/nvidia-docker.repo

sudo yum install -y nvidia-container-toolkit
sudo systemctl restart docker
```

#### openSUSE Leap 15.1 (docker-ce)

Since openSUSE Leap 15.1 still has Docker 18.06, you have two options:

**Option 1**: use the `Virtualization:containers` repository to fetch a more recent version of Docker

```console
# Upgrade Docker to 19.03+ first:
zypper ar https://download.opensuse.org/repositories/Virtualization:/containers/openSUSE_Leap_15.1/Virtualization:containers.repo
zypper install --allow-vendor-change 'docker >= 19.03'  # accept the new signature

# Add the package repositories
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
zypper ar https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.repo

sudo zypper install -y nvidia-container-toolkit
sudo systemctl restart docker
```

**Option 2**: stay with the deprecated `nvidia-docker2` package for now (see also below)

```console
# Add the package repositories
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
zypper ar https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.repo

sudo zypper install -y nvidia-docker2  # accept the overwrite of /etc/docker/daemon.json
sudo systemctl restart docker
```

## Usage
```
#### Test nvidia-smi with the latest official CUDA image
docker run --gpus all nvidia/cuda:10.0-base nvidia-smi

# Start a GPU enabled container on two GPUs
docker run --gpus 2 nvidia/cuda:10.0-base nvidia-smi

# Starting a GPU enabled container on specific GPUs
docker run --gpus '"device=1,2"' nvidia/cuda:10.0-base nvidia-smi
docker run --gpus '"device=UUID-ABCDEF,1"' nvidia/cuda:10.0-base nvidia-smi

# Specifying a capability (graphics, compute, ...) for my container
# Note this is rarely if ever used this way
docker run --gpus all,capabilities=utility nvidia/cuda:10.0-base nvidia-smi
```

## RHEL Docker or Podman

_Note that RHEL's fork of Docker is no longer supported on RHEL8._
_Note that for powerpc you will have to install the nvidia-container-runtime-hook_

RHEL's fork of docker doesn't support the --gpus option, in this case you should still install
the nvidia-container-toolkit package but you will have to use the old interface. e.g:
```bash
# Add the package repositories
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.repo | sudo tee /etc/yum.repos.d/nvidia-docker.repo

# On x86
sudo yum install -y nvidia-container-toolkit
# On PPC
sudo yum install -y nvidia-container-hook
sudo systemctl restart docker

# On RHEL 7/8
docker run -e NVIDIA_VISIBLE_DEVICES=all nvidia/cuda:10.0-base nvidia-smi

# With Podman
podman run -e NVIDIA_VISIBLE_DEVICES=all nvidia/cuda:10.0-base nvidia-smi
```

More information on the environment variables are available [on this page](https://github.com/NVIDIA/nvidia-container-runtime#environment-variables-oci-spec).

## Upgrading with nvidia-docker2 (Deprecated)

If you are running an old version of docker (< 19.03) check the instructions on installing the [`nvidia-docker2`](https://github.com/NVIDIA/nvidia-docker/wiki/Installation-(version-2.0)) package which supports Docker >= 1.12.
If you already have the old package installed (nvidia-docker2), updating to the latest Docker version (>= 19.03) will still work and will  give you access to the new CLI options for supporting GPUs:

```
# On debian based distributions: Ubuntu / Debian
sudo apt-get update
sudo apt-get --only-upgrade install docker-ce nvidia-docker2
sudo systemctl restart docker

# On RPM based distributions: Centos / RHEL / Amazon Linux
sudo yum upgrade -y nvidia-docker2
sudo systemctl restart docker

# All of the following options will continue working
docker run --gpus all nvidia/cuda:10.0-base nvidia-smi
docker run --runtime nvidia nvidia/cuda:10.0-base nvidia-smi
nvidia-docker run nvidia/cuda:10.0-base nvidia-smi
```

Note that in the future, nvidia-docker2 packages will no longer be supported.

## Changelog

* Friday September 20th:
  We changed the gpgkey, the new fingerprint is: `BC02 13EE FC50 D046 F1CE  0208 6128 B5C2 36CD EE96`
  We will add a webpage on docs.nvidia.com listing the keys and their fingerprints.
  In the future we will publish a keyring package. This will allow automatic updates to the repository keys.
  Future updates to the keys will be communicated in advance. We apologize for any inconvenience caused by the unexpected change to the keys


## Issues and Contributing

[Checkout the Contributing document!](CONTRIBUTING.md)

* Please let us know by [filing a new issue](https://github.com/NVIDIA/nvidia-docker/issues/new)
* You can contribute by opening a [pull request](https://help.github.com/articles/using-pull-requests/)
