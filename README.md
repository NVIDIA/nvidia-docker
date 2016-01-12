# NVIDIA Docker

This repository includes utilities to build and run NVIDIA Docker images.

![nvidia-gpu-docker](https://cloud.githubusercontent.com/assets/3028125/12213714/5b208976-b632-11e5-8406-38d379ec46aa.png)

> Example of how CUDA integrates with Docker

# Documentation

The full documentation is available on the [repository wiki](https://github.com/NVIDIA/nvidia-docker/wiki).  
A good place to start is to understand [why NVIDIA Docker](https://github.com/NVIDIA/nvidia-docker/wiki/Why%20NVIDIA%20Docker) is needed in the first place.

# Quick start

Assuming the NVIDIA drivers and Docker are properly installed (see [installation](https://github.com/NVIDIA/nvidia-docker/wiki/Installation)):

```sh
git clone https://github.com/NVIDIA/nvidia-docker

# Initial setup
sudo make install
sudo nvidia-docker volume setup

# Run nvidia-smi
nvidia-docker run nvidia/cuda nvidia-smi
```

# Issues and Contributing

**A signed copy of the [Contributor License Agreement](https://raw.githubusercontent.com/NVIDIA/nvidia-docker/master/CLA) needs to be provided to digits@nvidia.com before any change can be accepted.**

* Please let us know by [filing a new issue](https://github.com/NVIDIA/nvidia-docker/issues/new)
* You can contribute by opening a [pull request](https://help.github.com/articles/using-pull-requests/)
