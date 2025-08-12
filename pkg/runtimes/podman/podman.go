/*
Copyright © 2023 The k3d Author(s)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package podman

import (
	"bufio"
	"context"
	"io"
	"net/netip"
	"os"
	"time"

	"github.com/k3d-io/k3d/v5/pkg/runtimes/docker"
	runtimeTypes "github.com/k3d-io/k3d/v5/pkg/runtimes/types"
	k3d "github.com/k3d-io/k3d/v5/pkg/types"
)

// Podman implements the Runtime interface for Podman.
type Podman struct {
	dockerRuntime docker.Docker
}

// ID returns the identifier of the runtime.
func (p Podman) ID() string {
	return "podman"
}

// GetHost returns the Podman daemon host.
// For now, it delegates to the Docker runtime.
func (p Podman) GetHost() string {
	return p.dockerRuntime.GetHost()
}

// CreateNode creates a new containerized node.
// For now, it delegates to the Docker runtime.
func (p Podman) CreateNode(ctx context.Context, node *k3d.Node) error {
	return p.dockerRuntime.CreateNode(ctx, node)
}

// DeleteNode deletes a containerized node.
// For now, it delegates to the Docker runtime.
func (p Podman) DeleteNode(ctx context.Context, node *k3d.Node) error {
	return p.dockerRuntime.DeleteNode(ctx, node)
}

// RenameNode renames a containerized node.
// For now, it delegates to the Docker runtime.
func (p Podman) RenameNode(ctx context.Context, node *k3d.Node, newHostname string) error {
	return p.dockerRuntime.RenameNode(ctx, node, newHostname)
}

// GetNodesByLabel returns a list of nodes matching the given labels.
// For now, it delegates to the Docker runtime.
func (p Podman) GetNodesByLabel(ctx context.Context, labels map[string]string) ([]*k3d.Node, error) {
	return p.dockerRuntime.GetNodesByLabel(ctx, labels)
}

// GetNode returns a node matching the given node object.
// For now, it delegates to the Docker runtime.
func (p Podman) GetNode(ctx context.Context, node *k3d.Node) (*k3d.Node, error) {
	return p.dockerRuntime.GetNode(ctx, node)
}

// GetNodeStatus returns the status of a node.
// For now, it delegates to the Docker runtime.
func (p Podman) GetNodeStatus(ctx context.Context, node *k3d.Node) (bool, string, error) {
	return p.dockerRuntime.GetNodeStatus(ctx, node)
}

// GetNodesInNetwork returns a list of nodes connected to the given network.
// For now, it delegates to the Docker runtime.
func (p Podman) GetNodesInNetwork(ctx context.Context, networkName string) ([]*k3d.Node, error) {
	return p.dockerRuntime.GetNodesInNetwork(ctx, networkName)
}

// CreateNetworkIfNotPresent creates a new network if it doesn't exist.
// For now, it delegates to the Docker runtime.
func (p Podman) CreateNetworkIfNotPresent(ctx context.Context, network *k3d.ClusterNetwork) (*k3d.ClusterNetwork, bool, error) {
	return p.dockerRuntime.CreateNetworkIfNotPresent(ctx, network)
}

// GetKubeconfig returns the Kubeconfig for a given node.
// For now, it delegates to the Docker runtime.
func (p Podman) GetKubeconfig(ctx context.Context, node *k3d.Node) (io.ReadCloser, error) {
	return p.dockerRuntime.GetKubeconfig(ctx, node)
}

// DeleteNetwork deletes a network.
// For now, it delegates to the Docker runtime.
func (p Podman) DeleteNetwork(ctx context.Context, networkName string) error {
	return p.dockerRuntime.DeleteNetwork(ctx, networkName)
}

// StartNode starts an existing container.
// For now, it delegates to the Docker runtime.
func (p Podman) StartNode(ctx context.Context, node *k3d.Node) error {
	return p.dockerRuntime.StartNode(ctx, node)
}

// StopNode stops an existing container.
// For now, it delegates to the Docker runtime.
func (p Podman) StopNode(ctx context.Context, node *k3d.Node) error {
	return p.dockerRuntime.StopNode(ctx, node)
}

// CreateVolume creates a new volume.
// For now, it delegates to the Docker runtime.
func (p Podman) CreateVolume(ctx context.Context, volumeName string, labels map[string]string) error {
	return p.dockerRuntime.CreateVolume(ctx, volumeName, labels)
}

// DeleteVolume deletes a volume.
// For now, it delegates to the Docker runtime.
func (p Podman) DeleteVolume(ctx context.Context, volumeName string) error {
	return p.dockerRuntime.DeleteVolume(ctx, volumeName)
}

// GetVolume returns information about a volume.
// For now, it delegates to the Docker runtime.
func (p Podman) GetVolume(volumeName string) (string, error) {
	return p.dockerRuntime.GetVolume(volumeName)
}

// GetVolumesByLabel returns a list of volumes matching the given labels.
// For now, it delegates to the Docker runtime.
func (p Podman) GetVolumesByLabel(ctx context.Context, labels map[string]string) ([]string, error) {
	return p.dockerRuntime.GetVolumesByLabel(ctx, labels)
}

// GetImageStream downloads an image.
// For now, it delegates to the Docker runtime.
func (p Podman) GetImageStream(ctx context.Context, images []string) (io.ReadCloser, error) {
	return p.dockerRuntime.GetImageStream(ctx, images)
}

// GetRuntimePath returns the path of the Podman socket.
// For now, it delegates to the Docker runtime.
func (p Podman) GetRuntimePath() string {
	return p.dockerRuntime.GetRuntimePath()
}

// ExecInNode executes a command in a node.
// For now, it delegates to the Docker runtime.
func (p Podman) ExecInNode(ctx context.Context, node *k3d.Node, cmd []string) error {
	return p.dockerRuntime.ExecInNode(ctx, node, cmd)
}

// ExecInNodeWithStdin executes a command in a node with stdin.
// For now, it delegates to the Docker runtime.
func (p Podman) ExecInNodeWithStdin(ctx context.Context, node *k3d.Node, cmd []string, stdin io.ReadCloser) error {
	return p.dockerRuntime.ExecInNodeWithStdin(ctx, node, cmd, stdin)
}

// ExecInNodeGetLogs executes a command in a node and returns the logs.
// For now, it delegates to the Docker runtime.
func (p Podman) ExecInNodeGetLogs(ctx context.Context, node *k3d.Node, cmd []string) (*bufio.Reader, error) {
	return p.dockerRuntime.ExecInNodeGetLogs(ctx, node, cmd)
}

// GetNodeLogs returns the logs of a node.
// For now, it delegates to the Docker runtime.
func (p Podman) GetNodeLogs(ctx context.Context, node *k3d.Node, since time.Time, opts *runtimeTypes.NodeLogsOpts) (io.ReadCloser, error) {
	return p.dockerRuntime.GetNodeLogs(ctx, node, since, opts)
}

// GetImages returns a list of images.
// For now, it delegates to the Docker runtime.
func (p Podman) GetImages(ctx context.Context) ([]string, error) {
	return p.dockerRuntime.GetImages(ctx)
}

// CopyToNode copies a file to a node.
// For now, it delegates to the Docker runtime.
func (p Podman) CopyToNode(ctx context.Context, source string, destination string, node *k3d.Node) error {
	return p.dockerRuntime.CopyToNode(ctx, source, destination, node)
}

// WriteToNode writes data to a file on a node.
// For now, it delegates to the Docker runtime.
func (p Podman) WriteToNode(ctx context.Context, content []byte, destination string, filemode os.FileMode, node *k3d.Node) error {
	return p.dockerRuntime.WriteToNode(ctx, content, destination, filemode, node)
}

// ReadFromNode reads a file from a node.
// For now, it delegates to the Docker runtime.
func (p Podman) ReadFromNode(ctx context.Context, filepath string, node *k3d.Node) (io.ReadCloser, error) {
	return p.dockerRuntime.ReadFromNode(ctx, filepath, node)
}

// GetHostIP returns the IP address of the host.
// For now, it delegates to the Docker runtime.
func (p Podman) GetHostIP(ctx context.Context, networkName string) (netip.Addr, error) {
	return p.dockerRuntime.GetHostIP(ctx, networkName)
}

// ConnectNodeToNetwork connects a node to a network.
// For now, it delegates to the Docker runtime.
func (p Podman) ConnectNodeToNetwork(ctx context.Context, node *k3d.Node, networkName string) error {
	return p.dockerRuntime.ConnectNodeToNetwork(ctx, node, networkName)
}

// DisconnectNodeFromNetwork disconnects a node from a network.
// For now, it delegates to the Docker runtime.
func (p Podman) DisconnectNodeFromNetwork(ctx context.Context, node *k3d.Node, networkName string) error {
	return p.dockerRuntime.DisconnectNodeFromNetwork(ctx, node, networkName)
}

// Info returns runtime information.
// For now, it delegates to the Docker runtime.
func (p Podman) Info() (*runtimeTypes.RuntimeInfo, error) {
	return p.dockerRuntime.Info()
}

// GetNetwork returns network information.
// For now, it delegates to the Docker runtime.
func (p Podman) GetNetwork(ctx context.Context, network *k3d.ClusterNetwork) (*k3d.ClusterNetwork, error) {
	return p.dockerRuntime.GetNetwork(ctx, network)
}
