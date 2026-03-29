# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.1] - 2025-03-29

### Fixed
- Fixed install command not finding versions correctly

## [1.2.0] - 2025-03-28

### Added
- Enhanced version detection and matching
- Improved error handling and user feedback

### Changed
- Switched JSON parsing library for better performance

## [1.1.0] - 2024-11-18

### Added
- Interactive TUI for version list with keyboard navigation
- Real-time version filtering with `/` search
- One-click install/use/uninstall from TUI interface
- Progress indicators for downloads
- Spinner animations for long operations

### Changed
- Improved mirror registry abstraction with 3 implementations
- Enhanced version parsing with semantic version support

## [1.0.0] - 2024-09-09

### Added
- Core version management commands: `list`, `install`, `use`, `uninstall`
- Configuration management: `config get/set/list/unset`
- Project creation: `gvm new` command
- Self-upgrade: `gvm upgrade` command
- Multi-mirror support (official, Aliyun, USTC, HUST, NJU, etc.)
- Cross-platform support (Linux, macOS)
- Symlink-based version switching mechanism

## [Unreleased]

### Planned
- `.gvmrc` project-level version isolation
- `gvm doctor` environment diagnostic tool
- Shell completion (bash, zsh, fish)
- Project template system
- Plugin system architecture
- Intelligent version recommendation