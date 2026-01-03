# Maintainer: Your Name <your.email@example.com>
pkgname=psql-stats
pkgver=1.0.0
pkgrel=1
pkgdesc="A terminal UI tool for viewing PostgreSQL database statistics"
arch=('x86_64' 'aarch64')
url="https://github.com/yourusername/psql-stats"
license=('MIT')
depends=('glibc')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::https://github.com/yourusername/$pkgname/archive/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
	cd "$pkgname-$pkgver"
	go build -trimpath -ldflags="-s -w" -o "$pkgname" .
}

package() {
	cd "$pkgname-$pkgver"
	install -Dm755 "$pkgname" "$pkgdir/usr/bin/$pkgname"
	install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"
}

