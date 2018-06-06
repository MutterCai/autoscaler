/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ali

import (
	"fmt"
	"regexp"
	"strings"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
	"k8s.io/kubernetes/pkg/scheduler/schedulercache"
)


// 文件提供给CA通用接口调用的api接口。


const (
	// ProviderName is the cloud provider name for Ali
	ProviderName = "ali"
)

// aliCloudProvider implements CloudProvider interface.
type aliCloudProvider struct {
	aliManager      *AliManager
	resourceLimiter *cloudprovider.ResourceLimiter
}

// BuildAliCloudProvider builds CloudProvider implementation for Ali.
func BuildAliCloudProvider(aliManager *AliManager, resourceLimiter *cloudprovider.ResourceLimiter) (cloudprovider.CloudProvider, error) {
	ali := &aliCloudProvider{
		aliManager:      aliManager,
		resourceLimiter: resourceLimiter,
	}
	return ali, nil
}

// Cleanup stops the go routine that is handling the current view of the ASGs in the form of a cache
func (ali *aliCloudProvider) Cleanup() error {
	ali.aliManager.Cleanup()
	return nil
}

// Name returns name of the cloud provider.
func (ali *aliCloudProvider) Name() string {
	return ProviderName
}

// NodeGroups 返回Cloud Provider配置的所有节点组。
func (ali *aliCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	asgs := ali.aliManager.getAsgs()
	ngs := make([]cloudprovider.NodeGroup, len(asgs))
	for i, asg := range asgs {
		ngs[i] = &AliNodeGroup{
			asg:        asg,
			aliManager: ali.aliManager,
		}
	}

	return ngs
}

// NodeGroupForNode 返回给定节点的节点组.
func (ali *aliCloudProvider) NodeGroupForNode(node *apiv1.Node) (cloudprovider.NodeGroup, error) {
	// TODO set ali cloud provider_id
	if len(node.Spec.ProviderID) == 0 {
		node.Spec.ProviderID = "cn-hongkong." + node.Name
	}

	ref, err := AliRefFromProviderId(node.Spec.ProviderID)
	if err != nil {
		return nil, err
	}
	asg := ali.aliManager.GetAsgForInstance(*ref)

	if asg == nil {
		return nil, nil
	}

	return &AliNodeGroup{
		asg:        asg,
		aliManager: ali.aliManager,
	}, nil
}

// Pricing 返回Cloud Provider的定价模式，如果不可用返回错误.
func (ali *aliCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	return nil, cloudprovider.ErrNotImplemented
}

// GetAvailableMachineTypes 获取可从Cloud Provider获得的所有机器类型
func (ali *aliCloudProvider) GetAvailableMachineTypes() ([]string, error) {
	return []string{}, nil
}

// NewNodeGroup builds a theoretical node group based on the node definition provided. The node group is not automatically
// created on the cloud provider side. The node group is not returned by NodeGroups() until it is created.
func (ali *aliCloudProvider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string,
	taints []apiv1.Taint, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	return nil, cloudprovider.ErrNotImplemented
}

// GetResourceLimiter returns struct containing limits (max, min) for resources (cores, memory etc.).
func (ali *aliCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	return ali.resourceLimiter, nil
}

// Refresh is called before every main loop and can be used to dynamically update cloud provider state.
// In particular the list of node groups returned by NodeGroups can change as a result of CloudProvider.Refresh().
func (ali *aliCloudProvider) Refresh() error {
	return ali.aliManager.Refresh()
}

// AliRef contains a reference to some entity in AWS world.
type AliRef struct {
	Name string
}

// AliInstanceRef contains a reference to an instance in the AWS world.
type AliInstanceRef struct {
	ProviderID string
	Name       string
}

var validAliRefIdRegex = regexp.MustCompile(`^[-0-9a-z]*\.[-0-9a-z]*$`)

// AliRefFromProviderId creates InstanceConfig object from provider id which
// must be in format: <REGION_ID>.<ECS_ID>
func AliRefFromProviderId(id string) (*AliInstanceRef, error) {
	if validAliRefIdRegex.FindStringSubmatch(id) == nil {
		return nil, fmt.Errorf("Wrong id: expected format <REGION_ID>.<ECS_ID>, got %v", id)
	}
	splitted := strings.Split(id, ".")
	return &AliInstanceRef{
		ProviderID: id,
		Name:       splitted[1],
	}, nil
}

// AliNodeGroup implements NodeGroup interface.
type AliNodeGroup struct {
	aliManager *AliManager
	asg        *asg
}

// MaxSize returns maximum size of the node group.
func (ng *AliNodeGroup) MaxSize() int {
	return ng.asg.maxSize
}

// MinSize returns minimum size of the node group.
func (ng *AliNodeGroup) MinSize() int {
	return ng.asg.minSize
}

// TargetSize returns the current TARGET size of the node group. It is possible that the
// number is different from the number of nodes registered in Kubernetes.
func (ng *AliNodeGroup) TargetSize() (int, error) {
	return ng.asg.curSize, nil
}

