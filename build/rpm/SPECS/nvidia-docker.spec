Name: %{name}
Version: %{version}
Release: %{revision}
BuildArch: %{architecture}
Group: Development Tools

Vendor: %{vendor}
Packager: %{vendor} <%{email}>

Summary: NVIDIA Docker container tools
URL: https://github.com/NVIDIA/nvidia-docker
License: BSD

Source0: %{name}_%{version}_%{architecture}.tar.xz
Source1: %{name}.service
Source2: LICENSE

%{?systemd_requires}
BuildRequires: systemd
Requires: libcap

%define nvidia_docker_user %{name}
%define nvidia_docker_driver %{name}
%define nvidia_docker_root /var/lib/nvidia-docker

%description
NVIDIA Docker provides utilities to extend the Docker CLI allowing users
to build and run GPU applications as lightweight containers.

%prep
%autosetup -n %{name}
cp %{SOURCE1} %{SOURCE2} .

%install
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_unitdir}
mkdir -p %{buildroot}%{nvidia_docker_root}
install -m 755 -t %{buildroot}%{_bindir} nvidia-docker
install -m 755 -t %{buildroot}%{_bindir} nvidia-docker-plugin
install -m 644 -t %{buildroot}%{_unitdir} %{name}.service

%files
%license LICENSE
%dir %{nvidia_docker_root}
%{_bindir}/*
%{_unitdir}/*

%post
if [ $1 -eq 1 ]; then
    echo "Configuring user"
    id -u %{nvidia_docker_user} >/dev/null 2>&1 || \
    useradd -r -M -d %{nvidia_docker_root} -s /usr/sbin/nologin -c "NVIDIA Docker plugin" %{nvidia_docker_user}
fi
echo "Setting up permissions"
chown %{nvidia_docker_user}: %{nvidia_docker_root}
setcap cap_fowner+pe %{_bindir}/nvidia-docker-plugin
%systemd_post %{name}

%preun
if [ $1 -eq 0 ]; then
    echo "Purging NVIDIA volumes"
    docker volume ls | awk -v drv=%{nvidia_docker_driver} '{if ($1 == drv) print $2}' | xargs -r docker volume rm ||
        echo "Failed to remove NVIDIA volumes, ignoring"
    find %{nvidia_docker_root} ! -wholename %{nvidia_docker_root} -type d -empty -delete
fi
%systemd_preun %{name}

%postun
if [ $1 -eq 0 ]; then
    id -u %{nvidia_docker_user} >/dev/null 2>&1 && \
    userdel %{nvidia_docker_user}
fi
%systemd_postun_with_restart %{name}

%changelog
* Fri Mar 03 2017 NVIDIA CORPORATION <digits@nvidia.com> 1.0.1-1
- Support for Docker 17.03 including EE and CE (Closes: #323, #324)
- Load UVM unconditionally
- Fix Docker argument parsing (Closes: #295)
- Fix images pull output (Closes: #310)

* Wed Jan 18 2017 NVIDIA CORPORATION <digits@nvidia.com> 1.0.0-1
- Support for Docker 1.13
- Fix CPU affinity reporting on systems where NUMA is disabled (Closes: #198)
- Fix premature EOF in the remote API responses (Closes: #123)
- Add support for the VolumeDriver.Capabilities plugin endpoint
- Enable ppc64le library lookup (Closes: #194)
- Fix parsing of DOCKER_HOST for unix domain sockets (Closes: #119)

* Fri Jun 17 2016 NVIDIA CORPORATION <digits@nvidia.com> 1.0.0~rc.3-1
- Support for Docker 1.12
- Add volume mount options support to the nvidia package
- Export the nvidia-uvm-tools device
- Provide the libcuda.so symlink as part of the driver volume (Closes: #103)
- Use relative symlinks inside the volumes
- Disable CUDA unified memory

* Sat May 28 2016 NVIDIA CORPORATION <digits@nvidia.com> 1.0.0~rc.2-1
- Allow UUIDs to be used in NV_GPU and docker/cli RestAPI endpoint
- Change the plugin usage with version information (Closes: #90)
- Remove the volume setup command (Closes: #96)
- Add support for the Pascal architecture

* Tue May 03 2016 NVIDIA CORPORATION <digits@nvidia.com> 1.0.0~rc-1
- Add /docker/cli/json RestAPI endpoint (Closes: #39, #91)
- Fix support for Docker 1.9 (Closes: #83)
- Handle gracefully devices unsupported by NVML (Closes: #40)
- Improve error reporting
- Support for Docker 1.11 (Closes: #89, #84, #77, #73)
- Add NVIDIA Docker version output
- Improve init scripts and add support for systemd
- Query CPU affinity through sysfs instead of NVML (Closes: #65)
- Load UVM before anything else

* Mon Mar 28 2016 NVIDIA CORPORATION <digits@nvidia.com> 1.0.0~beta.3-1
- Remove driver hard dependency (NVML)
- Improve error handling and REST API output
- Support for 364 drivers
- Preventive removal of the plugin socket

* Mon Mar 07 2016 NVIDIA CORPORATION <digits@nvidia.com> 1.0.0~beta.2-1
- Support for Docker 1.10 (Closes: #46)
- Support for Docker plugin API v1.2
- Support for 361 drivers
- Add copy strategy for cross-device volumes (Closes: #47)

* Mon Feb 08 2016 NVIDIA CORPORATION <digits@nvidia.com> 1.0.0~beta-1
- Initial release (Closes: #33)
