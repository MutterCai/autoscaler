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
	"reflect"
	"sync"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/config/dynamic"
	"github.com/golang/glog"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
)

const scaleToZeroSupported = true

type asgCache struct {
	// 扩展组群
	registeredAsgs []*asg
	// 通过扩展组关联实例
	asgToInstances map[AliRef][]AliInstanceRef
	// 通过实例关联扩展组
	instanceToAsg  map[AliInstanceRef]*asg
	mutex          sync.Mutex
	service        autoScalingWrapper
	interrupt      chan struct{}

	asgAutoDiscoverySpecs []cloudprovider.ASGAutoDiscoveryConfig
	// 是否明确配置，就是当前是否有保存最大/最小值用于控制，而不是使用Api的最大最小值控制。
	explicitlyConfigured  map[AliRef]bool
}

type asg struct {
	AliRef

	minSize int
	maxSize int
	curSize int
	// 信息集
	ScalingGroupItem		ess.ScalingGroup
	// 配置信息
	ScalingConfiguration	[]ess.ScalingConfiguration
	activity 				*inScaling
	//AvailabilityZones       []string
	//LaunchConfigurationName string
	//Tags                    []*autoscaling.TagDescription
}

func newASGCache(service autoScalingWrapper, explicitSpecs []string, autoDiscoverySpecs []cloudprovider.ASGAutoDiscoveryConfig) (*asgCache, error) {
	registry := &asgCache{
		registeredAsgs:        make([]*asg, 0),
		service:               service,
		asgToInstances:        make(map[AliRef][]AliInstanceRef),
		instanceToAsg:         make(map[AliInstanceRef]*asg),
		interrupt:             make(chan struct{}),
		asgAutoDiscoverySpecs: autoDiscoverySpecs,
		explicitlyConfigured:  make(map[AliRef]bool),
	}
	// 将对应的explicitlyConfigured设置成true
	if err := registry.parseExplicitAsgs(explicitSpecs); err != nil {
		return nil, err
	}

	return registry, nil
}

// Fetch explicitly configured ASGs. These ASGs should never be unregistered
// during refreshes, even if they no longer exist in AWS.
func (m *asgCache) parseExplicitAsgs(specs []string) error {
	for _, spec := range specs {
		asg, err := m.buildAsgFromSpec(spec)
		if err != nil {
			return fmt.Errorf("failed to parse node group spec: %v", err)
		}
		m.explicitlyConfigured[asg.AliRef] = true
		m.register(asg)
	}

	return nil
}

// 注册ASG。如果ASG已注册，则返回true。
func (m *asgCache) register(asg *asg) bool {
	// 遍历已注册的ASGs
	for i := range m.registeredAsgs {
		// 当asg的名称存在于已注册的ASGs的对象中，检测是否深度相等，不相等则进行更新
		if existing := m.registeredAsgs[i]; existing.AliRef == asg.AliRef {
			// 检查是否深度相等
			if reflect.DeepEqual(existing, asg) {
				return false
			}

			glog.V(4).Infof("Updating ASG %s", asg.AliRef.Name)

			// 显式注册组应始终使用手动提供的最小/最大值，而不是API返回的值
			// 如果不明确配置，则使用api的返回值。
			if !m.explicitlyConfigured[asg.AliRef] {
				existing.minSize = asg.minSize
				existing.maxSize = asg.maxSize
			}

			existing.curSize = asg.curSize

			// Those information are mainly required to create templates when scaling
			// from zero
			//existing.AvailabilityZones = asg.AvailabilityZones
			//existing.LaunchConfigurationName = asg.LaunchConfigurationName
			//existing.Tags = asg.Tags
			existing.ScalingGroupItem = asg.ScalingGroupItem
			existing.activity = asg.activity
			return true
		}
	}
	glog.V(1).Infof("Registering ASG %s", asg.AliRef.Name)
	// 如果不存在已注册的ASG，进行注册添加
	m.registeredAsgs = append(m.registeredAsgs, asg)
	return true
}

