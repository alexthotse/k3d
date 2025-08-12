# Using Podman instead of Docker

K3d has native support for Podman as a container runtime, starting from k3d v5.x.x (refer to the specific version when this was introduced). This means you can use Podman directly without relying solely on its Docker API compatibility layer for k3d to function.

Podman v4.0 and higher is recommended.

!!! tip "Using Podman is Easy!"
    The simplest way to use Podman with k3d is to specify it as the runtime:
    *   **Via command-line flag:** Use the `--runtime podman` global flag.
        ```bash
        k3d cluster create mypodmancluster --runtime podman
        ```
    *   **Via configuration file:** Set the `runtime` field in your k3d configuration file.
        ```yaml
        # k3d-config.yaml
        apiVersion: k3d.io/v1alpha5
        kind: Simple
        metadata:
          name: mypodmancluster
        runtime: podman # <--- Add this line
        servers: 1
        agents: 0
        ```
        Then create the cluster with `k3d cluster create --config k3d-config.yaml`.

    When using `--runtime podman`, k3d will attempt to interact with Podman directly.

## Essential Podman Setup

Even with native support, some Podman-specific setup might be necessary depending on your environment, especially for rootless Podman or specific operating systems.

### 1. Ensure Podman Socket is Active

K3d interacts with Podman via its API socket.

**For rootful Podman (usually system-wide):**
Ensure the Podman system socket is enabled and active:
```bash
sudo systemctl enable --now podman.socket
# Or, to start the service daemonless (less common for k3d use):
# sudo podman system service --time=0 &
```

You might also need to disable the service timeout for long-running k3d operations, although this is less critical when k3d manages the lifecycle. If you encounter timeouts:
```bash
# Optional: Disable service timeout
sudo mkdir -p /etc/containers/containers.conf.d
echo 'service_timeout=0' | sudo tee /etc/containers/containers.conf.d/timeout.conf > /dev/null
```

**For rootless Podman (user-specific):**
Ensure the Podman user socket is enabled and active:
```bash
systemctl --user enable --now podman.socket
# Or, to start the service daemonless:
# podman system service --time=0 &
```
The socket path is typically `$XDG_RUNTIME_DIR/podman/podman.sock`.

### 2. Cgroup Configuration (Especially for Rootless Podman)

To run containers correctly, especially in rootless mode, Podman requires certain cgroup controllers to be delegated to the user. By default, a non-root user might only have memory and pids controllers. For k3s (and thus k3d) to function properly, CPU, CPUSET, and I/O delegation are often needed.

!!! note "Ensure you are using cgroup v2"
    Check if `/sys/fs/cgroup/cgroup.controllers` exists. If it does, you're on cgroup v2.

