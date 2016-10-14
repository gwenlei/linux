Name: cstemplate		
Version: 0.1
Release: 1%{?dist}
Summary: manipulate cloudstack templates

Group: Application/test		
License: GPLv3+	
Source0:$RPM_SOURCE_DIR/cstemplate-0.1.tar	

%description
manipulate cloudstack templates

%prep
%setup -q

%build
cd $RPM_BUILD_DIR/cstemplate
go build cstemplate.go

%install
cd $RPM_BUILD_DIR/cstemplate
mkdir -p $RPM_BUILD_ROOT%{_bindir} $RPM_BUILD_ROOT%_sysconfdir
install -D -m 0755 cstemplate $RPM_BUILD_ROOT%{_bindir}
install -D -m 0644 cstemplate.ini $RPM_BUILD_ROOT%_sysconfdir

%clean
rm -rf $RPM_BUILD_ROOT

%files
%defattr(-,root,root,-)
%{_bindir}/cstemplate
%_sysconfdir/cstemplate.ini


%post
/sbin/install-info %{_infodir}/%{name}.info %{_infodir}/dir || :

%preun
if [ $1 = 0 ] ; then
/sbin/install-info --delete %{_infodir}/%{name}.info %{_infodir}/dir || :
fi

%define __debug_install_post   \
   %{_rpmconfigdir}/find-debuginfo.sh %{?_find_debuginfo_opts} "%{_builddir}/%{?buildsubdir}"\
%{nil}

%changelog
* Thu Oct 13 2016 onecloud <mail@onecloud.cn> 0.1-1
- Initial version of the package


