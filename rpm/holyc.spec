Name:           holyc
Version:        0.1.1
Release:        1%{?dist}
Summary:        HolyC compiler written in Go

License:        GPL-3.0-or-later
URL:            https://github.com/Sh1n4r4/holyc-compiler
BuildArch:      x86_64

BuildRequires:  golang >= 1.21
Requires:       glibc

%description
HolyC is a programming language used in TempleOS. This compiler is
written in pure Go and compiles HolyC source code to x64 machine code.

Features:
- Lexical analysis (tokenizer)
- Parsing with AST generation
- x64 code generation
- Support for functions, classes, loops, conditions
- Command-line interface with multiple modes

%package -n holyc
Summary:        HolyC compiler binary
Requires:       %{name} = %{version}-%{release}

%prep
%setup -q -n holyc-go-%{version}

%build
export GOPATH=%{_builddir}/go
export GOFLAGS="-mod=vendor"

%go_build

%install
mkdir -p %{buildroot}%{_bindir}
install -m 755 holyc %{buildroot}%{_bindir}/holyc

mkdir -p %{buildroot}%{_docdir}/holyc
install -m 644 README.md %{buildroot}%{_docdir}/holyc/README.md
install -m 644 test.hc %{buildroot}%{_docdir}/holyc/test.hc

mkdir -p %{buildroot}%{_licensedir}
install -m 644 LICENSE %{buildroot}%{_licensedir}/holyc

%files
%{_bindir}/holyc
%doc %{_docdir}/holyc/README.md
%doc %{_docdir}/holyc/test.hc
%license %{_licensedir}/holyc

%changelog
* Sat Apr 04 2026 LT-SYAII <lt-syaii@users.noreply.github.com> - 0.1.1-1
- Add support for U0 type (void)
- Update parser to handle more HolyC features
- Improve error messages in parser
- Clean up redundant symbol type definitions
- Fix codegen issue where some operations were missing

* Sat Mar 28 2024 Sh1n4r4 <your.email@example.com> - 0.1.0-1
- Initial package
- HolyC compiler written in Go
