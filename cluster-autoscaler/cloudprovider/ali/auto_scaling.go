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
	//"fmt"
	//"github.com/golang/glog"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	//provider_ali "github.com/AliyunContainerService/alicloud-controller-manager/alicloud"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"fmt"
	"strconv"
)

const (
	// 伸缩规则的调整方式。可选值：
	//- QuantityChangeInCapacity：增加或减少指定数量的ECS实例。
	//- PercentChangeInCapacity： 增加或减少指定比例的ECS实例。
	//- TotalCapacity： 将当前伸缩组的ECS实例数量调整到指定数量。
	AdjustmentType = "TotalCapacity"
)

// CloudConfig wraps the settings for the AWS cloud provider.
type CloudConfig struct {
	Global struct {
		KubernetesClusterTag string

		AccessKeyID     string `json:"accessKeyID"`
		AccessKeySecret string `json:"accessKeySecret"`
		Region          string `json:"region"`
	}
}

// autoScalingWrapper 根据sdk封装一层接口，负责与阿里云的通讯
type autoScalingWrapper struct {
	ess ess.Client
	ecs ecs.Client
	cfg CloudConfig
}

// 获取该地域的所有ASG信息
func (m *autoScalingWrapper) getAutoscalingGroupsByNames(names []string) ([]ess.ScalingGroup, error) {
	// 构造查询伸缩组的请求
	request := ess.CreateDescribeScalingGroupsRequest()
	for i := 0; i < len(names); i++ {
		//if i == 0 {
		//	//request.ScalingGroupName = names[i]
		//	err := SetInfo(request, "ScalingGroupName", names[i])
		//	if err != nil {
		//		return nil, err
		//	}
		//}else{
		// TODO 这里阿里云sdk和官方的文档都有问题，ScalingGroupName字段无法使用，只能从ScalingGroupName1开始
		err := SetInfo(request, "ScalingGroupName"+strconv.Itoa(i+1), names[i])
		if err != nil {
			return nil, err
		}
		//}
	}
	// 执行查询伸缩组
	response, err := m.ess.DescribeScalingGroups(request)
	if err != nil {
		return nil, err
	}
	return response.ScalingGroups.ScalingGroup, nil
}

// 获取伸缩规则
//func (m *autoScalingWrapper) getAutoscalingRules(scalingGroupId string) ([]ess.ScalingRule, error) {
//	reqDSR := ess.CreateDescribeScalingRulesRequest()
//	reqDSR.RegionId = m.cfg.Global.Region
//	reqDSR.ScalingGroupId = scalingGroupId
//	respDSR, err := m.ess.DescribeScalingRules(reqDSR)
//	if err != nil {
//		return nil, err
//	}
//	return respDSR.ScalingRules.ScalingRule, nil
//}

