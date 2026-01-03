# Packaging Guide

This document explains how to package `psql-stats` for different distributions.

## AUR (Arch User Repository)

1. **Create AUR package repository:**
   ```bash
   git clone ssh://aur@aur.archlinux.org/psql-stats.git
   cd psql-stats
   ```

2. **Copy PKGBUILD:**
   ```bash
   cp /path/to/psql-stats/PKGBUILD .
   ```

3. **Update PKGBUILD:**
   - Set `pkgver` to the current version
   - Update `source` URL to point to your GitHub release
   - Calculate `sha256sums`:
     ```bash
     makepkg -g
     ```

4. **Test build:**
   ```bash
   makepkg -s
   ```

5. **Submit to AUR:**
   ```bash
   git add PKGBUILD .SRCINFO
   git commit -m "Initial package"
   git push
   ```

## Debian/Ubuntu

1. **Install build dependencies:**
   ```bash
   sudo apt-get install build-essential devscripts debhelper golang-go
   ```

2. **Build package:**
   ```bash
   dpkg-buildpackage -us -uc
   ```

3. **Create PPA (optional):**
   - Upload to Launchpad PPA
   - Or use `dput` to upload to Debian

## Flatpak

1. **Install Flatpak and Builder:**
   ```bash
   sudo apt install flatpak flatpak-builder
   flatpak install org.freedesktop.Platform//23.08 org.freedesktop.Sdk//23.08
   ```

2. **Build:**
   ```bash
   flatpak-builder build com.github.yourusername.psql-stats.yml
   ```

3. **Test:**
   ```bash
   flatpak-builder --run build com.github.yourusername.psql-stats.yml psql-stats
   ```

4. **Publish to Flathub:**
   - Fork https://github.com/flathub/flathub
   - Add your manifest
   - Submit PR

## Homebrew (macOS)

1. **Create formula:**
   Create `Formula/psql-stats.rb` in your tap:
   ```ruby
   class PsqlStats < Formula
     desc "Terminal UI tool for PostgreSQL statistics"
     homepage "https://github.com/yourusername/psql-stats"
     url "https://github.com/yourusername/psql-stats/archive/v1.0.0.tar.gz"
     sha256 "YOUR_SHA256"
     license "MIT"
   
     depends_on "go" => :build
   
     def install
       system "make", "build"
       bin.install "psql-stats"
     end
   
     test do
       system "#{bin}/psql-stats", "--help"
     end
   end
   ```

2. **Test:**
   ```bash
   brew install --build-from-source Formula/psql-stats.rb
   ```

## Version Management

Update version in:
- `PKGBUILD` (pkgver)
- `debian/changelog`
- `com.github.yourusername.psql-stats.yml` (tag)
- `README.md`

## Release Process

1. Update version numbers
2. Create git tag: `git tag v1.0.0`
3. Push tag: `git push --tags`
4. GitHub Actions will build and release binaries
5. Update package files with new version
6. Submit to package repositories

