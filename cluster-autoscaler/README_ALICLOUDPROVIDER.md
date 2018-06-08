# CA的阿里云支持

## 阿里云sdk
* go get -u github.com/aliyun/alibaba-cloud-sdk-go/sdk

## 环境变量
* AccessKeyID
    进入阿里云，账户上的AccessKey，进入后可以获取到相关的信息
* AccessKeySecret
    进入阿里云，账户上的AccessKey，进入后可以获取到相关的信息
* Region
    根据阿里云的ess部署的区域设置
    
## 启动参数
* --v=4
    设置日志的等级，越高日志等级越低 
* --cloud-provider=ali
    设置cloud provider的归属，这里使用ali
* --skip-nodes-with-local-storage=false 
    忽略本地存储节点
* --node-group-auto-discovery=asg:tag=cn-hongkong.kubernetes-hongkong-group,cn-hongkong.kubernetes-hongkong-group2 
    设置管理的伸缩组，多个组用','间隔
* --kubeconfig kubeconfigFile
    设置管理的K8S集群，对应的kubeconfigFile为k8s的master的~/.kube/config文件
    注意如果是外网应用应该将server的ip改成公网ip