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

//go:generate go run ec2_instance_types/gen.go

package ali

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
	"github.com/golang/glog"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/utils/gpu"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
	//provider_ali "github.com/AliyunContainerService/alicloud-controller-manager/alicloud"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"encoding/json"
)

const (
	operationWaitTimeout    = 5 * time.Second
	operationPollInterval   = 100 * time.Millisecond
	maxRecordsReturnedByAPI = 100
	refreshInterval         = 10 * time.Second
)


// AliManager is handles ali communication and data caching.
type AliManager struct {
	service     autoScalingWrapper
	asgCache    *asgCache
	lastRefresh time.Time
}

type asgTemplate struct {
	InstanceType 	*instanceType
	Region       	string
	Configuration	ess.ScalingConfiguration
}

// createAwsManagerInternal allows for a customer autoScalingWrapper to be passed in by tests
func createAliManagerInternal(
	configReader io.Reader,
	discoveryOpts cloudprovider.NodeGroupDiscoveryOptions,
	service *autoScalingWrapper,
) (*AliManager, error) {
	var cfg CloudConfig
	//if configReader != nil {
	//	if err := gcfg.ReadInto(&cfg, configReader); err != nil {
	//		glog.Errorf("Couldn't read config: %v", err)
	//		return nil, err
	//	}
	//}
	cfg.Global.Region = "cn-hongkong"
	cfg.Global.AccessKeyID = "LTAITh6uSXRCv14k"
	cfg.Global.AccessKeySecret = "zm4qbsuKnDQ3VRIxAQzFKvMt9aOBuI"
	if service == nil {
		// ess的API实例
		essClient, err := ess.NewClientWithAccessKey(
			cfg.Global.Region,
			cfg.Global.AccessKeyID,
			cfg.Global.AccessKeySecret)
		if err != nil {
			return nil, err
		}
		// ecs的API实例
		ecsClient, err := ecs.NewClientWithAccessKey(
			cfg.Global.Region,
			cfg.Global.AccessKeyID,
			cfg.Global.AccessKeySecret)
		if err != nil {
			return nil, err
		}

		service = &autoScalingWrapper{
			ess: *essClient,
			ecs: *ecsClient,
			cfg: cfg,
		}
	}

	specs, err := discoveryOpts.ParseASGAutoDiscoverySpecs()
	if err != nil {
		return nil, err
	}

	cache, err := newASGCache(*service, discoveryOpts.NodeGroupSpecs, specs)
	if err != nil {
		return nil, err
	}

	manager := &AliManager{
		service:  *service,
		asgCache: cache,
	}

	if err := manager.forceRefresh(); err != nil {
		return nil, err
	}

	return manager, nil
}

// CreateAwsManager constructs awsManager object.
func CreateAliManager(configReader io.Reader, discoveryOpts cloudprovider.NodeGroupDiscoveryOptions) (*AliManager, error) {
	return createAliManagerInternal(configReader, discoveryOpts, nil)
}

// Refresh is called before every main loop and can be used to dynamically update cloud provider state.
// In particular the list of node groups returned by NodeGroups can change as a result of CloudProvider.Refresh().
func (m *AliManager) Refresh() error {
	if m.lastRefresh.Add(refreshInterval).After(time.Now()) {
		return nil
	}
	return m.forceRefresh()
}

func (m *AliManager) forceRefresh() error {
	if err := m.asgCache.regenerate(); err != nil {
		glog.Errorf("Failed to regenerate ASG cache: %v", err)
		return err
	}
	m.lastRefresh = time.Now()
	glog.V(2).Infof("Refreshed ASG list, next refresh after %v", m.lastRefresh.Add(refreshInterval))
	return nil
}

// GetAsgForInstance returns AsgConfig of the given Instance
func (m *AliManager) GetAsgForInstance(instance AliInstanceRef) *asg {
	return m.asgCache.FindForInstance(instance)
}

// Cleanup the ASG cache.
func (m *AliManager) Cleanup() {
	m.asgCache.Cleanup()
}

func (m *AliManager) getAsgs() []*asg {
	return m.asgCache.Get()
}

func (m *AliManager) SetAsgSize(asg *asg, size int) error {
	// TODO 在循环开始应该将所有的规则清空最好。
	// 创建规则、应用、删除规则
	ruleResp, err := m.service.createAutoscalingRule(asg.ScalingGroupItem.ScalingGroupId, size)
	if err != nil {
		return err
	}
	glog.V(0).Infof("伸缩规则应用：Setting asg %s size to %d", asg.Name, size)
	// scalingActivityId 暂时用不上
	_, err = m.service.executeAutoscalingRule(ruleResp.ScalingRuleAri)
	if err != nil {
		return err
	}
	m.service.deleteAutoscalingRule(ruleResp.ScalingRuleId)
	return nil
}

