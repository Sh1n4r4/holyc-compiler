# HolyC Compiler

**HolyC compiler written in pure Go** - Compiles HolyC source code to x64 machine code.

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/Sh1n4r4/holyc-compiler)](https://goreportcard.com/report/github.com/Sh1n4r4/holyc-compiler)

## Overview

HolyC is a programming language originally created by Terry Davis for TempleOS. This project is a modern reimplementation of a HolyC compiler written in Go, capable of compiling HolyC source code to native x64 machine code.

## Features

- **Lexical Analysis**: Complete tokenizer for HolyC syntax
- **Parser**: AST generation with operator precedence parsing
- **Code Generation**: x64 machine code emission
- **Multiple Output Modes**: Compile, tokenize, or parse only
- **Cross-Platform**: Works on all Linux distributions

### Supported Language Features

#### Types
- `void`, `Bool`
- `U8`, `U16`, `U32`, `U64`
- `I8`, `I16`, `I32`, `I64`
- `F64`

#### Operators
- Arithmetic: `+`, `-`, `*`, `/`, `%`
- Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
- Logical: `&&`, `||`, `!`
- Bitwise: `&`, `|`, `^`, `~`, `<<`, `>>`
- Assignment: `=`, `+=`, `-=`, `*=`, `/=`

#### Control Flow
- `if` / `else`
- `while`, `for`, `do/while`
- `switch` / `case` / `default`
- `break`, `continue`, `return`, `goto`

#### Declarations
- Functions with parameters
- Classes with members
- Variables (local and global)
- Extern functions

## Installation

### From Package Managers

#### Arch Linux / Manjaro
```bash
# Using AUR (with yay or paru)
yay -S holyc

# Or build from PKGBUILD
git clone https://github.com/Sh1n4r4/holyc-compiler.git
cd holyc-go
makepkg -si
```

#### Debian / Ubuntu / Linux Mint
```bash
# Download .deb package
wget https://github.com/Sh1n4r4/holyc-compiler/releases/download/v0.1.0/holyc_0.1.0_amd64.deb
sudo dpkg -i holyc_0.1.0_amd64.deb

# Or build from source
sudo apt install golang-go debhelper
git clone https://github.com/Sh1n4r4/holyc-compiler.git
cd holyc-go
make pkg-deb
sudo dpkg -i ../holyc_0.1.0_amd64.deb
```

#### Fedora / RHEL / openSUSE
```bash
# Download .rpm package
wget https://github.com/Sh1n4r4/holyc-compiler/releases/download/v0.1.0/holyc-0.1.0-1.x86_64.rpm
sudo dnf install holyc-0.1.0-1.x86_64.rpm

# Or build from source
sudo dnf install golang
git clone https://github.com/Sh1n4r4/holyc-compiler.git
cd holyc-go
make pkg-rpm
sudo dnf install rpm-build/RPMS/x86_64/holyc-0.1.0-1.x86_64.rpm
```

### From Source (Any Linux)
```bash
# Requires Go 1.21+
git clone https://github.com/Sh1n4r4/holyc-compiler.git
cd holyc-go
make build
sudo make install

# Verify installation
holyc --version
```

### Using Go Install
```bash
go install github.com/Sh1n4r4/holyc-compiler@latest
```

## Usage

```bash
# Compile a HolyC file
holyc program.hc

# Tokenize only (show all tokens)
holyc -t program.hc

# Parse only (show AST)
holyc -p program.hc

# Show help
holyc --help

# Show version
holyc --version
```

## Example

Create a file `hello.hc`:

```holyC
I64 Add(I64 a, I64 b)
{
    return a + b;
}

I64 Factorial(I64 n)
{
    if (n <= 1) {
        return 1;
    }
    return n * Factorial(n - 1);
}

class Point
{
    I64 x;
    I64 y;
};

I64 Main()
{
    I64 result;
    result = Add(10, 20);
    result = Factorial(5);
    return 0;
}
```

Compile it:

```bash
holyc hello.hc
```

Output:
```
+------------------------------------------+
|       HolyC Compiler                     |
|       Go Implementation v0.1.0           |
+------------------------------------------+

Compiling hello.hc...
Source size: 350 bytes
Phase 1: Lexical Analysis
Phase 2: Parsing
  Found 3 functions, 1 classes
Phase 3: Code Generation
  Generated 512 bytes of machine code
Output written to: hello.bin

Compilation successful!
```

## Development

### Prerequisites
- Go 1.21 or higher
- GNU Make (optional, for build scripts)
- Git

### Build from Source
```bash
git clone https://github.com/Sh1n4r4/holyc-compiler.git
cd holyc-go

# Build
make build

# Build with debug symbols
make build-debug

# Run tests
make test

# Run tests with coverage
make test-coverage

# Install locally
make install
```

### Project Structure
```
holyc-go/
├── main.go              # CLI entry point
├── lexer.go             # Lexical analyzer
├── parser.go            # Parser and AST
├── codegen.go           # x64 code generator
├── compiler.go          # Main compiler logic
├── Makefile             # Build system
├── LICENSE              # GPL-3.0 License
├── README.md            # This file
├── test.hc              # Example HolyC code
├── PKGBUILD             # Arch Linux package
├── debian/              # Debian/Ubuntu package files
└── rpm/                 # RPM package spec file
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Documentation

### Compiler Phases

1. **Lexical Analysis**: Source code -> Tokens
2. **Parsing**: Tokens -> Abstract Syntax Tree (AST)
3. **Code Generation**: AST -> x64 Machine Code

### Output Format

The compiler generates `.bin` files containing raw x64 machine code that can be:
- Loaded and executed directly (with proper loader)
- Linked with other object files
- Analyzed with disassemblers

## Limitations (v0.1.0)

- Single-pass compilation (no optimization)
- Limited standard library support
- Preprocessor directives (`#include`, `#define`) are skipped
- Array indexing partially implemented
- String literals in expressions partially supported
- Member access (`.` and `->`) partially implemented

## Roadmap

- [ ] Multi-pass compilation with optimizations
- [ ] Complete array support
- [ ] Full string literal support
- [ ] Standard library functions
- [ ] Debug info generation (DWARF)
- [ ] Better error messages with source snippets
- [ ] More HolyC features (templates, etc.)
- [ ] Linker support for standalone executables

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go best practices
- Write tests for new features
- Update documentation as needed
- Ensure `make test` passes

## License

This project is licensed under the **GNU General Public License v3.0 or later** (GPL-3.0-or-later).

```
Copyright (C) 2024  Sh1n4r4

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```

## Support

For issues and questions:
- **GitHub Issues**: [Create an issue](https://github.com/Sh1n4r4/holyc-compiler/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Sh1n4r4/holyc-compiler/discussions)

## Links

- [HolyC Wikipedia](https://en.wikipedia.org/wiki/HolyC)
- [TempleOS](https://en.wikipedia.org/wiki/TempleOS)
- [Go Programming Language](https://go.dev/)

---

Made with Go
