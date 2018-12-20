_The template below is mostly useful for bug reports and support questions. Feel free to remove anything which doesn't apply to you and add more information where it makes sense._

_Also, before reporting a new issue, please make sure that:_

- _You read carefully the [documentation and frequently asked questions](https://github.com/NVIDIA/nvidia-docker/wiki)._
- _You [searched](https://github.com/NVIDIA/nvidia-docker/issues?utf8=%E2%9C%93&q=is%3Aissue) for a similar issue and this is not a duplicate of an existing one._
- _This issue is not related to [NGC](https://github.com/NVIDIA/nvidia-docker/wiki/NGC), otherwise, please use the [devtalk forums](https://devtalk.nvidia.com/default/board/200/nvidia-gpu-cloud-ngc-users/) instead._
- _You went through the [troubleshooting](https://github.com/NVIDIA/nvidia-docker/wiki/Troubleshooting) steps._

---

### 1. Issue or feature description

### 2. Steps to reproduce the issue

### 3. Information to [attach](https://help.github.com/articles/file-attachments-on-issues-and-pull-requests/) (optional if deemed irrelevant)

 - [ ] Some nvidia-container information: `nvidia-container-cli -k -d /dev/tty info`
 - [ ] Kernel version from `uname -a`
 - [ ] Any relevant kernel output lines from `dmesg`
 - [ ] Driver information from `nvidia-smi -a`
 - [ ] Docker version from `docker version`
 - [ ] NVIDIA packages version from `dpkg -l '*nvidia*'` _or_ `rpm -qa '*nvidia*'`
 - [ ] NVIDIA container library version from `nvidia-container-cli -V`
 - [ ] NVIDIA container library logs (see [troubleshooting](https://github.com/NVIDIA/nvidia-docker/wiki/Troubleshooting))
 - [ ] Docker command, image and tag used