// 取消注册ASG。如果ASG取消注册，则返回true。
func (m *asgCache) unregister(a *asg) bool {
	updated := make([]*asg, 0, len(m.registeredAsgs))
	changed := false
	// 遍历ASGs注册表，去除asg的信息
	for _, existing := range m.registeredAsgs {
		if existing.AliRef == a.AliRef {
			glog.V(1).Infof("Unregistered ASG %s", a.AliRef.Name)
			changed = true
			continue
		}
		updated = append(updated, existing)
	}
	m.registeredAsgs = updated
	return changed
}

func (m *asgCache) buildAsgFromSpec(spec string) (*asg, error) {
	s, err := dynamic.SpecFromString(spec, scaleToZeroSupported)
	if err != nil {
		return nil, fmt.Errorf("failed to parse node group spec: %v", err)
	}
	asg := &asg{
		AliRef:  AliRef{Name: s.Name},
		minSize: s.MinSize,
		maxSize: s.MaxSize,
	}
	return asg, nil
}

// Get returns the currently registered ASGs
func (m *asgCache) Get() []*asg {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.registeredAsgs
}

// FindForInstance returns AsgConfig of the given Instance
func (m *asgCache) FindForInstance(instance AliInstanceRef) *asg {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if asg, found := m.instanceToAsg[instance]; found {
		return asg
	}

	return nil
}

// InstancesByAsg returns the nodes of an ASG
func (m *asgCache) InstancesByAsg(ref AliRef) ([]AliInstanceRef, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if instances, found := m.asgToInstances[ref]; found {
		return instances, nil
	}

	return nil, fmt.Errorf("Error while looking for instances of ASG: %s", ref)
}

// Fetch automatically discovered ASGs. These ASGs should be unregistered if
// they no longer exist in AWS.
func (m *asgCache) fetchAutoAsgNames() ([]string, error) {
	groupNames := make([]string, 0)

	for _, spec := range m.asgAutoDiscoverySpecs {
		for name := range spec.Tags{
			groupNames = append(groupNames, name)
		}
	}
	return groupNames, nil
}

func (m *asgCache) buildAsgNames() ([]string, error) {
	// Collect explicitly specified names
	refreshNames := make([]string, len(m.explicitlyConfigured))
	i := 0
	// 遍历明确配置的对象，并获取对应的name
	for k := range m.explicitlyConfigured {
		refreshNames[i] = k.Name
		i++
	}

	// Append auto-discovered names
	autoDiscoveredNames, err := m.fetchAutoAsgNames()
	if err != nil {
		return nil, err
	}
	for _, name := range autoDiscoveredNames {
		autoRef := AliRef{Name: name}

		if m.explicitlyConfigured[autoRef] {
			// This ASG was already explicitly configured, we only need to fetch it once
			continue
		}

		refreshNames = append(refreshNames, name)
	}

	return refreshNames, nil
}

// 反射提取struct的字段并赋值
func SetInfo(o interface{}, fieldName string, setValue string) error {
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Ptr && !v.Elem().CanSet() { //判断是否为指针类型 元素是否可以修改
		return fmt.Errorf("无法设置结构体")
	} else {
		v = v.Elem() //实际取得的对象
	}
	//判断字段是否存在
	f := v.FieldByName(fieldName)
	if !f.IsValid() {
		return fmt.Errorf("无效的字段名称")
	}
	//设置字段
	if f := v.FieldByName(fieldName); f.Kind() == reflect.String {
		f.SetString(setValue)
	}
	return nil
}

