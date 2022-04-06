package v1

import (
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func TestBuildNodeSelector(t *testing.T) {
	maps := make(map[string]string)

	res := buildNodeSelector(maps)

	assert.Equal(t, "", res)

	maps["flink"] = "run"

	res = buildNodeSelector(maps)

	assert.Equal(t, "flink:run", res)

	maps["disk"] = "ssd"

	res = buildNodeSelector(maps)

	if res != "flink:run,disk:ssd" && res != "disk:ssd,flink:run" {
		assert.Fail(t, "map 拼接失败！")
	}
}

func TestPodTemplate(t *testing.T) {
	session := &FlinkSession{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flink",
			Namespace: "test",
		},
		Spec: FlinkSessionSpec{
			BalancedSchedule: PreferredDuringScheduling,
		},
	}

	template := session.GeneratePodTemplate()
	assert.Equal(t, Expected, template)
}
