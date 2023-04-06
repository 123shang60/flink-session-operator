package v1

import (
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

const (
	Expected1 = `apiVersion: v1
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
  containers:
  - name: flink-main-container
    resources: {}
status: {}
`
	Expected2 = `apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
spec:
  affinity:
    podAntiAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
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
  containers:
  - name: flink-main-container
    resources: {}
    volumeMounts:
    - mountPath: /tmp
      name: test
  volumes:
  - emptyDir: {}
    name: test
status: {}
`

	ExpectCommand1 = `$FLINK_HOME/bin/kubernetes-session.sh -Dkubernetes.cluster-id=flink -Dkubernetes.namespace=test -Dtaskmanager.numberOfTaskSlots=0 -Dstate.backend=filesystem -Ds3.endpoint=http:// -Ds3.path.style.access=true -Dstate.checkpoints.dir=s3:///flink/flink/checkpoints -Dstate.savepoints.dir=s3:///flink/flink/savepoints -Dhistoryserver.archive.fs.dir=s3:///flink/flink/completed-jobs -Djobmanager.archive.fs.dir=s3:///flink/flink/archive -Dstate.backend.incremental=true -Dfs.overwrite-files=true -Dhigh-availability=org.apache.flink.kubernetes.highavailability.KubernetesHaServicesFactory -Dhigh-availability.storageDir=s3:///flink/flink/ha/metadata -Denv.java.opts="-XX:+UseG1GC" -Dkubernetes.rest-service.exposed.type=NodePort `
	ExpectCommand2 = `$FLINK_HOME/bin/kubernetes-session.sh -Dkubernetes.cluster-id=flink -Dkubernetes.namespace=test -Dtaskmanager.numberOfTaskSlots=0 -Dstate.backend=filesystem -Ds3.endpoint=http:// -Ds3.path.style.access=true -Dstate.checkpoints.dir=s3:///flink/flink/checkpoints -Dstate.savepoints.dir=s3:///flink/flink/savepoints -Dhistoryserver.archive.fs.dir=s3:///flink/flink/completed-jobs -Djobmanager.archive.fs.dir=s3:///flink/flink/archive -Dstate.backend.incremental=true -Dfs.overwrite-files=true -Dhigh-availability=KUBERNETES -Dhigh-availability.storageDir=s3:///flink/flink/ha/metadata -Denv.java.opts="-XX:+UseG1GC" -Dkubernetes.rest-service.exposed.type=NodePort `
	ExpectCommand3 = `$FLINK_HOME/bin/kubernetes-session.sh -Dkubernetes.cluster-id=flink -Dkubernetes.namespace=test -Dtaskmanager.numberOfTaskSlots=0 -Dstate.backend=filesystem -Ds3.endpoint=http:// -Ds3.path.style.access=true -Dstate.checkpoints.dir=s3:///flink/flink/checkpoints -Dstate.savepoints.dir=s3:///flink/flink/savepoints -Dhistoryserver.archive.fs.dir=s3:///flink/flink/completed-jobs -Djobmanager.archive.fs.dir=s3:///flink/flink/archive -Dstate.backend.incremental=true -Dfs.overwrite-files=true -Dhigh-availability.type=KUBERNETES -Dhigh-availability.storageDir=s3:///flink/flink/ha/metadata -Denv.java.opts="-XX:+UseG1GC" -Dkubernetes.rest-service.exposed.type=NodePort `
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
	assert.Equal(t, Expected1, template)

	session = &FlinkSession{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flink",
			Namespace: "test",
		},
		Spec: FlinkSessionSpec{
			BalancedSchedule: RequiredDuringScheduling,
			Volumes: []apiv1.Volume{
				{
					Name: "test",
					VolumeSource: apiv1.VolumeSource{
						EmptyDir: &apiv1.EmptyDirVolumeSource{},
					},
				},
			},
			VolumeMounts: []apiv1.VolumeMount{
				{
					Name:      "test",
					MountPath: "/tmp",
				},
			},
		},
	}

	template = session.GeneratePodTemplate()
	assert.Equal(t, Expected2, template)
}

func TestHAFlink(t *testing.T) {
	session := &FlinkSession{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flink",
			Namespace: "test",
		},
		Spec: FlinkSessionSpec{
			FlinkVersion: nil,
			HA: FlinkHA{
				Typ:    "kubernetes",
				Quorum: "",
				Path:   "",
			},
			BalancedSchedule: NoneScheduling,
		},
	}

	template, err := session.GenerateCommand()
	assert.Nil(t, err)
	assert.Equal(t, ExpectCommand1, template)

	version := "1.16"
	session = &FlinkSession{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flink",
			Namespace: "test",
		},
		Spec: FlinkSessionSpec{
			FlinkVersion: &version,
			HA: FlinkHA{
				Typ:    "kubernetes",
				Quorum: "",
				Path:   "",
			},
			BalancedSchedule: NoneScheduling,
		},
	}

	template, err = session.GenerateCommand()
	assert.Nil(t, err)
	assert.Equal(t, ExpectCommand2, template)

	version = "1.17"
	session = &FlinkSession{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flink",
			Namespace: "test",
		},
		Spec: FlinkSessionSpec{
			FlinkVersion: &version,
			HA: FlinkHA{
				Typ:    "kubernetes",
				Quorum: "",
				Path:   "",
			},
			BalancedSchedule: NoneScheduling,
		},
	}

	template, err = session.GenerateCommand()
	assert.Nil(t, err)
	assert.Equal(t, ExpectCommand3, template)

	version = "1.15"
	session = &FlinkSession{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flink",
			Namespace: "test",
		},
		Spec: FlinkSessionSpec{
			FlinkVersion: &version,
			HA: FlinkHA{
				Typ:    "kubernetes",
				Quorum: "",
				Path:   "",
			},
			BalancedSchedule: NoneScheduling,
		},
	}

	template, err = session.GenerateCommand()
	assert.Nil(t, err)
	assert.Equal(t, ExpectCommand1, template)
}
