# flink-session-operator
flink session 集群部署

[![Go](https://github.com/123shang60/flink-session-operator/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/123shang60/flink-session-operator/actions/workflows/go.yml)
[![CodeQL](https://github.com/123shang60/flink-session-operator/actions/workflows/codeql-analysis.yml/badge.svg?branch=main)](https://github.com/123shang60/flink-session-operator/actions/workflows/codeql-analysis.yml)

基于 kubebuilder 构建的、方便管理与创建 flink session 集群的 operator

当前实现进度：

- [x] controllers 脚手架构建
- [x] webhook 适配
- [x] finalizer 适配
- [x] 解决 operator 自身 update 对象时重复调用逻辑问题
- [x] 真实场景 flink crd 配置构建
- [x] deleteExternalResources 函数实现集群卸载
- [x] updateExternalResources 函数实现集群清理 + 集群部署
- [x] k8s event 事件记录
- [x] status 展示
- [x] k8s 原生 ha 支持
- [x] 自定义 config 配置文件能力
- [x] 支持基于 pod-template 的多种均衡节点调度策略
- [x] 支持可选更新时删除 ha 及 minio 状态后端信息

未完成进度：

- [ ] webhook 校验能力

需要支持的扩展能力：

- [ ] 多架构构建适配
- [ ] 支持 flink 集群 kerberos 认证
