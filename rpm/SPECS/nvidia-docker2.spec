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
# Note: The -3 revision in the required toolkit version is to handle the released versions of
# The nvidia-container-toolkit 1.5.1 package. This can be replaced with '-1' in subsequent releases.
Requires: nvidia-container-toolkit > %{toolkit_version}-3
Requires: %{docker_version}

%description
Replaces nvidia-docker with a new implementation based on the NVIDIA Container Toolkit

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
%config /etc/docker/daemon.json

%changelog
* Mon Sep 06 2021 NVIDIA CORPORATION <cudatools@nvidia.com> 2.6.1-0.1.rc.1
- [BUILD] Allow for TAG to be specified in Makfile to match other projects
- Replace nvidia-container-runtime dependece with nvidia-container-toolit >= 1.5.2

* Thu Apr 29 2021 NVIDIA CORPORATION <cudatools@nvidia.com> 2.6.0-1
- Add dependence on nvidia-container-runtime >= 3.5.0
- Add Jenkinsfile for building packages

* Wed Sep 16 2020 NVIDIA CORPORATION <cudatools@nvidia.com> 2.5.0-1
- Bump version to v2.5.0
- Add dependence on nvidia-container-runtime >= 3.4.0
- Update readme to point to the official documentatio
- Add %config directive to daemon.json for RPM installations

* Wed Jul 08 2020 NVIDIA CORPORATION <cudatools@nvidia.com> 2.4.0-1
- 09a01276 Update package license to match source license
- b9c70155 Update dependence on nvidia-container-runtime to 3.3.0

* Fri May 15 2020 NVIDIA CORPORATION <cudatools@nvidia.com> 2.3.0-1
- 0d3b049a Update build system to support multi-arch builds
- 8557216d Require new MIG changes
