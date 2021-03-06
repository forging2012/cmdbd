# =============================================================================
%define		name	cmdbd
%define		version	3.2.1
%define		release	1
%define		branch  master
%define		gecos	CMDBd Service
%define		summary	Configuration Management Database Daemon
%define		author	John Scherff <jscherff@24hourfit.com>
%define		package	github.com/jscherff/%{name}
%define		gopath	%{_builddir}/go
%define		docdir	%{_docdir}/%{name}-%{version}
%define		logdir	%{_var}/log/%{name}
%define		syslib	%{_prefix}/lib/systemd/system
%define		confdir %{_sysconfdir}/%{name}
# =============================================================================

Name:		%{name}
Version:	%{version}
Release:	%{release}%{?dist}
Summary:	%{summary}

Group:		Applications/System
License:	ASL 2.0
URL:		https://www.24hourfitness.com
Vendor:		24 Hour Fitness, Inc.
Prefix:		%{_sbindir}
Packager: 	%{packager}
BuildRoot:	%{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)
Distribution:	el

BuildRequires:    golang >= 1.8.0
Requires(pre):    %{_sbindir}/useradd, %{_bindir}/getent
Requires(postun): %{_sbindir}/userdel

%description
The Configuration Management Database Daemon, %{name}, is a lightweight HTTP
daemon that provides a REST API for clients installed on Windows endpoints.
The clients collect information about attached devices and send it to the
server for storage in the database. Clients can register attached devices
with the server, obtain unique serial numbers from the server for devices
that support serial number configuration, perform audits against previous
device configurations, and report any configuration changes found during
the audit to the server for later analysis.

%prep

%build

  export GOPATH=%{gopath}
  export GIT_DIR=%{gopath}/src/%{package}/.git

  go get %{package}
  git checkout %{branch}

  go build -ldflags='-X main.version=%{version}-%{release}' %{package}
  go build -ldflags='-X main.version=%{version}-%{release}' %{package}/bcrypt

%install

  test %{buildroot} != / && rm -rf %{buildroot}/*

  mkdir -p %{buildroot}{%{_sbindir},%{_bindir}}
  mkdir -p %{buildroot}{%{confdir},%{syslib},%{logdir},%{docdir}}

  install -s -m 755 %{_builddir}/%{name} %{buildroot}%{_sbindir}/
  install -s -m 755 %{_builddir}/bcrypt %{buildroot}%{_bindir}/
  install -m 640 %{gopath}/src/%{package}/deploy/ddl/%{name}.sql %{buildroot}%{docdir}/
  install -m 640 %{gopath}/src/%{package}/deploy/dml/reset.sql %{buildroot}%{docdir}/
  install -m 644 %{gopath}/src/%{package}/deploy/svc/* %{buildroot}%{syslib}/
  install -m 644 %{gopath}/src/%{package}/{LICENSE,*.md} %{buildroot}%{docdir}/

  cp -R %{gopath}/src/%{package}/config/* %{buildroot}%{confdir}/

%clean

  test %{buildroot} != / && rm -rf %{buildroot}/*
  test %{_builddir} != / && rm -rf %{_builddir}/*

%files

  %defattr(-,root,root)
  %license %{docdir}/LICENSE
  %{_sbindir}/*
  %{_bindir}/*
  %{syslib}/*
  %{docdir}/*

  %defattr(640,%{name},%{name},750)
  %config %{confdir}/*

  %defattr(644,%{name},%{name},755)
  %{logdir}

%pre

  # Tasks to perform FROM NEW RPM before install (1) or upgrade (2)

  case ${1} in

    1)
      %{_sbindir}/useradd -Mrd %{_sbindir} -c '%{gecos}' -s /sbin/nologin %{name}
      ;;

    2)
      systemctl --quiet is-active %{name} && systemctl --quiet stop %{name}
      systemctl --quiet is-enabled %{name} && systemctl --quiet disable %{name}
      ;;

  esac

  : Force zero return code

%post

  # Tasks to perform FROM NEW RPM after install (1) or upgrade (2)

  case ${1} in

    1)
      systemctl --quiet is-enabled %{name} || systemctl --quiet enable %{name} 
      ;;

    2)
      systemctl --quiet is-enabled %{name} || systemctl --quiet enable %{name} 
      systemctl --quiet is-active %{name} || systemctl --quiet start %{name}
      ;;

  esac

  : Force zero return code

%preun

  # Tasks to perform FROM OLD RPM before uninstall (0) or upgrade (1)

  case ${1} in

    0)
      systemctl --quiet is-active %{name} && systemctl --quiet stop %{name}
      systemctl --quiet is-enabled %{name} && systemctl --quiet disable %{name}
      ;;

    1)
      ;;

  esac

  : Force zero return code

%postun

  # Tasks to perform FROM OLD RPM after uninstall (0) or upgrade (1)

  case ${1} in

    0)
      %{_sbindir}/userdel %{name}
      test %{logdir} != / && rm -rf %{logdir}
      ;;

    1)
      ;;

  esac

  : Force zero return code

%changelog
* Fri Jan 26 2018 - jscherff@gmail.com
- Separated DML and DDL
- Modified RPM spec file to enhance GO build process
* Wed Jan 17 2018 - jscherff@gmail.com
- Comprehensive refactor to make code resusable and easier to maintain
- Converted model to lightweight ORM using sqlx
- Segregated server components into 'server' package
- Segregated common services into 'service' package
- Segregated store components into 'store' package
- Segregated API components into 'api' package
- Created separate v1 and v2 APIs for backward compatibility
* Mon Nov 13 2017 - jscherff@gmail.com
- Modified queries to use tables directly versus views
- Added DATETIME columns to inserts with time.Now() as value
- Modified Loc (location) database parameter to 'Local'
- Removed unnecessary views from DDL
* Wed Nov 8 2017 - jscherff@gmail.com
- Added cmdb_users table for authentication
- Added authentication API to support basic authentication
- Added authentication JWT support for protected API endpoints
- Added authentication JWT validation middleware to protect API endpoints
* Thu Oct 19 2017 - jscherff@gmail.com
- Added SQL script to truncate all tables
* Fri Oct 13 2017 - jscherff@gmail.com
- Refactored and streamlined
- Added API endpoints for device information lookups
* Mon Oct 9 2017 - jscherff@gmail.com
- Modified table, view, and stored procedure names
- Added column to each table for the JSON object
- Modified changes column in changes table to be datatype JSON
* Sat Oct 7 2017 - jscherff@gmail.com
- Added v1 prefix to URLs and handlers
* Sat Sep 30 2017 - jscherff@gmail.com
- Tightened file permissions mode on config.json
