# Repository configuration

In order to setup the nvidia-docker repository for your distribution, follow the instructions below.

If you feel something is missing or requires additional information, please let us know by [filing a new issue](https://github.com/NVIDIA/nvidia-docker/issues/new).

## Debian-based distributions

#### Ubuntu 14.04/16.04, Debian Jessie/Stretch
```bash
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | \
  sudo apt-key add -
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | \
  sudo tee /etc/apt/sources.list.d/nvidia-docker.list
sudo apt-get update
```

## CentOS/RHEL-based distributions

#### CentOS 7, RHEL 7.4, Amazon Linux 2
```bash
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.repo | \
  sudo tee /etc/yum.repos.d/nvidia-docker.repo
```

#### Amazon Linux 1
```bash
curl -s -L https://nvidia.github.io/nvidia-docker/amzn1/nvidia-docker.repo | \
  sudo tee /etc/yum.repos.d/nvidia-docker.repo
```