// 创建伸缩规则
func (m *autoScalingWrapper) createAutoscalingRule(scalingGroupID string, size int) (*ess.CreateScalingRuleResponse, error) {
	request := ess.CreateCreateScalingRuleRequest()
	request.ScalingGroupId = scalingGroupID
	request.AdjustmentType = AdjustmentType
	request.AdjustmentValue = requests.NewInteger(size)
	//reqCSR.ScalingRuleName = "tempAutoscalingRule"
	response, err := m.ess.CreateScalingRule(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 修改无法获取ari信息，暂时不用
// 修改伸缩规则(修改后需要执行应用规则才生效)
//func (m *autoScalingWrapper) modifyAutoscalingRule(scalingRuleId string, size int) (string, error) {
//	reqMSR := ess.CreateModifyScalingRuleRequest()
//	reqMSR.RegionId = m.cfg.Global.Region
//	reqMSR.ScalingRuleId = scalingRuleId
//	reqMSR.AdjustmentType = AdjustmentType
//	reqMSR.AdjustmentValue = requests.Integer(size)
//	respMSR, err := m.ess.ModifyScalingRule(reqMSR)
//	if err != nil {
//		return "", err
//	}
//	return respMSR.RequestId, nil
//}

// 删除规则
func (m *autoScalingWrapper) deleteAutoscalingRule(scalingRuleID string) (bool, error) {
	if len(scalingRuleID) == 0 {
		return false, fmt.Errorf("规则id为空")
	}
	request := ess.CreateDeleteScalingRuleRequest()
	request.ScalingRuleId = scalingRuleID
	_, err := m.ess.DeleteScalingRule(request)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 应用规则
func (m *autoScalingWrapper) executeAutoscalingRule(scalingRuleAri string) (string, error) {
	if len(scalingRuleAri) == 0 {
		return "", fmt.Errorf("规则Ari为空")
	}
	request := ess.CreateExecuteScalingRuleRequest()
	request.ScalingRuleAri = scalingRuleAri
	response, err := m.ess.ExecuteScalingRule(request)
	if err != nil {
		return "", err
	}
	// 返回伸缩活动的id
	return response.ScalingActivityId, nil
}

// 重新构造方法：通过名称获取扩展组实例数据
func (m *autoScalingWrapper) getAutoscalingGroupInstancesByGroupID(scalingGroupID string) ([]ess.ScalingInstance, error) {
	if len(scalingGroupID) == 0 {
		return nil, nil
	}
	// 构造查询实例的请求
	request := ess.CreateDescribeScalingInstancesRequest()
	// 使用ASG的ID进行索引
	request.ScalingGroupId = scalingGroupID
	response, err := m.ess.DescribeScalingInstances(request)
	if err != nil {
		return nil, err
	}
	return response.ScalingInstances.ScalingInstance, nil
}

func (m *autoScalingWrapper) getAutoscalingGroupConfigurationByGroupID(scalingGroupID string) ([]ess.ScalingConfiguration, error) {
	if len(scalingGroupID) == 0 {
		return nil, fmt.Errorf("伸缩组id为空")
	}
	request := ess.CreateDescribeScalingConfigurationsRequest()
	request.ScalingGroupId = scalingGroupID
	response, err := m.ess.DescribeScalingConfigurations(request)
	if err != nil {
		return nil, err
	}
	return response.ScalingConfigurations.ScalingConfiguration, nil
}

func (m *autoScalingWrapper) createInstance(configuration ess.ScalingConfiguration) (string, error) {
	reqCI := ecs.CreateCreateInstanceRequest()
	reqCI.RegionId = m.cfg.Global.Region
	reqCI.ImageId = configuration.ImageId
	reqCI.InstanceType = configuration.InstanceType
	reqCI.SecurityGroupId = configuration.SecurityGroupId
	respCI, err := m.ecs.CreateInstance(reqCI)
	if err != nil {
		return "", err
	}
	return respCI.InstanceId, nil
}

// 删除实例
func (m *autoScalingWrapper) removeInstances(scalingGroupID string, instances []string) (string, error) {
	request := ess.CreateRemoveInstancesRequest()
	request.ScalingGroupId = scalingGroupID
	for i := 0; i < len(instances); i++ {
		// reset instanceID to sdk format. x-xxxxxxx(when get instance we format xzxxxxxxxz)
		sdkInstance := instances[i][:1]+"-"+instances[i][2:len(instances[i])-1]
		err := SetInfo(request, "InstanceId"+strconv.Itoa(i+1), sdkInstance)
		if err != nil {
			return "", err
		}
	}
	response, err := m.ess.RemoveInstances(request)
	if err != nil {
		return "", err
	}
	return response.ScalingActivityId, nil
}

func (m *autoScalingWrapper) getAutoscalingGroupActivitiesByGroupID(scalingGroupID string) ([]ess.ScalingActivity, error) {
	request := ess.CreateDescribeScalingActivitiesRequest()
	request.ScalingGroupId = scalingGroupID
	response, err := m.ess.DescribeScalingActivities(request)
	if err != nil {
		return nil, err
	}
	return response.ScalingActivities.ScalingActivity, nil
}

func (m *autoScalingWrapper) getAutoscalingGroupActivities(scalingGroupID string, scalingActivityID string) ([]ess.ScalingActivity, error) {
	request := ess.CreateDescribeScalingActivitiesRequest()
	request.ScalingGroupId = scalingGroupID
	request.ScalingActivityId1 = scalingActivityID
	response, err := m.ess.DescribeScalingActivities(request)
	if err != nil {
		return nil, err
	}
	return response.ScalingActivities.ScalingActivity, nil
}
