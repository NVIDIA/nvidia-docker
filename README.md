# NVIDIA Container Toolkit

[![GitHub license](https://img.shields.io/github/license/NVIDIA/nvidia-docker?style=flat-square)](https://raw.githubusercontent.com/NVIDIA/nvidia-docker/main/LICENSE)
[![Documentation](https://img.shields.io/badge/documentation-wiki-blue.svg?style=flat-square)](https://github.com/NVIDIA/nvidia-docker/wiki)
[![Package repository](https://img.shields.io/badge/packages-repository-b956e8.svg?style=flat-square)](https://nvidia.github.io/nvidia-docker)

![nvidia-gpu-docker](https://cloud.githubusercontent.com/assets/3028125/12213714/5b208976-b632-11e5-8406-38d379ec46aa.png)

## Introduction
The NVIDIA Container Toolkit allows users to build and run GPU accelerated Docker containers. The toolkit includes a container runtime [library](https://github.com/NVIDIA/libnvidia-container) and utilities to automatically configure containers to leverage NVIDIA GPUs.

Product documentation including an architecture overview, platform support, installation and usage guides can be found in the [documentation repository](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/overview.html).

Frequently asked questions are available on the [wiki](https://github.com/NVIDIA/nvidia-docker/wiki).

## Getting Started

On the host system:
1. Install the [NVIDIA driver](https://docs.nvidia.com/datacenter/tesla/tesla-installation-notes/index.html) and Docker engine for your Linux distribution.
2. Install the NVIDIA Container Toolkit, following the [installation guide](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html).

Note:
- The NVIDIA Container Toolkit's Linux package is called `nvidia-docker2`. Older versions are distributed under `nvidia-docker`. See the installation guide for more details.
- You do **not** need to install the CUDA Toolkit on the host system. Only in the image/container.


## Usage

The [user guide](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/user-guide.html) provides information on the configuration and command line options available when running GPU containers with Docker.

## Issues and Contributing

[Checkout the Contributing document!](CONTRIBUTING.md)

* Please let us know by [filing a new issue](https://github.com/NVIDIA/nvidia-docker/issues/new)
* You can contribute by opening a [merge request](https://gitlab.com/nvidia/container-toolkit/nvidia-docker/-/merge_requests)
