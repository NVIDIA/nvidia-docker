# Docker Engine Utility for NVIDIA GPUs

**We have now transitioned to [nvidia-docker 2.0](https://github.com/NVIDIA/nvidia-docker/tree/master).**

![nvidia-gpu-docker](https://cloud.githubusercontent.com/assets/3028125/12213714/5b208976-b632-11e5-8406-38d379ec46aa.png)

# Documentation

The full documentation is available on the [repository wiki](https://github.com/NVIDIA/nvidia-docker/wiki), section "Version 1.0".  
A good place to start is to understand [why nvidia-docker](https://github.com/NVIDIA/nvidia-docker/wiki/Motivation) is needed in the first place.


# Quick start

Assuming the NVIDIA drivers and Docker® Engine are properly installed (see [installation](https://github.com/NVIDIA/nvidia-docker/wiki/Installation-(version-1.0)))

#### _Ubuntu distributions_
```sh
# Install nvidia-docker and nvidia-docker-plugin
wget -P /tmp https://github.com/NVIDIA/nvidia-docker/releases/download/v1.0.1/nvidia-docker_1.0.1-1_amd64.deb
sudo dpkg -i /tmp/nvidia-docker*.deb && rm /tmp/nvidia-docker*.deb

# Test nvidia-smi
nvidia-docker run --rm nvidia/cuda nvidia-smi
```

#### _CentOS distributions_
```sh
# Install nvidia-docker and nvidia-docker-plugin
wget -P /tmp https://github.com/NVIDIA/nvidia-docker/releases/download/v1.0.1/nvidia-docker-1.0.1-1.x86_64.rpm
sudo rpm -i /tmp/nvidia-docker*.rpm && rm /tmp/nvidia-docker*.rpm
sudo systemctl start nvidia-docker

# Test nvidia-smi
nvidia-docker run --rm nvidia/cuda nvidia-smi
```

#### _Other distributions_
```sh
# Install nvidia-docker and nvidia-docker-plugin
wget -P /tmp https://github.com/NVIDIA/nvidia-docker/releases/download/v1.0.1/nvidia-docker_1.0.1_amd64.tar.xz
sudo tar --strip-components=1 -C /usr/bin -xvf /tmp/nvidia-docker*.tar.xz && rm /tmp/nvidia-docker*.tar.xz

# Run nvidia-docker-plugin
sudo -b nohup nvidia-docker-plugin > /tmp/nvidia-docker.log

# Test nvidia-smi
nvidia-docker run --rm nvidia/cuda nvidia-smi
```

#### _ppc64le (POWER) Archictecture_
There is limited build support for ppc64le. Running `make deb` will build the nvidia-docker deb for ppc64le (if run on a ppc64le system). If the deb install fails because you have the 'docker.io' (>= v1.9) package installed, but not the 'docker-engine' package, you can force-install. There is currently no docker-provided docker-engine repository for ppc64le.

Not all the build targets for ppc64le have been implemented. If you would like for a Dockerfile to be created to enable a ppc64le target, please open an issue.

# Issues and Contributing

**A signed copy of the [Contributor License Agreement](https://raw.githubusercontent.com/NVIDIA/nvidia-docker/master/CLA) needs to be provided to digits@nvidia.com before any change can be accepted.**

* Please let us know by [filing a new issue](https://github.com/NVIDIA/nvidia-docker/issues/new)
* You can contribute by opening a [pull request](https://help.github.com/articles/using-pull-requests/)
