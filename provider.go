package fleeting_plugin_k8s

import (
	"context"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"gitlab.com/gitlab-org/fleeting/fleeting/provider"
	core "k8s.io/api/core/v1"

	"fleeting_plugin_k8s/internal/kubectl"
)

var _ provider.InstanceGroup = (*InstanceGroup)(nil)

type InstanceGroup struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`

	size int

	client *kubectl.Kubectl

	log hclog.Logger

	settings provider.Settings
}

func deref[T any](v *T) T {
	if v == nil {
		var def T
		return def
	}
	return *v
}

// Init implements provider.InstanceGroup
func (g *InstanceGroup) Init(ctx context.Context, logger hclog.Logger, settings provider.Settings) (provider.ProviderInfo, error) {
	g.log = logger.With("statefulset", g.Name, "namespace", g.Namespace)
	g.settings = settings

	if g.client == nil {
		k, err := kubectl.New()

		if err != nil {
			g.log.Error("failed to create kubectl client", "error", err)
			return provider.ProviderInfo{}, err
		}

		g.client = k
	}

	if err := g.setScaleSetSize(ctx, true); err != nil {
		g.log.Error("failed to sync scale set size", "error", err)
		return provider.ProviderInfo{}, err
	}

	return provider.ProviderInfo{
		ID:        path.Join(g.Namespace, g.Name),
		MaxSize:   10,
		Version:   Version.String(),
		BuildInfo: Version.BuildInfo(),
	}, nil
}

// ConnectInfo implements provider.InstanceGroup
func (g *InstanceGroup) ConnectInfo(ctx context.Context, id string) (provider.ConnectInfo, error) {
	info := provider.ConnectInfo{
		ConnectorConfig: g.settings.ConnectorConfig,
		ID: id,
	}

	namespace, name, ok := podNameFromID(id)

	if !ok {
		g.log.Error("unable to determine namespace and name for id", "id", id)
		return info, fmt.Errorf("unable to determine namespace and name for id (%s)", id)
	}

	pod, err := g.client.Pod(ctx, namespace, name)

	if err != nil {
		g.log.Error("failed to fetch pod", "namespace", namespace, "name", name, "error", err)
		return info, fmt.Errorf("fetching pod: %w", err)
	}

	if info.OS == "" {
		info.OS = "linux"
	}

	if info.Arch == "" {
		info.Arch = "amd64"
		//if strings.HasSuffix(deref(instance.Properties.StorageProfile.ImageReference.SKU), "arm64") {
		//	info.Arch = "arm64"
		//}
	}

	if info.Username == "" {
		//info.Username = deref(instance.Properties.OSProfile.AdminUsername)
		info.Username = "root"
	}

	if info.Password == "" {
		info.Password = "blub"
		//info.Password = deref(instance.Properties.OSProfile.AdminPassword)
	}

	if info.Protocol == "" {
		info.Protocol = provider.ProtocolSSH
		//if info.OS == "windows" {
		//	info.Protocol = provider.ProtocolWinRM
		//}
	}

	info.InternalAddr = pod.Status.PodIP
	info.ExternalAddr = pod.Status.PodIP
	info.UseStaticCredentials = true

	if info.Keepalive == 0 {
		info.Keepalive = time.Second * 60
	}

	if info.Timeout == 0 {
		info.Timeout = time.Minute * 5
	}

	expires := time.Now().Add(2 * time.Minute)
	info.Expires = &expires

	//if err := g.populateNetwork(ctx, &info, instance.VirtualMachineScaleSetVM); err != nil {
	//	return provider.ConnectInfo{}, err
	//}

	g.log.Info("connect info", "id", id, "info", info)

	return info, nil
}

// Decrease implements provider.InstanceGroup
func (g *InstanceGroup) Decrease(ctx context.Context, instances []string) ([]string, error) {
	var deleted []string

	if len(instances) == 0 {
		g.log.Info("no instances to delete")
		return deleted, nil
	}

	if g.size == 0 {
		g.log.Info("scale set is already at minimum capacity")
		return deleted, nil
	}

	sort.SliceStable(instances, func(i, j int) bool {
		return instances[i] > instances[j]
	})

	delta := 0

	for _, instance := range instances {
		_, name, ok := podNameFromID(instance)

		if !ok {
			g.log.Info("unable to determine namespace and name for id", "id", instance)
			continue
		}

		if !strings.HasSuffix(name, strconv.Itoa(g.size-1)) {
			g.log.Info("instance is not the last one in the scale set, skipping deletion request", "id", instance, "size", g.size)
			continue
		}

		delta++
		deleted = append(deleted, instance)
	}

	if err := g.client.UpdateStatefulSetScale(ctx, g.Namespace, g.Name, g.size-delta); err != nil {
		g.log.Error("failed to decrease replica set capacity", "namespace", g.Namespace, "name", g.Name, "error", err)
		return []string{}, fmt.Errorf("request to decrease replica set capacity: %w", err)
	}

	g.size -= delta

	g.log.Info("decreased replica set capacity", "size", g.size, "delta", delta)

	return deleted, nil
}

func podNameFromID(id string) (string, string, bool) {
	parts := strings.Split(id, "/")

	if len(parts) == 2 {
		return parts[0], parts[1], true
	}

	return "", "", false
}

func instanceIDFromName(namespace, name string) string {
	return path.Join(namespace, name)
}

// Increase implements provider.InstanceGroup
func (g *InstanceGroup) Increase(ctx context.Context, delta int) (int, error) {
	if err := g.setScaleSetSize(ctx, false); err != nil {
		g.log.Error("failed to sync scale set size", "error", err)
		return 0, err
	}

	err := g.client.UpdateStatefulSetScale(ctx, g.Namespace, g.Name, g.size+delta)

	if err != nil {
		g.log.Error("failed to increase replica set capacity", "namespace", g.Namespace, "name", g.Name, "error", err, "delta", delta)
		return 0, fmt.Errorf("request to increase replica set capacity: %w", err)
	}

	g.size += delta

	g.log.Info("increased replica set capacity", "size", g.size, "delta", delta, "namespace", g.Namespace, "name", g.Name)

	return delta, nil
}

// Update implements provider.InstanceGroup
func (g *InstanceGroup) Update(ctx context.Context, update func(instance string, state provider.State)) error {
	pods, err := g.client.StatefulSetPods(ctx, g.Namespace, g.Name)

	if err != nil {
		g.log.Error("failed to list statefulset pods", "namespace", g.Namespace, "name", g.Name, "error", err)
		return err
	}

	for _, pod := range pods {
		state := provider.StateRunning

		if pod.Status.Phase != core.PodSucceeded {
			if pod.Status.Phase != core.PodRunning {
				state = provider.StateCreating
			}

			for _, c := range pod.Status.Conditions {
				switch c.Type {
				case core.ContainersReady, core.PodInitialized, core.PodReady, core.PodScheduled:
					if c.Status != core.ConditionTrue {
						state = provider.StateCreating
					}
				}
			}

			for _, cs := range pod.Status.ContainerStatuses {
				if !cs.Ready {
					if cs.State.Waiting != nil {
						state = provider.StateCreating
					}
				}
			}

			if pod.Status.Phase == core.PodFailed {
				state = provider.StateDeleting
			}
		}

		instanceId := instanceIDFromName(pod.Namespace, pod.Name)

		g.log.Info("updating instance state", "id", instanceId, "state", state)
		update(instanceId, state)
	}

	return nil
}

func (g *InstanceGroup) setScaleSetSize(ctx context.Context, initial bool) error {
	ss, err := g.client.StatefulSet(ctx, g.Namespace, g.Name)

	if err != nil {
		return fmt.Errorf("getting scale set size: %w", err)
	}

	capacity := ss.Spec.Replicas
	size := int(deref(capacity))

	if !initial && size != g.size {
		g.log.Error("out-of-sync capacity", "expected", g.size, "actual", size)
	}

	g.size = size
	return nil
}