// Exist checks if the node group really exists on the cloud provider side. Allows to tell the
// theoretical node group from the real one.
func (ng *AliNodeGroup) Exist() bool {
	return true
}

// Create creates the node group on the cloud provider side.
func (ng *AliNodeGroup) Create() error {
	return cloudprovider.ErrAlreadyExist
}

// Autoprovisioned returns true if the node group is autoprovisioned.
func (ng *AliNodeGroup) Autoprovisioned() bool {
	return false
}

// Delete deletes the node group on the cloud provider side.
// This will be executed only for autoprovisioned node groups, once their size drops to 0.
func (ng *AliNodeGroup) Delete() error {
	return cloudprovider.ErrNotImplemented
}

// IncreaseSize increases Asg size
func (ng *AliNodeGroup) IncreaseSize(delta int) error {
	if delta <= 0 {
		return fmt.Errorf("size increase must be positive")
	}
	size := ng.asg.curSize
	if size+delta > ng.asg.maxSize {
		return fmt.Errorf("size increase too large - desired:%d max:%d", size+delta, ng.asg.maxSize)
	}
	return ng.aliManager.SetAsgSize(ng.asg, size+delta)
}

// DecreaseTargetSize decreases the target size of the node group. This function
// doesn't permit to delete any existing node and can be used only to reduce the
// request for new nodes that have not been yet fulfilled. Delta should be negative.
// It is assumed that cloud provider will not delete the existing nodes if the size
// when there is an option to just decrease the target.
func (ng *AliNodeGroup) DecreaseTargetSize(delta int) error {
	if delta >= 0 {
		return fmt.Errorf("size decrease size must be negative")
	}

	size := ng.asg.curSize
	nodes, err := ng.aliManager.GetAsgNodes(ng.asg.AliRef)
	if err != nil {
		return err
	}
	if int(size)+delta < len(nodes) {
		return fmt.Errorf("attempt to delete existing nodes targetSize:%d delta:%d existingNodes: %d",
			size, delta, len(nodes))
	}
	return ng.aliManager.SetAsgSize(ng.asg, size+delta)
}

// Belongs returns true if the given node belongs to the NodeGroup.
func (ng *AliNodeGroup) Belongs(node *apiv1.Node) (bool, error) {
	ref, err := AliRefFromProviderId(node.Spec.ProviderID)
	if err != nil {
		return false, err
	}
	targetAsg := ng.aliManager.GetAsgForInstance(*ref)
	if targetAsg == nil {
		return false, fmt.Errorf("%s doesn't belong to a known asg", node.Name)
	}
	if targetAsg.AliRef != ng.asg.AliRef {
		return false, nil
	}
	return true, nil
}

// DeleteNodes deletes the nodes from the group.
func (ng *AliNodeGroup) DeleteNodes(nodes []*apiv1.Node) error {
	size := ng.asg.curSize
	if int(size) <= ng.MinSize() {
		return fmt.Errorf("min size reached, nodes will not be deleted")
	}
	refs := make([]*AliInstanceRef, 0, len(nodes))
	for _, node := range nodes {
		belongs, err := ng.Belongs(node)
		if err != nil {
			return err
		}
		if belongs != true {
			return fmt.Errorf("%s belongs to a different asg than %s", node.Name, ng.Id())
		}
		aliref, err := AliRefFromProviderId(node.Spec.ProviderID)
		if err != nil {
			return err
		}
		refs = append(refs, aliref)
	}
	return ng.aliManager.DeleteInstances(refs)
}

// Id returns asg id.
func (ng *AliNodeGroup) Id() string {
	return ng.asg.Name
}

// Debug returns a debug string for the Asg.
func (ng *AliNodeGroup) Debug() string {
	return fmt.Sprintf("%s (%d:%d)", ng.Id(), ng.MinSize(), ng.MaxSize())
}

// Nodes returns a list of all nodes that belong to this node group.
func (ng *AliNodeGroup) Nodes() ([]string, error) {
	asgNodes, err := ng.aliManager.GetAsgNodes(ng.asg.AliRef)
	if err != nil {
		return nil, err
	}

	nodes := make([]string, len(asgNodes))

	for i, asgNode := range asgNodes {
		nodes[i] = asgNode.ProviderID
	}
	return nodes, nil
}

// TemplateNodeInfo returns a node template for this node group.
func (ng *AliNodeGroup) TemplateNodeInfo() (*schedulercache.NodeInfo, error) {
	template, err := ng.aliManager.getAsgTemplate(ng.asg)
	if err != nil {
		return nil, err
	}

	node, err := ng.aliManager.buildNodeFromTemplate(ng.asg, template)
	if err != nil {
		return nil, err
	}

	nodeInfo := schedulercache.NewNodeInfo(cloudprovider.BuildKubeProxy(ng.asg.Name))
	nodeInfo.SetNode(node)
	return nodeInfo, nil
}
