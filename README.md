# 🐋 miniDocker

A custom container runtime and management system built in Go.

---

## Phase 1: Core Runtime

- [x] Process isolation (CLONE_NEWUTS, CLONE_NEWPID, CLONE_NEWNS)
- [x] Filesystem jailing (chroot + chdir)
- [x] Virtual filesystem mounting (/proc)
- [x] Resource constraints via Cgroups (memory + PID limits)
- [x] Parent → child re-execution pattern (/proc/self/exe)

## Phase 2: Container Management

- [ ] CLI separation (Manager vs Runtime, like docker-cli vs runc)
- [ ] Persistent storage directory (`/var/lib/minidocker/`) for container metadata
- [ ] Container state storage (ID, PID, StartTime, Config as JSON)
- [ ] `minidocker ps` — list active containers
- [ ] `minidocker inspect <id>` — view container metadata
- [ ] Volume support via bind-mounting (`-v /host/path:/container/path`)

## Phase 3: Advanced Features & OCI Compliance

- [ ] OCI image management — pull manifests from Docker Hub / OCI registries
- [ ] Download and unpack filesystem layers (.tar) as dynamic container root
- [ ] `minidocker exec <id> <command>` — join running container namespaces (setns)
- [ ] Networking — veth pairs + bridge networking for container IPs

## Phase 4: Lifecycle & Cleanup

- [ ] Graceful shutdown — unmount /proc, volumes, and delete cgroup dirs on exit
- [ ] Container logs — redirect stdout/stderr to a persistent host file