// DeleteInstances deletes the given instances. All instances must be controlled by the same ASG.
func (m *AliManager) DeleteInstances(instances []*AliInstanceRef) error {
	if len(instances) == 0 {
		return nil
	}
	commonAsg := m.asgCache.FindForInstance(*instances[0])
	if commonAsg == nil {
		return fmt.Errorf("can't delete instance %s, which is not part of an ASG", instances[0].Name)
	}

	instancesList := map[string][]string{}

	for _, instance := range instances {
		asg := m.asgCache.FindForInstance(*instance)

		if asg != commonAsg {
			instanceIds := make([]string, len(instances))
			for i, instance := range instances {
				instanceIds[i] = instance.Name
			}

			return fmt.Errorf("can't delete instances %s as they belong to at least two different ASGs (%s and %s)", strings.Join(instanceIds, ","), commonAsg.Name, asg.Name)
		}
		// 以autoscalingGroupID划分组
		instancesList[asg.ScalingGroupItem.ScalingGroupId] = append(instancesList[asg.ScalingGroupItem.ScalingGroupId], instance.Name)
	}
	jdata, _ := json.Marshal(instancesList)
	glog.Infof("instancesList: %s", jdata)
	// 删掉实例，并且缩容
	for autoscalingGroupID, instances := range instancesList {
		scalingActivityId, err := m.service.removeInstances(autoscalingGroupID, instances)
		if err != nil {
			return err
		}
		glog.V(4).Infof("删除请求完毕。ASGID：%s，instances：%v。scalingActivityId：%s", autoscalingGroupID, instances,scalingActivityId)
	}

	return nil
}

// GetAsgNodes returns Asg nodes.
func (m *AliManager) GetAsgNodes(ref AliRef) ([]AliInstanceRef, error) {
	return m.asgCache.InstancesByAsg(ref)
}

// 获取ASG的模板信息(阿里云下应该是伸缩配置)
func (m *AliManager) getAsgTemplate(asg *asg) (*asgTemplate, error) {

	cfg, err := m.service.getAutoscalingGroupConfigurationByGroupID(asg.ScalingGroupItem.ScalingGroupId)
	if err != nil {
		return nil, err
	}
	// TODO 这里获取的配置是一个数组，暂时以第一个数据作为模板
	return &asgTemplate{
		InstanceType: InstanceTypes[cfg[0].InstanceType],
		Configuration: cfg[0],
		Region: asg.ScalingGroupItem.RegionId,
	}, nil
}

func (m *AliManager) buildNodeFromTemplate(asg *asg, template *asgTemplate) (*apiv1.Node, error) {
	node := apiv1.Node{}
	nodeName := fmt.Sprintf("%s-asg-%d", asg.Name, rand.Int63())

	node.ObjectMeta = metav1.ObjectMeta{
		Name:     nodeName,
		SelfLink: fmt.Sprintf("/api/v1/nodes/%s", nodeName),
		Labels:   map[string]string{},
	}

	node.Status = apiv1.NodeStatus{
		Capacity: apiv1.ResourceList{},
	}

	// TODO: get a real value.
	node.Status.Capacity[apiv1.ResourcePods] = *resource.NewQuantity(110, resource.DecimalSI)
	node.Status.Capacity[apiv1.ResourceCPU] = *resource.NewQuantity(template.InstanceType.VCPU, resource.DecimalSI)
	node.Status.Capacity[gpu.ResourceNvidiaGPU] = *resource.NewQuantity(template.InstanceType.GPU, resource.DecimalSI)
	node.Status.Capacity[apiv1.ResourceMemory] = *resource.NewQuantity(template.InstanceType.MemoryMb*1024*1024, resource.DecimalSI)

	// TODO: use proper allocatable!!
	node.Status.Allocatable = node.Status.Capacity

	// NodeLabels
	//node.Labels = cloudprovider.JoinStringMaps(node.Labels, extractLabelsFromAsg(template.Tags))
	// GenericLabels
	node.Labels = cloudprovider.JoinStringMaps(node.Labels, buildGenericLabels(template, nodeName))

	//node.Spec.Taints = extractTaintsFromAsg(template.Tags)

	node.Status.Conditions = cloudprovider.BuildReadyConditions()
	return &node, nil
}

func buildGenericLabels(template *asgTemplate, nodeName string) map[string]string {
	result := make(map[string]string)
	// TODO: extract it somehow
	result[kubeletapis.LabelArch] = cloudprovider.DefaultArch
	result[kubeletapis.LabelOS] = cloudprovider.DefaultOS

	result[kubeletapis.LabelInstanceType] = template.InstanceType.InstanceType

	result[kubeletapis.LabelZoneRegion] = template.Region
	result[kubeletapis.LabelZoneFailureDomain] = template.Region
	result[kubeletapis.LabelHostname] = nodeName
	return result
}

func extractLabelsFromAsg(tags []*autoscaling.TagDescription) map[string]string {
	result := make(map[string]string)

	for _, tag := range tags {
		k := *tag.Key
		v := *tag.Value
		splits := strings.Split(k, "k8s.io/cluster-autoscaler/node-template/label/")
		if len(splits) > 1 {
			label := splits[1]
			if label != "" {
				result[label] = v
			}
		}
	}

	return result
}

func extractTaintsFromAsg(tags []*autoscaling.TagDescription) []apiv1.Taint {
	taints := make([]apiv1.Taint, 0)

	for _, tag := range tags {
		k := *tag.Key
		v := *tag.Value
		splits := strings.Split(k, "k8s.io/cluster-autoscaler/node-template/taint/")
		if len(splits) > 1 {
			values := strings.SplitN(v, ":", 2)
			taints = append(taints, apiv1.Taint{
				Key:    splits[1],
				Value:  values[0],
				Effect: apiv1.TaintEffect(values[1]),
			})
		}
	}
	return taints
}
