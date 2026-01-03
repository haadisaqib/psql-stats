# Quick Start: Packaging Your App

## 1. Update Package Information

Before submitting, update these files with your information:

- **PKGBUILD**: Change `yourusername` to your GitHub username
- **debian/control**: Update Maintainer email
- **com.github.yourusername.psql-stats.yml**: Change app-id and URLs
- **README.md**: Update GitHub URLs

## 2. AUR (Arch Linux)

```bash
# Clone AUR repository (create it first on AUR website)
git clone ssh://aur@aur.archlinux.org/psql-stats.git aur-psql-stats
cd aur-psql-stats

# Copy your files
cp ../PKGBUILD .
cp ../.SRCINFO .

# Generate .SRCINFO (if you have makepkg)
makepkg --printsrcinfo > .SRCINFO

# Test build
makepkg -s

# Commit and push
git add PKGBUILD .SRCINFO
git commit -m "Initial package release"
git push
```

## 3. Debian/Ubuntu Package

```bash
# Build the package
dpkg-buildpackage -us -uc

# This creates: ../psql-stats_1.0.0-1_amd64.deb

# Test install
sudo dpkg -i ../psql-stats_1.0.0-1_amd64.deb
```

## 4. Flatpak

```bash
# Install Flatpak builder
sudo apt install flatpak flatpak-builder

# Install runtime
flatpak install org.freedesktop.Platform//23.08 org.freedesktop.Sdk//23.08

# Build
flatpak-builder build com.github.yourusername.psql-stats.yml

# Test
flatpak-builder --run build com.github.yourusername.psql-stats.yml psql-stats

# To install locally
flatpak-builder --install --user build com.github.yourusername.psql-stats.yml
```

## 5. GitHub Releases (Universal)

```bash
# Tag your release
git tag v1.0.0
git push --tags

# GitHub Actions will automatically build and release binaries
# Check .github/workflows/release.yml
```

## 6. Homebrew (macOS)

Create a tap repository:
```bash
# Create formula in your tap
# See PACKAGING.md for formula template
```

## Testing Checklist

- [ ] AUR: `makepkg -s` succeeds
- [ ] Debian: `dpkg-buildpackage` succeeds
- [ ] Flatpak: Builds and runs
- [ ] Binary works on target system
- [ ] All URLs point to correct repository
- [ ] Version numbers are consistent

## Next Steps

1. Create GitHub releases with binaries
2. Submit AUR package
3. Submit to Flathub (if desired)
4. Create Homebrew tap
5. Update README with installation instructions

