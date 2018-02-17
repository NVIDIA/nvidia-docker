# Repository configuration

In order to setup the nvidia-docker repository for your distribution, follow the instructions below.

If you feel something is missing or requires additional information, please let us know by [filing a new issue](https://github.com/NVIDIA/nvidia-docker/issues/new).

## Ubuntu distributions

#### Xenial x86_64

```bash
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | \
  sudo apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/ubuntu16.04/amd64/nvidia-docker.list | \
  sudo tee /etc/apt/sources.list.d/nvidia-docker.list
sudo apt-get update
```

#### Xenial ppc64le

```bash
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | \
  sudo apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/ubuntu16.04/ppc64el/nvidia-docker.list | \
  sudo tee /etc/apt/sources.list.d/nvidia-docker.list
sudo apt-get update
```

## Debian distributions

#### Stretch x86_64

```bash
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | \
  sudo apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/debian9/amd64/nvidia-docker.list | \
  sudo tee /etc/apt/sources.list.d/nvidia-docker.list
sudo apt-get update
```

## CentOS distributions

#### RHEL7 x86_64

```bash
curl -s -L https://nvidia.github.io/nvidia-docker/centos7/x86_64/nvidia-docker.repo | \
  sudo tee /etc/yum.repos.d/nvidia-docker.repo
```

#### RHEL7 ppc64le
```bash
curl -s -L https://nvidia.github.io/nvidia-docker/centos7/ppc64le/nvidia-docker.repo | \
  sudo tee /etc/yum.repos.d/nvidia-docker.repo
```