If you're on cgroup v2 and using rootless Podman, you might need to enable delegation:
```bash
# Check current delegations: cat /sys/fs/cgroup/user.slice/user-$(id -u).slice/user@$(id -u).service/cgroup.controllers
# To enable necessary delegations:
sudo mkdir -p /etc/systemd/system/user@.service.d
cat <<EOF | sudo tee /etc/systemd/system/user@.service.d/delegate.conf > /dev/null
[Service]
Delegate=cpu cpuset io memory pids
EOF
sudo systemctl daemon-reload
# You might need to restart your user session or the user@.service for changes to take full effect.
# E.g., systemctl --user daemon-reload ; systemctl --user restart podman.socket
```
Reference: [Rootless Containers - cgroup2](https://rootlesscontaine.rs/getting-started/common/cgroup2/#enabling-cpu-cpuset-and-io-delegation)

!!! warning "Missing cpuset cgroup controller"
    If you experience an error regarding a missing cpuset cgroup controller, ensure the user unit `xdg-document-portal.service` is not interfering. Running `systemctl --user stop xdg-document-portal.service` might help. See [this issue](https://github.com/systemd/systemd/issues/18293#issuecomment-831397578) for context.

### 3. Podman Network for DNS

The default `podman` network might have DNS disabled. For k3d nodes to communicate correctly using DNS, create a network with DNS enabled:
```bash
podman network create k3d-network # Or any name you prefer
podman network inspect k3d-network -f '{{ .DNSEnabled }}' # Should output true
```
Then, you can specify this network in your k3d cluster configuration or via the `--network` flag:
```bash
k3d cluster create mycluster --runtime podman --network k3d-network
# Or in config:
# network: k3d-network
```

## Alternative Connection Methods (If Not Using `--runtime podman`)

If you are **not** using the `--runtime podman` flag, or if you need k3d to think it's talking to Docker (e.g., for older k3d versions or specific tooling), you can use the traditional methods of exposing the Podman socket as if it were the Docker socket:

*   **Symlinking (Rootful):**
    ```bash
    # Ensure /var/run/docker.sock does not exist or is a symlink you can overwrite
    sudo ln -s /run/podman/podman.sock /var/run/docker.sock
    # Then run k3d commands as usual
    # sudo k3d cluster create mycluster
    ```
    Some distributions offer a `podman-docker` package that might handle this.

*   **Using `DOCKER_HOST` (Rootful or Rootless):**
    This tells k3d (and other Docker client libraries) where to find the API socket.

    For rootful Podman:
    ```bash
    export DOCKER_HOST=unix:///run/podman/podman.sock
    # DOCKER_SOCK is not standard but might be used by some tools; k3d primarily respects DOCKER_HOST
    # export DOCKER_SOCK=/run/podman/podman.sock
    sudo --preserve-env=DOCKER_HOST k3d cluster create mycluster
    ```

    For rootless Podman:
    ```bash
    XDG_RUNTIME_DIR=${XDG_RUNTIME_DIR:-/run/user/$(id -u)}
    export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/podman/podman.sock
    k3d cluster create mycluster
    ```

## Platform-Specific Instructions

### macOS (Podman Machine)

1.  **Initialize and Start Podman Machine:**
    ```bash
    podman machine init
    podman machine start
    ```

2.  **Connection Details:**
    Use `podman system connection ls` to find the SSH URI for the Podman machine.
    Example output:
    ```
    Name                         URI                                                         Identity                                      Default
    podman-machine-default       ssh://core@localhost:PORT/run/user/USERID/podman/podman.sock  /Users/username/.ssh/podman-machine-default  true
    podman-machine-default-root  ssh://root@localhost:PORT/run/podman/podman.sock           /Users/username/.ssh/podman-machine-default  false
    ```
    Ensure your SSH configuration (`~/.ssh/config`) correctly points to the `IdentityFile` for `localhost` if needed.

3.  **Using k3d with Podman on macOS:**

    *   **With `--runtime podman` (Recommended):**
        K3d (v5.x.x+) should ideally pick up the correct Podman environment if `podman system connection default` points to your desired Podman machine. You might still need to ensure `DOCKER_HOST` is unset or points to the Podman socket if k3d's auto-detection isn't sufficient.
        ```bash
        # Ensure DOCKER_HOST is unset or correctly points to Podman if needed
        # unset DOCKER_HOST
        k3d cluster create mymacoscluster --runtime podman
        ```
        For rootless mode on the Podman machine, additional setup might be needed for Kubelet in user namespaces (this is advanced k3s territory):
        ```bash
        # Example for rootless, may require specific k3s versions and Podman machine setup
        # k3d cluster create mymacoscluster --runtime podman --k3s-arg '--kubelet-arg=feature-gates=KubeletInUserNamespace=true@server:*'
        ```

    *   **Manually setting `DOCKER_HOST`:**
        For rootless mode on the Podman VM:
        ```bash
        export DOCKER_HOST="ssh://core@localhost:PORT" # Replace PORT from 'podman system connection ls'
        # Optional: Point DOCKER_SOCK to the socket path inside the VM, though DOCKER_HOST is primary
        # export DOCKER_SOCK="/run/user/USERID/podman/podman.sock"
        k3d cluster create mymacoscluster
        ```
        For rootful mode on the Podman VM:
        ```bash
        export DOCKER_HOST="ssh://root@localhost:PORT" # Replace PORT
        k3d cluster create mymacoscluster
        ```

### Remote Podman Host

If Podman runs on a remote machine:
1.  Ensure Podman is correctly set up for remote access on that host.
2.  Set `DOCKER_HOST` to point to the remote Podman service, typically via SSH.
    ```bash
    export DOCKER_HOST=ssh://username@hostname[:PORT]
    # export DOCKER_SOCK="/path/to/podman.sock" # Path on the remote machine, if needed by specific tools
    k3d cluster create myremotecluster
    ```
    If using `--runtime podman`, ensure your Podman client configuration (`containers.conf` or `podman system connection default`) is set up to connect to this remote host.

## Local Registries with Podman

When creating a k3d-managed local registry for use with Podman:

*   **If k3d creates the registry alongside the cluster (e.g., `k3d cluster create --registry-create ... --runtime podman`):**
    K3d should handle connecting it to the Podman network used by the cluster.

*   **If creating the registry separately (`k3d registry create ...`):**
    Podman does not have a default "bridge" network like Docker. You might need to ensure the registry and the subsequent cluster use a shared network.
    ```bash
    # 1. Create a network if you haven't already (see "Podman Network for DNS" section)
    # podman network create k3d-network

    # 2. Create the registry, attaching it to your chosen network
    k3d registry create myregistry --default-network k3d-network

    # 3. Create your cluster, using the same network and referencing the registry
    k3d cluster create mycluster --runtime podman --network k3d-network --registry-use myregistry.localhost
    ```

!!! note "Legacy `--registry-create` Incompatibility"
    The old note about `--registry-create` assuming a "bridge" network is less relevant when `--runtime podman` is used, as k3d's Podman logic should manage networking. However, if you encounter issues, explicitly creating the registry and network as shown above is a reliable workaround.

This page will be updated as Podman support in k3d evolves. If you encounter issues, please [file an issue](https://github.com/k3d-io/k3d/issues/new?labels=bug,runtime/podman&template=bug_report.md&title=%5BBUG%5D+Podman%3A%20).
