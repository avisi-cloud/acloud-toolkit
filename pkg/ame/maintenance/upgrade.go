package maintenance

import (
	"context"
	"fmt"

	"gitlab.avisi.cloud/ame/acloud-toolkit/pkg/ame/nodedrain"
	v1 "k8s.io/api/core/v1"
)

func (c *NodeUpgraderClient) UpgradeNode(ctx context.Context, nodeName string) error {
	return c.runMaintenanceJobOnNode(ctx, "node-upgrade", nodeName, upgradeSecretSecretName, "upgrade.sh", func(node v1.Node) {
		fmt.Printf("upgrading node %s...\n", nodeName)
		fmt.Printf("node version info: kubelet: %s, container run time: %s\n", node.Status.NodeInfo.KubeletVersion, node.Status.NodeInfo.ContainerRuntimeVersion)
	}, func(node v1.Node) {
		fmt.Printf("upgraded node version info: kubelet: %s, container run time: %s\n", node.Status.NodeInfo.KubeletVersion, node.Status.NodeInfo.ContainerRuntimeVersion)
	})
}

func (c *NodeUpgraderClient) RebootNode(ctx context.Context, nodeName string) error {
	return c.runMaintenanceJobOnNodeWithScript(ctx, "node-reboot", nodeName, "test -f /var/run/reboot-required && reboot now || echo ok", func(node v1.Node) {
		nodedrain.DrainNode(ctx, []string{nodeName}, nodedrain.DrainOptions{})
	}, func(node v1.Node) {
		nodedrain.UncordonNode(ctx, []string{nodeName}, nodedrain.DrainOptions{})
	})
}

func (c *NodeUpgraderClient) RestartKubelet(ctx context.Context, nodeName string) error {
	return c.runMaintenanceJobOnNodeWithScript(ctx, "kubelet-restart", nodeName, "systemctl restart kubelet", func(node v1.Node) {}, func(node v1.Node) {})
}

func (c *NodeUpgraderClient) RestartContainerd(ctx context.Context, nodeName string) error {
	return c.runMaintenanceJobOnNodeWithScript(ctx, "containerd-restart", nodeName, "systemctl restart containerd", func(node v1.Node) {}, func(node v1.Node) {})
}
