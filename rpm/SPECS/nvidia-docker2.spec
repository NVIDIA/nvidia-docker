Name: nvidia-docker2
Version: %{version}
Release: %{release}
BuildArch: noarch
Group: Development Tools

Vendor: NVIDIA CORPORATION
Packager: NVIDIA CORPORATION <cudatools@nvidia.com>

Summary: nvidia-docker CLI wrapper
URL: https://github.com/NVIDIA/nvidia-docker
License: ASL 2.0

Source0: LICENSE

Conflicts: nvidia-docker < 2.0.0
Requires: nvidia-container-toolkit >= %{toolkit_version}

%description
A meta-package that allows installation flows expecting the nvidia-docker2
to be migrated to installing the NVIDIA Container Toolkit packages directly.
The wrapper script provided in earlier versions of this package should be
considered deprecated.

The nvidia-container-toolkit-base package provides an nvidia-ctk CLI that can be
used to update the docker config in-place to allow for the NVIDIA Container
Runtime to be used.

%prep
cp %{SOURCE0} .

%install

%files
%license LICENSE

%changelog
# As of 2.7.0-1 we generate the release information automatically
* %{release_date} NVIDIA CORPORATION <cudatools@nvidia.com> %{version}-%{release}
- As of 2.7.0-1 the package changelog is generated automatically. This means that releases since 2.7.0-1 all contain this same changelog entry updated for the version being released.
- Bump nvidia-container-toolkit dependency to %{toolkit_version}
- Docker dependency to %{docker_version}
