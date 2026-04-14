#!/usr/bin/env bash
set -euo pipefail

usage() {
    cat <<'EOF'
Usage:
  ./scripts/release.sh <version> [options]

Examples:
  ./scripts/release.sh 1.2.2
  ./scripts/release.sh v1.2.2 --notes-file release-notes.md
  ./scripts/release.sh 1.2.2 --skip-tests

Options:
  --notes-file <path>   Show a suggested gh release command with this notes file
  --skip-tests          Skip go test ./...
  --create-tag          Create/update the local annotated tag after build
  -h, --help            Show this help message
EOF
}

log() {
    printf '[release] %s\n' "$*"
}

die() {
    printf '[release] %s\n' "$*" >&2
    exit 1
}

require_command() {
    command -v "$1" >/dev/null 2>&1 || die "Missing required command: $1"
}

normalize_version() {
    local version="${1#v}"
    if [[ ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+([.-][0-9A-Za-z.-]+)?$ ]]; then
        die "Invalid version: $1"
    fi
    printf '%s' "$version"
}

ensure_clean_tree() {
    if ! git diff --quiet || ! git diff --cached --quiet; then
        git status --short
        die "Working tree is not clean. Commit or stash your changes before preparing a release."
    fi
}

ensure_repo_root() {
    local root
    root="$(git rev-parse --show-toplevel 2>/dev/null)" || die "Not inside a git repository"
    cd "$root"
}

ensure_local_tag_state() {
    local head_commit local_tag_commit
    head_commit="$(git rev-parse HEAD)"

    if git rev-parse "$TAG" >/dev/null 2>&1; then
        local_tag_commit="$(git rev-list -n 1 "$TAG")"
        [[ "$local_tag_commit" == "$head_commit" ]] || die "Local tag $TAG already exists but does not point to HEAD"
        TAG_EXISTS_LOCAL=true
    fi
}

create_tag_if_needed() {
    if [[ "$CREATE_TAG" == false ]]; then
        log "Skipping local tag creation"
        return
    fi
    if [[ "$TAG_EXISTS_LOCAL" == true ]]; then
        log "Local tag $TAG already exists"
        return
    fi
    log "Creating local tag $TAG"
    git tag -a "$TAG" -m "release: $TAG"
}

print_next_steps() {
    log "Release artifacts are ready in dist/"
    printf '\n'
    printf 'Next steps you can run manually:\n'
    printf '  1. Review artifacts in dist/\n'
    if [[ "$CREATE_TAG" == false ]]; then
        printf '  2. Create local tag: git tag -a %s -m "release: %s"\n' "$TAG" "$TAG"
        printf '  3. Push branch/tag when you are ready:\n'
    else
        printf '  2. Push branch/tag when you are ready:\n'
    fi
    printf '     git push origin <branch>\n'
    printf '     git push origin %s\n' "$TAG"
    if [[ -n "$NOTES_FILE" ]]; then
        printf '  4. Create GitHub release:\n'
        printf '     gh release create %s dist/* --title "%s" --notes-file %s\n' "$TAG" "$TAG" "$NOTES_FILE"
    else
        printf '  4. Create GitHub release:\n'
        printf '     gh release create %s dist/* --title "%s" --generate-notes\n' "$TAG" "$TAG"
    fi
}

VERSION=""
TAG=""
NOTES_FILE=""
SKIP_TESTS=false
CREATE_TAG=false
TAG_EXISTS_LOCAL=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        --notes-file)
            shift
            NOTES_FILE="${1:-}"
            [[ -n "$NOTES_FILE" ]] || die "Missing value for --notes-file"
            ;;
        --skip-tests)
            SKIP_TESTS=true
            ;;
        --create-tag)
            CREATE_TAG=true
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        -*)
            die "Unknown option: $1"
            ;;
        *)
            if [[ -n "$VERSION" ]]; then
                die "Only one version argument is allowed"
            fi
            VERSION="$(normalize_version "$1")"
            ;;
    esac
    shift
done

[[ -n "$VERSION" ]] || {
    usage
    exit 1
}

TAG="v${VERSION}"

require_command git
require_command go

ensure_repo_root
ensure_clean_tree
ensure_local_tag_state

if [[ -n "$NOTES_FILE" && ! -f "$NOTES_FILE" ]]; then
    die "Release notes file not found: $NOTES_FILE"
fi

if [[ "$SKIP_TESTS" == false ]]; then
    log "Running go test ./..."
    go test ./...
else
    log "Skipping tests"
fi

log "Building release artifacts for ${VERSION}"
./build.sh "$VERSION"

create_tag_if_needed
print_next_steps

log "Local release preparation finished for ${TAG}"
