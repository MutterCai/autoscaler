# CA的阿里云支持

CA原本的MD文件改名成READMD_BACKUP.md

* 协定：在阿里云上创建的伸缩组名称，必须以Region开头  
例如香港的伸缩组(Region=cn-hongkong): cn-hongkong.kubernetes-hongkong-group

## 调试

#### 阿里云sdk
* go get -u github.com/aliyun/alibaba-cloud-sdk-go/sdk

#### 环境变量
* AccessKeyID  
    进入阿里云，账户上的AccessKey，进入后可以获取到相关的信息
* AccessKeySecret  
    进入阿里云，账户上的AccessKey，进入后可以获取到相关的信息
* Region  
    根据阿里云的ess部署的区域设置
    
#### 启动参数
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
    
## docker和k8s

#### docker镜像制作
1. go build main.go  
    编译ca二进制文件
2. mv main cluster-autoscaler  
    重命名二进制文件
3. cp $GOROOT/lib/time/zoneinfo.zip .  
    将zoneinfo拉到目录下，供dockerfile打包使用；zoneinfo.zip是go的时区库
4. sudo docker build -t dockerHubRepository:version  
    创建docker镜像
5. sudo docker push  
    推送docker镜像到仓库

#### k8s嵌入
1. kubectl apply -f ca-sa.yaml  
    创建CA的serverAccount
    ```yaml
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: cluster-autoscaler
      namespace: default
    ```
2. kubectl apply -f ali-ca.yaml  
    创建CA的deployment
    ```yaml
    apiVersion: extensions/v1beta1
    kind: Deployment
    metadata:
      name: cluster-autoscaler
      namespace: kube-system
      labels:
        app: cluster-autoscaler
    spec:
      replicas: 1
      selector:
        matchLabels:
          app: cluster-autoscaler
      template:
        metadata:
          labels:
            app: cluster-autoscaler
        spec:
          containers:
            - image: muttercai/cluster-autoscaler:v1.2
              name: cluster-autoscaler
              resources:
                limits:
                  cpu: 100m
                  memory: 300Mi
                requests:
                  cpu: 100m
                  memory: 300Mi
              command:
                - ./cluster-autoscaler
                - --v=4
                - --cloud-provider=ali
                - --skip-nodes-with-local-storage=false
                - --node-group-auto-discovery=asg:tag=cn-hongkong.kubernetes-hongkong-group,cn-hongkong.kubernetes-hongkong-group2
              env:
              - name: Region
                value: cn-hongkong
              - name: AccessKeyID
                value: LTAITh6uSXRCv14k
              - name: AccessKeySecret
                value: zm4qbsuKnDQ3VRIxAQzFKvMt9aOBuI
              imagePullPolicy: "Always"
    ```
    
