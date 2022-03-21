# flink-session-operator
flink session 集群部署

[![Go](https://github.com/123shang60/flink-session-operator/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/123shang60/flink-session-operator/actions/workflows/go.yml)
[![CodeQL](https://github.com/123shang60/flink-session-operator/actions/workflows/codeql-analysis.yml/badge.svg?branch=main)](https://github.com/123shang60/flink-session-operator/actions/workflows/codeql-analysis.yml)

基于 kubebuilder 构建的、方便管理与创建 flink session 集群的 operator

主要能力：

- 自动化 flink session 集群部署流程
- 自动化修改 flink session 集群资源配置
- 自动化清理 flink session 集群

目录

- [flink-session-operator](#flink-session-operator)
  - [部署方法](#部署方法)
  - [实现原理](#实现原理)
  - [使用方法](#使用方法)
    - [新增部署集群](#新增部署集群)
    - [修改集群配置](#修改集群配置)
    - [卸载集群](#卸载集群)
  - [参数说明](#参数说明)
  - [特性功能](#特性功能)
    - [配置变更时清理 ha 及状态后端信息](#配置变更时清理-ha-及状态后端信息)
  - [注意事项](#注意事项)
  - [FAQ](#faq)

## 部署方法

1. 部署 [cert-manager](https://cert-manager.io/docs/installation/)
2. 执行命令  `kubectl apply -f https://github.com/123shang60/flink-session-operator/releases/download/v0.2.0/install.yaml`

## 实现原理

部署原理参考 [flink 官方文档](https://nightlies.apache.org/flink/flink-docs-release-1.14/docs/deployment/resource-providers/native_kubernetes/)，配置来源为官方文档。

修改资源配置生效的原理为，将原有集群卸载，并根据最新的配置动态生成启动 `job`，以达到自动修改集群资源配置的功能

## 使用方法

### 新增部署集群

1. 首先需要单独准备运行 flink 所需的外部组件，例如 s3\zookeeper 等，确保 k8s 与相关组件的连通性；
2. 配置部署的 `namespaces` 以及 `service account` 等基础配置，可以参考 [示例](./config/samples/rbac.yaml)
3. 配置 `flinkSession` 对象，可以参考 [示例](./config/samples/flink_v1_flinksession.yaml)
4. 使用 `kubectl`  或其他方式 `apply` 相关配置，即可将集群部署到 k8s 内部

### 修改集群配置

直接使用 `kubectl edit flinkSession ` 命令修改 `flinkSession` 对象的相关内容，保存后即可生效

### 卸载集群

直接使用 `kubectl delete flinkSession ` 命令删除 `flinkSession` 对象，即可将集群完全卸载

## 参数说明

详细参数说明见 [参数说明](./doc/param.md)

## 特性功能

### 配置变更时清理 ha 及状态后端信息

功能开关为 `{{spec.autoClean}}` 配置项。此配置项开启时，每次修改 flink 集群配置或者重新部署时，都会自动清理 `flink` 残留的 `ha` 及状态后端信息，实现完全重新部署 `session` 集群的目的；当功能开关为关闭的情况下，修改或者重新部署将不对 `ha` 信息做任何修改。

## 注意事项

1. 建议使用 `minio` 等 `s3` 后端作为状态后端及 `ha` 配置存储，暂不支持 `hadoop` 模式；
2. 作为状态后端的 `bucket` 请务必保证不被其他程序使用。开启清理 `ha` 信息功能会，重启 `flink` 集群时对应 `bucket` 数据会被完整删除

## FAQ

参见 [FAQ](./doc/faq.md)