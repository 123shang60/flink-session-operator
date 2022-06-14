# 参数说明

`spec` :

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|image|string|""| |填写 `flink` 运行镜像 ,这个镜像为 `flink` `taskmanager` 以及 `jobmanager` 共用，不支持双镜像，必填|  
|imageSecret|string|null| |填写 `flink` 镜像拉取 `secret`|  
|autoClean|bool|false||是否在修改配置重启 `flink` 时自动清除 `ha` 信息，必填|
|sa|string|""||填写集群运行的 `k8s service account` 配置|
|resource|FlinkResource|||`flink` 运行资源配置|
|numberOfTaskSlots|int|1||`taskmanager` 可用槽位，对应 `taskmanager.numberOfTaskSlots`|
|s3|FlinkS3|||`S3` 配置|
|ha|FlinkHA|||`flink` `ha` 配置|
|config|FlinkConfig|||`Flink` 自定义配置项|
|nodeSelector|map[string]string|null||为 `flink` 增加基于 `node label` 的 `nodeSelector`|
|balancedSchedule|enum|None|{Required,Preferred,None}|均衡调度策略 ，可选值： `Required` 必须每个节点调度一个 `Preferred` 尽可能每个节点调度一个 `None` 不设置均衡调度|
|volumes|[]apiv1.Volume|null||为 `flink` 增加 卷挂载|
|volumeMounts|[]apiv1.VolumeMount|null||为 `flink` 所有 `container` 配置卷挂载|
|security|FlinkSecurity|null||`flink` 安全性配置|

`FlinkResource` :

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|jobManager|JobManagerFlinkResource|||`jobmanager` 资源配置|
|taskManager|TaskManagerFlinkResource|||`taskmanager` 资源配置|

`JobManagerFlinkResource` :

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|cpu|string|""||等价 `flink` 原始配置 `kubernetes.jobmanager.cpu`|
|memory|string|""||等价 `flink` 原始配置 `jobmanager.memory.process.size`|
|jvm-metaspace|string|""||等价 `flink` 原始配置 `jobmanager.memory.jvm-metaspace.size`|
|off-heap|string|""||等价 `flink` 原始配置 `jobmanager.memory.off-heap.size`|

`TaskManagerFlinkResource` :

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|cpu|string|""||等价 `flink` 原始配置 `kubernetes.taskmanager.cpu`|
|memory|string|""||等价 `flink` 原始配置 `taskmanager.memory.process.size`|
|jvm-metaspace|string|""||等价 `flink` 原始配置 `taskmanager.memory.jvm-metaspace.size`|
|framework|TaskManagerFrameworkFlinkResource||| `farmework` 资源配置|
|task|TaskManagerTaskFlinkResource||| `task` 资源配置|
|netWork|TaskManagerNetWorkFlinkResource||| `network` 资源配置|
|managed|TaskManagerManagedFlinkResource||| `managed` 资源配置|

`TaskManagerFrameworkFlinkResource` :

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|heap|string|||等价 `flink` 原始配置 `taskmanager.memory.framework.heap.size`|
|off-heap|string|||等价 `flink` 原始配置 `taskmanager.memory.framework.off-heap.size`|

`TaskManagerTaskFlinkResource` :

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|off-heap|string|||等价 `flink` 原始配置 `taskmanager.memory.task.off-heap.size`|

`TaskManagerNetWorkFlinkResource` :

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|fraction|string|||等价 `flink` 原始配置 `taskmanager.memory.network.fraction` ， 与 `max` `min` 仅其中一个配置生效，默认走`max` `min`|
|min|string|||等价 `flink` 原始配置 `taskmanager.memory.network.min` ， 与 `max` `min` 仅其中一个配置生效，默认走`max` `min`|
|max|string|||等价 `flink` 原始配置 `taskmanager.memory.network.max` ， 与 `max` `min` 仅其中一个配置生效，默认走`max` `min`|

`TaskManagerManagedFlinkResource` :

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|fraction|string|||等价 `flink` 原始配置 `taskmanager.memory.managed.fraction` ， 与 `max` `min` 仅其中一个配置生效，默认走`max` `min`|
|min|string|||等价 `flink` 原始配置 `taskmanager.memory.managed.min` ， 与 `max` `min` 仅其中一个配置生效，默认走`max` `min`|
|max|string|||等价 `flink` 原始配置 `taskmanager.memory.managed.max` ， 与 `max` `min` 仅其中一个配置生效，默认走`max` `min`|

`FlinkS3` : 

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|endPoint|string|||等价 `flink` 原始配置 `s3.endpoint`|
|accessKey|string|||等价 `flink` 原始配置 `s3.access-key`|
|secretKey|string|||等价 `flink` 原始配置 `s3.secret-key`|
|bucket|string|||指定 `flink` 部署使用的 `bucket`，当 `bucket` 不存在时会自动创建。|

`FlinkHA` :

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|type|enum|none|{zookeeper,kubernetes,none}|指定 `flink` 集群 `ha` 开启情况，允许使用 `zookeeper` 或者 `org.apache.flink.kubernetes.highavailability.KubernetesHaServicesFactory`|
|quorum|string|||仅 zookeeper ha 生效，配置 zk 地址|
|path|string|||仅 zookeeper ha 生效，配置 ha 路径前缀|

`FlinkConfig` :

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|flink-conf.yaml|string|||对应 `$FLINK_HOME/conf/flink-conf.yaml` ，不写使用镜像内预制配置文件|
|log4j-console.properties||||对应 `$FLINK_HOME/conf/log4j-console.properties` ，不写使用镜像内预制配置文件|
|logback-console.xml||||对应 `$FLINK_HOME/conf/logback-console.xml` ，不写使用镜像内预制配置文件|

`FlinkSecurity` :

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|kerberos|KerberosConf|||`Kerberos` 相关配置|

`KerberosConf` :

|参数|类型|默认值|枚举值|说明|  
|-|-|-|-|-|
|krb5|string|||`Krb5` 文件|
|contexts|string||`Client`,`KafkaClient`,`Client,KafkaClient`,`KafkaClient,Client`|登录上下文列表，根据 `Flink` 文档支持 `ZK` 及 `KAFKA`|
|principal|string|||`Kerberos` 主体名称|
|base64Keytab|string|||`base64` 编码的 `Keytab` 文件|
|useTicketCache|bool|`true`||`UseTicketCache`，是否使用 `Ticket` 缓存|

`status` :

|参数|类型|说明|  
|-|-|-|
|ready|bool|`operator` 操作是否已经完成|
|port|int|`flink session` 集群对外 `nodeport` 端口号|
