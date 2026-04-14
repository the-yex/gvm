#!/usr/bin/env bash
set -euo pipefail

resolve_release() {
    if [[ -n "${1:-}" ]]; then
        printf '%s' "$1"
        return
    fi
    if [[ -n "${GVM_RELEASE:-}" ]]; then
        printf '%s' "${GVM_RELEASE}"
        return
    fi
    if git_tag="$(git describe --tags --abbrev=0 2>/dev/null)"; then
        printf '%s' "${git_tag#v}"
        return
    fi
    echo "Unable to determine release version. Pass it explicitly: ./build.sh 1.2.3" >&2
    exit 1
}

RELEASE="$(resolve_release "${1:-}")"

TARGETS=(
    "darwin_amd64" "darwin_arm64"
    "linux_386" "linux_amd64" "linux_arm" "linux_arm64" "linux_s390x" "linux_riscv64"
#    "windows_386" "windows_amd64" "windows_arm" "windows_arm64"
)

OUTPUT_DIR="./dist"
SHA_FILE="${OUTPUT_DIR}/sha256sum.txt"
VERSION_PKG="github.com/the-yex/gvm/internal/consts"

rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

function package() {
    local release="$1"
    local osarch="$2"
    local os="${osarch%%_*}"
    local arch="${osarch##*_}"

    local tmp_bin="./gvm"
    if [[ "$os" == "windows" ]]; then
        tmp_bin="${tmp_bin}.exe"
    fi

    echo "[→] Building ${os}-${arch}..."

    GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 GO111MODULE=on GOPROXY="https://goproxy.cn,direct" \
        go build -ldflags "-X ${VERSION_PKG}.Version=${release}" -o "$tmp_bin" .

    if [[ "$os" == "windows" ]]; then
        local pkg_name="gvm${release}.${os}-${arch}.zip"
        zip -j "${OUTPUT_DIR}/${pkg_name}" "$tmp_bin" > /dev/null
    else
        local pkg_name="gvm${release}.${os}-${arch}.tar.gz"
        tar -czf "${OUTPUT_DIR}/${pkg_name}" -C "$(dirname "$tmp_bin")" "$(basename "$tmp_bin")"
    fi

    shasum -a 256 "${OUTPUT_DIR}/${pkg_name}" > "${OUTPUT_DIR}/${pkg_name}.sha256"

    rm -f "$tmp_bin"
    echo "[✓] $pkg_name built"
}


if [[ ! "$RELEASE" =~ ^[0-9]+\.[0-9]+\.[0-9]+([.-][0-9A-Za-z.-]+)?$ ]]; then
    echo "Invalid release version: ${RELEASE}" >&2
    echo "Usage: ./build.sh 1.2.3" >&2
    exit 1
fi

echo "Building gvm version ${RELEASE} for multiple platforms..."
echo "Output dir: $OUTPUT_DIR"

for target in "${TARGETS[@]}"; do
    package "$RELEASE" "$target" &
done

wait

cat ${OUTPUT_DIR}/*.sha256 > "$SHA_FILE"
rm -f ${OUTPUT_DIR}/*.sha256

go clean
echo "All builds done. SHA256 sums written to $SHA_FILE"
