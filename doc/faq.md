# FAQ

## operator 支持的 flink 版本范围是什么？

理论上，从 `flink 1.13.5` 以后的版本均可支持使用此 `operator` 部署。但是为了部署兼容性及稳定性考虑，建议使用 `flink 1.14.4` 或者更高版本

## 是否可以不重启 session 集群动态修改 flink 集群的资源配置参数？

如果要想实现此功能，仅控制 `k8s` 是不够的，还需要 `flink` 自身具备如下能力:

1. 重新加载配置文件
2. 不重启 `jvm` ，动态调整 `jvm` 内存分配

目前，flink 并不能提供上述能力，因此，相关功能在实现时，暂时只能通过重启 `flink session` 集群的方式进行

## 重启 flink 集群会对集群中运行的任务产生何种影响？

这个问题分为两种情况：

1. 开启 ha  
   在开启 ha 的情况下，`operator` 并没有操作 `flink` 集群关闭已有的 `task` ，因此集群的重启会被 `flink` 处理为常规的 `jobmanager` 异常，并在新集群启动后重启启动相关任务。在 `autoClean` 参数配置为 `false` 的情况下，可以理解为对运行任务无严重影响，可以依赖 `flink` 的 `ha` 机制正常重启；如果 `autoClean` 参数被配置为 `true`，则会在集群关闭后清理 `ha` 信息，新启动的集群与未开启 `ha` 时保持一致

2. 未开启 ha  
   如果未开启 `ha` ，则重新启动 `flink` 集群与 `jobmanager` 重启效果一致，全部任务会丢失，需要重新提交任务

因此，建议在修改集群配置前，手动停止 `flink session` 集群中所有任务，并在修改后重新提交，以确保正常运行

## 是否会支持 kerberos 认证？

目前最新版本已经可以支持 `kerberos` 认证。

技术方案分析：

当前 `Flink` 对 `Kerberos` 认证的支持主要体现在两部分:

- `Kafka` 连接支持 `Kerberos` 认证
- `Flink` 本身使用 `Zookeeper` 支持 `Kerberos` 认证

对于本工具来说，因为需要操作 `Flink` 的 `HA` 信息，因此必须解决 `Zookeeper` 连接的 `Kerberos` 认证问题。当前社区，可以找到的支持 `go-zookeeper` 的可行方案有:

- [gosasl-zk](https://github.com/wenbingshen/zk/tree/kerb)
- [gokrb5-zk](https://github.com/zhuliquan/zk)

第一种方案是使用 [gosasl](https://github.com/beltran/gosasl) 库，依赖 `kinit` 命令进行认证；第二种方案是使用 [gokrb5](https://github.com/jcmturner/gokrb5) 这种原生方案进行认证。

在 `operator` 场景中，更适合使用第二种原生方案。因此，可以在 [gokrb5-zk](https://github.com/zhuliquan/zk) 方案的基础上，针对 `Flink` 连接 `Kerberos` 的特殊场景进行针对性改造。具体修改可以参考

- [zk-kerberos](https://github.com/123shang60/zk/tree/kerberos)

