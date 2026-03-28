# Publishing Guide for HolyC Compiler

This guide explains how to publish the HolyC compiler to various Linux distribution repositories.

## Prerequisites

Before publishing, ensure you have:

1. A GitHub account and repository
2. GPG key for signing packages
3. Accounts on respective package repositories

## Setup GPG Key (for signing packages)

```bash
# Generate a new GPG key
gpg --full-generate-key

# List your keys
gpg --list-secret-keys --keyid-format LONG

# Export your public key
gpg --armor --export mitonohikari74@zohomail.com > gpg_public_key.asc

# Upload to keyserver
gpg --keyserver keyserver.ubuntu.com --send-keys YOUR_KEY_ID
```

## GitHub Release

### 1. Prepare Release

```bash
cd /home/karin/TempleOS/holyc-compiler

# Update version in files
# - main.go (version constant)
# - PKGBUILD (pkgver)
# - debian/changelog
# - rpm/holyc.spec
# - Makefile (VERSION)

# Build all release binaries
make release-all

# Create git tag
git tag -a v0.1.0 -m "Release version 0.1.0"
git push origin v0.1.0
```

### 2. Create GitHub Release

```bash
# Go to: https://github.com/Sh1n4r4/holyc-compiler/releases/new
# Tag: v0.1.0
# Title: HolyC Compiler v0.1.0
# Upload artifacts from dist/ folder
```

## Arch Linux (AUR)

### Publish to AUR

```bash
cd /home/karin/TempleOS/holyc-compiler

# Install aur-publish tools if needed
sudo pacman -S aur-publish

# Update PKGBUILD with correct source URL and checksums
updpkgsums PKGBUILD

# Publish to AUR
aur-publish
```

### Manual AUR Upload

```bash
# Clone AUR repo
git clone ssh://aur@aur.archlinux.org/holyc.git
cd holyc

# Copy files from holyc-compiler
cp ../holyc-compiler/PKGBUILD .
cp ../holyc-compiler/.SRCINFO .

# Generate .SRCINFO if not exists
makepkg --printsrcinfo > .SRCINFO

# Commit and push
git add .
git commit -m "Initial release v0.1.0"
git push origin master
```

## Debian/Ubuntu

### Option 1: Launchpad PPA

```bash
# Install packaging tools
sudo apt install devscripts debhelper build-essential

# Setup dput for PPA
mkdir -p ~/.dput.cf
cat >> ~/.dput.cf << EOF
[ppa]
fqdn = ppa.launchpad.net
method = ftp
incoming = ~your-username/ubuntu
login = anonymous
EOF

# Build source package
cd /home/karin/TempleOS/holyc-compiler
debuild -S -sa

# Sign and upload
dput ppa:your-username/ppa holyc_0.1.0_source.changes
```

### Option 2: Direct .deb Distribution

```bash
# Build .deb package
cd /home/karin/TempleOS/holyc-compiler
make pkg-deb

# Upload to GitHub Releases
# Attach .deb file to release
```

## Fedora/RHEL (COPR)

### Option 1: Fedora COPR

```bash
# Install COPR CLI
sudo dnf install copr-cli

# Build SRPM
cd /home/karin/TempleOS/holyc-compiler
make pkg-rpm

# Submit to COPR
copr-cli create holyc
copr-cli build holyc rpm-build/SRPMS/holyc-0.1.0-1.src.rpm
```

### Option 2: Direct .rpm Distribution

```bash
# Build RPM
cd /home/karin/TempleOS/holyc-compiler
make pkg-rpm

# Upload .rpm to GitHub Releases
```

## openSUSE (OBS)

### Open Build Service

```bash
# Install osc
sudo zypper install osc

# Configure osc
osc config --user your-username --pass your-password

# Create project
osc meta prj -e home:your-username:holyc

# Add package
cd /home/karin/TempleOS/holyc-compiler
osc checkout home:your-username:holyc
cp rpm/holyc.spec home:your-username:holyc/holyc/
cp -r * home:your-username:holyc/holyc/

# Commit
cd home:your-username:holyc/holyc
osc addremove
osc commit -m "Initial release v0.1.0"
```

