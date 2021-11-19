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
Requires: nvidia-container-toolkit >= %{toolkit_version}
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
# As of 2.7.0-1 we generate the release information automatically
* %{release_date} NVIDIA CORPORATION <cudatools@nvidia.com> %{version}-%{release}
- As of 2.7.0-1 the package changelog is generated automatically. This means that releases since 2.7.0-1 all contain this same changelog entry updated for the version being released.
- Bump nvidia-container-toolkit dependency to %{toolkit_version}
- Docker dependency to %{docker_version}
