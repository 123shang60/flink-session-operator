package controllers

import (
	flinkv1 "github.com/123shang60/flink-session-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	Expected = `apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
spec:
  affinity:
    podAntiAffinity:
      preferredDuringSchedulingIgnoredDuringExecution:
      - podAffinityTerm:
          labelSelector:
            matchExpressions:
            - key: app
              operator: In
              values:
              - flink
            - key: type
              operator: In
              values:
              - flink-native-kubernetes
          namespaces:
          - test
          topologyKey: kubernetes.io/hostname
        weight: 100
  containers: []
status: {}
`
)

func TestPodTemplate(t *testing.T) {
	res := generatePodTemplate(flinkv1.PreferredDuringScheduling, "flink", "test")
	assert.Equal(t, Expected, res)
}
