#!/usr/bin/env bash
set -euo pipefail

GVM_HOME="${HOME}/.gvm"
SOURCE="github.com"
RELEASE_VERSION=""

fetch() {
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$1"
        return
    fi
    wget -qO- "$1"
}

normalize_version() {
    local version="$1"
    version="${version#v}"
    printf '%s' "$version"
}

release_api_candidates() {
    case "${SOURCE}" in
        gitee.com)
            printf '%s\n' \
                "https://gitee.com/api/v5/repos/the-yex/gvm/releases/latest" \
                "https://api.github.com/repos/the-yex/gvm/releases/latest"
            ;;
        *)
            printf '%s\n' \
                "https://api.github.com/repos/the-yex/gvm/releases/latest" \
                "https://gitee.com/api/v5/repos/the-yex/gvm/releases/latest"
            ;;
    esac
}

resolve_release_version() {
    if [[ -n "${RELEASE_VERSION}" ]]; then
        normalize_version "${RELEASE_VERSION}"
        return
    fi

    local release_api tag response
    while IFS= read -r release_api; do
        if [[ -z "${release_api}" ]]; then
            continue
        fi
        if response="$(fetch "${release_api}" 2>/dev/null)"; then
            tag="$(printf '%s' "${response}" | tr -d '\r' | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)"
            if [[ -n "${tag}" ]]; then
                normalize_version "${tag}"
                return
            fi
        fi
    done < <(release_api_candidates)

    echo "Failed to resolve latest release version from release APIs" >&2
    echo "You can pin a version manually with: ./install.sh --version 1.2.3" >&2
    exit 1
}
# 获取系统架构
get_arch() {
    case "$(uname -m)" in
        x86_64 | amd64) echo "amd64" ;;
        i386 | i486 | i586) echo "386" ;;
        aarch64 | arm64) echo "arm64" ;;
        armv6l | armv7l) echo "arm" ;;
        s390x) echo "s390x" ;;
        riscv64) echo "riscv64" ;;
        *)
            echo "Unsupported architecture: $(uname -m)" >&2
            exit 1
            ;;
    esac
}

function get_os() {
    echo $(uname -s | awk '{print tolower($0)}')
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --source)
            shift
            if [[ "$1" == "gitee" ]]; then
                SOURCE="gitee.com"
            elif [[ "$1" == "github" ]]; then
                SOURCE="github.com"
            else
                echo "Unknown source: $1. Use 'gitee' or 'github'."
                exit 1
            fi
            shift
            ;;
        --version)
            shift
            RELEASE_VERSION="$(normalize_version "${1:-}")"
            if [[ -z "${RELEASE_VERSION}" ]]; then
                echo "Missing value for --version" >&2
                exit 1
            fi
            shift
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done
# 安装 gvm
install_gvm() {
    local os arch dest_file url resolved_version
    os=$(get_os)
    arch=$(get_arch)
    resolved_version="$(resolve_release_version)"
    dest_file="${GVM_HOME}/gvm${resolved_version}.${os}-${arch}.tar.gz"
    url="https://${SOURCE}/the-yex/gvm/releases/download/v${resolved_version}/gvm${resolved_version}.${os}-${arch}.tar.gz"

    echo "[0/3] Resolved gvm version ${resolved_version}"
    echo "[1/3] Downloading ${url}"
    mkdir -p "${GVM_HOME}"
    rm -f "${dest_file}"

    if command -v wget >/dev/null 2>&1; then
        wget -q -O "${dest_file}" "${url}"
    else
        curl -sSL -o "${dest_file}" "${url}"
    fi

    echo "[2/3] Installing gvm to ${GVM_HOME}"
    tar -xzf "${dest_file}" -C "${GVM_HOME}"
    chmod +x "${GVM_HOME}/gvm"
    rm -f "${dest_file}"

    echo "[3/3] Configuring shell environment"
    local shell_config
    for shell_config in "${HOME}/.bashrc" "${HOME}/.zshrc"; do
        if [ -f "$shell_config" ] || [ -w "$HOME" ]; then
            if ! grep -q "gvm shell setup" "$shell_config"; then
                cat >>"$shell_config" <<-'EOF'

# gvm shell setup
# 清理旧的 Go 环境变量
unset GOROOT
unset GO_ROOT
unset GOPATH
# 从 PATH 中移除任何旧的 go/bin
export PATH=$(echo "$PATH" | awk -v RS=: -v ORS=: '$0 !~ /\/go\/bin/ && $0 !~ /\/usr\/local\/go\/bin/' | sed 's/:$//')

# 设置 gvm 的环境变量
export GVM_HOME="${HOME}/.gvm"
export GOROOT="${GVM_HOME}/go"
[ -z "$GOPATH" ] && export GOPATH="${HOME}/go"

case ":$PATH:" in
  *":${GVM_HOME}:"*) ;;
  *) export PATH="${GVM_HOME}:${GOROOT}/bin:${GOPATH}/bin:$PATH" ;;
esac

EOF
            fi
        fi
    done

    echo -e "\nInstallation completed. Please restart your terminal or source your shell configuration file."
    echo "Installed version: ${resolved_version}"
}
main() {
    install_gvm
}

main
