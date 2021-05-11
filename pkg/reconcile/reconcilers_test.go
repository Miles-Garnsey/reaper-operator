package reconcile

import (
	"testing"

	api "github.com/k8ssandra/reaper-operator/api/v1alpha1"
	mlabels "github.com/k8ssandra/reaper-operator/pkg/labels"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestNewService(t *testing.T) {
	reaper := newReaperWithCassandraBackend()
	key := types.NamespacedName{Namespace: reaper.Namespace, Name: GetServiceName(reaper.Name)}

	service := newService(key, reaper)

	assert.Equal(t, key.Name, service.Name)
	assert.Equal(t, key.Namespace, service.Namespace)
	assert.Equal(t, createLabels(reaper), service.Labels)

	assert.Equal(t, createLabels(reaper), service.Spec.Selector)
	assert.Equal(t, 1, len(service.Spec.Ports))

	port := corev1.ServicePort{
		Name:     "app",
		Protocol: corev1.ProtocolTCP,
		Port:     8080,
		TargetPort: intstr.IntOrString{
			Type:   intstr.String,
			StrVal: "app",
		},
	}
	assert.Equal(t, port, service.Spec.Ports[0])
}

func TestNewDeployment(t *testing.T) {
	assert := assert.New(t)
	image := "test/reaper:latest"
	reaper := newReaperWithCassandraBackend()
	reaper.Spec.Image = image
	reaper.Spec.ImagePullPolicy = "Always"
	reaper.Spec.ServerConfig.AutoScheduling = &api.AutoScheduler{Enabled: true}

	labels := createLabels(reaper)
	deployment := newDeployment(reaper, "target-datacenter-service")

	assert.Equal(reaper.Namespace, deployment.Namespace)
	assert.Equal(reaper.Name, deployment.Name)
	assert.Equal(labels, deployment.Labels)

	selector := deployment.Spec.Selector
	assert.Equal(0, len(selector.MatchLabels))
	assert.ElementsMatch(selector.MatchExpressions, []metav1.LabelSelectorRequirement{
		{
			Key:      mlabels.ManagedByLabel,
			Operator: metav1.LabelSelectorOpIn,
			Values:   []string{mlabels.ManagedByLabelValue},
		},
		{
			Key:      mlabels.ReaperLabel,
			Operator: metav1.LabelSelectorOpIn,
			Values:   []string{reaper.Name},
		},
	})

	assert.Equal(labels, deployment.Spec.Template.Labels)

	podSpec := deployment.Spec.Template.Spec
	assert.Equal(1, len(podSpec.Containers))

	container := podSpec.Containers[0]
	assert.Equal(image, container.Image)
	assert.Equal(corev1.PullAlways, container.ImagePullPolicy)
	assert.ElementsMatch(container.Env, []corev1.EnvVar{
		{
			Name:  "REAPER_STORAGE_TYPE",
			Value: "cassandra",
		},
		{
			Name:  "REAPER_ENABLE_DYNAMIC_SEED_LIST",
			Value: "false",
		},
		{
			Name:  "REAPER_CASS_CONTACT_POINTS",
			Value: "[target-datacenter-service]",
		},
		{
			Name:  "REAPER_AUTH_ENABLED",
			Value: "false",
		},
		{
			Name:  "REAPER_AUTO_SCHEDULING_ENABLED",
			Value: "true",
		},
	})

	reaper.Spec.ServerConfig.AutoScheduling = &api.AutoScheduler{
		Enabled:              false,
		InitialDelay:         "PT10S",
		PeriodBetweenPolls:   "PT5M",
		BeforeFirstSchedule:  "PT10M",
		ScheduleSpreadPeriod: "PT6H",
		ExcludedClusters:     []string{"a", "b"},
		ExcludedKeyspace:     []string{"system.powers"},
	}

	deployment = newDeployment(reaper, "target-datacenter-service")
	podSpec = deployment.Spec.Template.Spec
	container = podSpec.Containers[0]
	assert.Equal(4, len(container.Env))

	reaper.Spec.ServerConfig.AutoScheduling.Enabled = true
	deployment = newDeployment(reaper, "target-datacenter-service")
	podSpec = deployment.Spec.Template.Spec
	container = podSpec.Containers[0]
	assert.Equal(11, len(container.Env))

	assert.Contains(container.Env, corev1.EnvVar{
		Name:  "REAPER_AUTO_SCHEDULING_PERIOD_BETWEEN_POLLS",
		Value: "PT5M",
	})

	assert.Contains(container.Env, corev1.EnvVar{
		Name:  "REAPER_AUTO_SCHEDULING_TIME_BEFORE_FIRST_SCHEDULE",
		Value: "PT10M",
	})

	assert.Contains(container.Env, corev1.EnvVar{
		Name:  "REAPER_AUTO_SCHEDULING_INITIAL_DELAY_PERIOD",
		Value: "PT10S",
	})

	assert.Contains(container.Env, corev1.EnvVar{
		Name:  "REAPER_AUTO_SCHEDULING_EXCLUDED_CLUSTERS",
		Value: "[a, b]",
	})

	assert.Contains(container.Env, corev1.EnvVar{
		Name:  "REAPER_AUTO_SCHEDULING_EXCLUDED_KEYSPACES",
		Value: "[system.powers]",
	})

	probe := &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/healthcheck",
				Port: intstr.FromInt(8081),
			},
		},
		InitialDelaySeconds: 45,
		PeriodSeconds:       15,
	}
	assert.Equal(probe, container.LivenessProbe)
	assert.Equal(probe, container.ReadinessProbe)
}

func TestTolerations(t *testing.T) {
	image := "test/reaper:latest"
	tolerations := []corev1.Toleration{
		{
			Key:      "key1",
			Operator: corev1.TolerationOpEqual,
			Value:    "value1",
			Effect:   corev1.TaintEffectNoSchedule,
		},
		{
			Key:      "key2",
			Operator: corev1.TolerationOpEqual,
			Value:    "value2",
			Effect:   corev1.TaintEffectNoSchedule,
		},
	}

	reaper := newReaperWithCassandraBackend()
	reaper.Spec.Image = image
	reaper.Spec.Tolerations = tolerations

	deployment := newDeployment(reaper, "target-datacenter-service")

	assert.ElementsMatch(t, tolerations, deployment.Spec.Template.Spec.Tolerations)
}

func newReaperWithCassandraBackend() *api.Reaper {
	namespace := "service-test"
	reaperName := "test-reaper"
	dcName := "dc1"

	return &api.Reaper{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      reaperName,
		},
		Spec: api.ReaperSpec{
			ServerConfig: api.ServerConfig{
				StorageType: api.StorageTypeCassandra,
				CassandraBackend: &api.CassandraBackend{
					CassandraDatacenter: api.CassandraDatacenterRef{
						Name: dcName,
					},
					Keyspace: api.DefaultKeyspace,
					Replication: api.ReplicationConfig{
						NetworkTopologyStrategy: &map[string]int32{
							"DC1": 3,
						},
					},
				},
			},
		},
	}
}
