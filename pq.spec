%global project_version %(make print-vars | awk -F = '/RPM_VERSION/ {print $2}' | tr -d '\n')
%global debug_package %{nil}
 
Name:           pq
Version:        %{project_version} 
Release:        1%{?dist}
Summary:        A tool to manage podman quadlets 

License:        Apache-2.0
URL:            https://github.com/rgolangh/pq           
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang 
BuildRequires:  make 

%description
A tool to manage podman quadlets

%prep
%autosetup

%build
%make_build build INSTALL_PATH=./ GIT_COMMIT=000 GIT_COMMIT_DATE=000 GIT_TREE_STATE=clean GIT_VERSION=%{version} RPM_VERSION=%{version}

%install
#%make_install GIT_COMMIT=000 GIT_COMMIT_DATE=000 GIT_TREE_STATE=clean GIT_VERSION=%{version} INSTALL_PATH=%{name}
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}


%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}

%changelog
* Sat Oct 26 2024 Roy Golan <rgolan@redhat.com>
- 
