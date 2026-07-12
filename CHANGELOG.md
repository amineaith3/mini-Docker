# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- **Dynamic Root Filesystem**: Implemented dynamic path resolution for the container root filesystem. It now reads from the `MINIDOCKER_ROOTFS` environment variable or falls back to `/var/lib/minidocker/rootfs`, replacing hardcoded user-specific paths.

### Fixed
- **Mount Namespace Leakage**: Moved `setupVolumeMounts` into the child process. Volumes are now securely mounted *after* the process enters its isolated mount namespace, preventing container mounts from polluting the host machine.
- **Cgroup Isolation**: Container IDs are now generated prior to cgroup creation. Containers are placed in isolated cgroup directories (`/sys/fs/cgroup/miniDocker-<id>`), rather than sharing a single cgroup and conflicting with one another during exit cleanup.
- **State Initialization**: Corrected `GenJSON` behavior to unconditionally accept the provided ID so that the saved metadata JSON perfectly aligns with the generated container ID used for resource tracking.

### Changed
- **Cgroup Hierarchy**: Flattened the cgroup hierarchy to avoid `permission denied` errors when writing to resource controllers (like `memory.max`) under Cgroup v2. Containers are now placed directly at the root cgroup level where memory controllers are delegated by default.