## Snap Package

### Create Snap

```bash
# Install snapcraft
sudo snap install snapcraft --classic

# Create snapcraft.yaml
cd /home/karin/TempleOS/holyc-compiler
mkdir snap
cat > snap/snapcraft.yaml << EOF
name: holyc
version: '0.1.0'
summary: HolyC compiler written in Go
description: |
  HolyC compiler written in pure Go. Compiles HolyC
  source code to x64 machine code.

base: core22
confinement: strict

parts:
  holyc:
    plugin: go
    source: .
    go-buildtags: [netgo]

apps:
  holyc:
    command: bin/holyc
    plugs:
      - home
      - removable-media
EOF

# Build snap
snapcraft

# Upload to Snap Store
snapcraft login
snapcraft upload holyc_0.1.0_amd64.snap --release stable
```

## Flatpak Package

### Create Flatpak

```bash
# Install flatpak-builder
sudo apt install flatpak-builder

# Clone runtime
flatpak install flathub org.freedesktop.Platform//22.08
flatpak install flathub org.freedesktop.Sdk//22.08

# Create flatpak manifest
cat > holyc.json << EOF
{
    "app-id": "com.holyc.compiler",
    "runtime": "org.freedesktop.Platform",
    "runtime-version": "22.08",
    "sdk": "org.freedesktop.Sdk",
    "command": "holyc",
    "modules": [
        {
            "name": "holyc",
            "buildsystem": "simple",
            "build-commands": [
                "go build -o holyc .",
                "install -Dm755 holyc /app/bin/holyc"
            ],
            "sources": [
                {
                    "type": "git",
                    "url": "https://github.com/Sh1n4r4/holyc-compiler.git",
                    "tag": "v0.1.0"
                }
            ]
        }
    ]
}
EOF

# Build flatpak
flatpak-builder build-dir holyc.json --force-clean
flatpak build-export export-dir build-dir
flatpak build-update-repo export-dir
flatpak build-bundle export-dir holyc.flatpak com.holyc.compiler
```

## Post-Publishing Checklist

- [ ] GitHub release created with all binaries
- [ ] AUR package published and installable
- [ ] PPA/COPR builds successful
- [ ] README.md updated with installation instructions
- [ ] CHANGELOG.md updated
- [ ] Website/documentation updated
- [ ] Announce on social media/forums

## Announcement Template

```
HolyC Compiler v0.1.0 Released!

A HolyC compiler written in pure Go is now available for:
- Arch Linux (AUR)
- Debian/Ubuntu (PPA/.deb)
- Fedora/RHEL (COPR/.rpm)
- Direct binary downloads

Install with:
- Arch: yay -S holyc
- Debian: sudo add-apt-repository ppa:your-username/ppa && sudo apt install holyc
- Fedora: sudo dnf copr enable your-username/holyc && sudo dnf install holyc

GitHub: https://github.com/Sh1n4r4/holyc-compiler
Docs: https://github.com/Sh1n4r4/holyc-compiler#readme

#HolyC #TempleOS #Compiler #Go #OpenSource
```

## Troubleshooting

### Build fails on some architectures
- Check Go compatibility
- Update CGO settings
- Add architecture-specific patches

### Package rejected from repository
- Review repository guidelines
- Fix linting issues
- Improve documentation

### Users report installation issues
- Check dependencies
- Verify GPG signatures
- Update installation docs

---

For more help, see:
- [Arch Wiki - Creating packages](https://wiki.archlinux.org/title/Creating_packages)
- [Debian Policy Manual](https://www.debian.org/doc/debian-policy/)
- [Fedora Packaging Guidelines](https://docs.fedoraproject.org/en-US/packaging-guidelines/)
