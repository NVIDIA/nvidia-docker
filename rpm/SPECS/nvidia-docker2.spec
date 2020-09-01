Name: nvidia-docker2
Version: %{version}
Release: %{release}
BuildArch: noarch
Group: Development Tools

Vendor: NVIDIA CORPORATION
Packager: NVIDIA CORPORATION <cudatools@nvidia.com>

Summary: nvidia-docker CLI wrapper
URL: https://github.com/NVIDIA/nvidia-docker
License: BSD

Source0: nvidia-docker
Source1: daemon.json
Source2: LICENSE

Conflicts: nvidia-docker < 2.0.0
Requires: nvidia-container-runtime >= %{runtime_version}
Requires: %{docker_version}

%description
Replaces nvidia-docker with a new implementation based on nvidia-container-runtime

%prep
cp %{SOURCE0} %{SOURCE1} %{SOURCE2} .

%install
mkdir -p %{buildroot}%{_bindir}
install -m 755 -t %{buildroot}%{_bindir} nvidia-docker
mkdir -p %{buildroot}/etc/docker
install -m 644 -t %{buildroot}/etc/docker daemon.json

%files
%license LICENSE
%{_bindir}/nvidia-docker
/etc/docker/daemon.json

%changelog
* Wed Jul 08 2020 NVIDIA CORPORATION <cudatools@nvidia.com> 2.4.0-1
- 09a01276 Update package license to match source license
- b9c70155 Update dependence on nvidia-container-runtime to 3.3.0

* Fri May 15 2020 NVIDIA CORPORATION <cudatools@nvidia.com> 2.3.0-1
- 0d3b049a Update build system to support multi-arch builds
- 8557216d Require new MIG changes
