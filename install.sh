#!/usr/bin/env bash
set -euo pipefail

GVM_RELEASE="1.1.1"
GVM_HOME="${HOME}/.gvm"
SOURCE="github.com"
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
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
    shift
done
# 安装 gvm
install_gvm() {
    local os arch dest_file url
    os=$(get_os)
    arch=$(get_arch)
    dest_file="${GVM_HOME}/gvm${GVM_RELEASE}.${os}-${arch}.tar.gz"
    url="https://${SOURCE}/the-yex/gvm/releases/download/v${GVM_RELEASE}/gvm${GVM_RELEASE}.${os}-${arch}.tar.gz"

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
}
main() {
    install_gvm
}

main