// 通过自动发现ASG和显示配置，重新生成缓存
// 这里改成，通过地域ID信息，重新生成缓存
func (m *asgCache) regenerate() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	newInstanceToAsgCache := make(map[AliInstanceRef]*asg)
	newAsgToInstancesCache := make(map[AliRef][]AliInstanceRef)
	newScalingConfiguration := make(map[AliRef][]ess.ScalingConfiguration)

	// 从提供的参数，获取对应的组名
	refreshNames, err := m.buildAsgNames()
	if err != nil {
		return err
	}
	// Fetch details of all ASGs
	glog.V(4).Infof("Regenerating instance to ASG map for ASGs: %v", refreshNames)

	// 通过组名获取ASGs的信息
	groups, err := m.service.getAutoscalingGroupsByNames(refreshNames)
	if err != nil {
		return err
	}

	// 注册更新ASGs信息
	exists := make(map[AliRef]bool)
	for _, group := range groups {
		// 通过获取的asg信息分别构建asg对象
		asg, err := m.buildAsgFromAli(group)
		if err != nil {
			return err
		}
		// 执行状态
		asg.activity.PendingCapacity = group.PendingCapacity
		asg.activity.RemovingCapacity = group.RemovingCapacity

		// 标注API列表存在该ASG
		exists[asg.AliRef] = true
		// 注册asg(如果asg已经注册过，则会进行数据更新)
		m.register(asg)

		// 获取ASG的配置信息
		configuration, err := m.service.getAutoscalingGroupConfigurationByGroupID(group.ScalingGroupId)
		newScalingConfiguration[asg.AliRef] = configuration

		// 获取asg的instances
		instances, err := m.service.getAutoscalingGroupInstancesByGroupID(group.ScalingGroupId)
		if err != nil {
			return err
		}

		newAsgToInstancesCache[asg.AliRef] = make([]AliInstanceRef, len(instances))
		// 更新下instances和asg的缓存关系
		for i, instance := range instances {
			// 创建关联对象
			ref := m.buildInstanceRefFromAli(&instance, asg.ScalingGroupItem.RegionId)
			// 创建instance关联扩展组
			newInstanceToAsgCache[ref] = asg
			// 创建asg关联实例
			newAsgToInstancesCache[asg.AliRef][i] = ref
		}
	}

	// 取消注册长期不存在的(无法自动发现的)ASG
	for _, asg := range m.registeredAsgs {
		// api获取的数据里没有registeredAsgs里的某个ASG信息 并且 没有明确配置
		if !exists[asg.AliRef] && !m.explicitlyConfigured[asg.AliRef] {
			// 取消注册asg(就是将ASGs注册表里的和asg名称相同的删除掉)
			m.unregister(asg)
		}else if val, ok := newScalingConfiguration[asg.AliRef]; ok{
			// asg获取配置信息
			asg.ScalingConfiguration = val
		}
	}
	// 对应更新下instance和asg的缓存关系
	m.asgToInstances = newAsgToInstancesCache
	m.instanceToAsg = newInstanceToAsgCache
	return nil
}

func (m *asgCache) buildAsgFromAli(g ess.ScalingGroup) (*asg, error) {
	spec := dynamic.NodeGroupSpec{
		Name:               g.ScalingGroupName,
		MinSize:            g.MinSize,
		MaxSize:            g.MaxSize,
		// 是否支持从0创建
		SupportScaleToZero: scaleToZeroSupported,
	}
	if verr := spec.Validate(); verr != nil {
		return nil, fmt.Errorf("failed to create node group spec: %v", verr)
	}
	activity := &inScaling{
		PendingCapacity: 0,
		RemovingCapacity: 0,
	}
	asg := &asg{
		AliRef:  AliRef{Name: spec.Name},
		minSize: spec.MinSize,
		maxSize: spec.MaxSize,
		// TODO curSize所需的容量，这里暂时用正常运行的容量
		curSize:                 g.ActiveCapacity,
		//AvailabilityZones:       g.AvailabilityZones,
		//LaunchConfigurationName: g.LaunchConfigurationName,
		//Tags: g.Tags,
		// 增加扩展组的完整信息对象
		ScalingGroupItem:	g,
		activity: activity,
	}
	return asg, nil
}

func (m *asgCache) buildInstanceRefFromAli(instance *ess.ScalingInstance, regionID string) AliInstanceRef {
	// This is the format of "REGION.NODEID"
	if instance.InstanceId[1] == '-'{
		instance.InstanceId = instance.InstanceId[:1]+"z"+instance.InstanceId[2:]+"z"
	}
	providerID := fmt.Sprintf("%s.%s", regionID, instance.InstanceId)
	return AliInstanceRef{
		ProviderID: providerID,
		Name:       instance.InstanceId,
	}
}

// Cleanup closes the channel to signal the go routine to stop that is handling the cache
func (m *asgCache) Cleanup() {
	close(m.interrupt)
}
