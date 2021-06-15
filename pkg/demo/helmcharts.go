package demo

import (
	helmv1 "github.com/k3s-io/helm-controller/pkg/apis/helm.cattle.io/v1"
	demov1alpha1 "github.com/rancher/opni/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func BuildMinioHelmChart(spec *demov1alpha1.OpniDemo) *helmv1.HelmChart {
	values := map[string]intstr.IntOrString{
		"accessKey":                intstr.FromString(spec.Spec.MinioAccessKey),
		"secretKey":                intstr.FromString(spec.Spec.MinioSecretKey),
		"persistence.storageClass": intstr.FromString("local-path"),
	}
	for k, v := range spec.Spec.Components.Opni.Minio.Set {
		values[k] = v
	}
	return &helmv1.HelmChart{
		ObjectMeta: v1.ObjectMeta{
			Name: "minio",
		},
		Spec: helmv1.HelmChartSpec{
			Chart:   "minio",
			Repo:    "https://helm.min.io/",
			Version: spec.Spec.MinioVersion,
			Set:     values,
		},
	}
}

func BuildNatsHelmChart(spec *demov1alpha1.OpniDemo) *helmv1.HelmChart {
	values := map[string]intstr.IntOrString{
		"auth.enabled":  intstr.FromString("true"),
		"auth.password": intstr.FromString(spec.Spec.NatsPassword),
		"replicaCount":  intstr.FromInt(spec.Spec.NatsReplicas),
		"maxPayload":    intstr.FromInt(spec.Spec.NatsMaxPayload),
	}
	for k, v := range spec.Spec.Components.Opni.Nats.Set {
		values[k] = v
	}
	return &helmv1.HelmChart{
		ObjectMeta: v1.ObjectMeta{
			Name: "nats",
		},
		Spec: helmv1.HelmChartSpec{
			Chart: "nats",
			Repo:  "https://charts.bitnami.com/bitnami",
			Set:   values,
		},
	}
}

func BuildElasticHelmChart(spec *demov1alpha1.OpniDemo) *helmv1.HelmChart {
	values := map[string]intstr.IntOrString{
		"elasticsearch.master.persistence.enabled":      intstr.FromString("true"),
		"elasticsearch.master.persistence.storageClass": intstr.FromString("local-path"),
		"elasticsearch.data.persistence.enabled":        intstr.FromString("true"),
		"elasticsearch.data.persistence.storageClass":   intstr.FromString("local-path"),
		"elasticsearch.username":                        intstr.FromString(spec.Spec.ElasticsearchUser),
		"elasticsearch.password":                        intstr.FromString(spec.Spec.ElasticsearchPassword),
	}
	for k, v := range spec.Spec.Components.Opni.Elastic.Set {
		values[k] = v
	}
	return &helmv1.HelmChart{
		ObjectMeta: v1.ObjectMeta{
			Name: "opendistro-es",
		},
		Spec: helmv1.HelmChartSpec{
			Chart:   "opendistro-es",
			Repo:    "https://raw.githubusercontent.com/rancher/opni-charts/main",
			Version: "1.13.201+up1.13.2",
			Set:     values,
		},
	}
}

func BuildRancherLoggingCrdHelmChart() *helmv1.HelmChart {
	return &helmv1.HelmChart{
		ObjectMeta: v1.ObjectMeta{
			Name: "rancher-logging-crd",
		},
		Spec: helmv1.HelmChartSpec{
			Chart:   "rancher-logging-crd",
			Repo:    "https://raw.githubusercontent.com/rancher/opni-charts/main",
			Version: "3.10.0",
		},
	}
}

func BuildRancherLoggingHelmChart(spec *demov1alpha1.OpniDemo) *helmv1.HelmChart {
	values := map[string]intstr.IntOrString{
		"additionalLoggingSources.k3s.enabled": intstr.FromString("true"),
	}
	for k, v := range spec.Spec.Components.Opni.RancherLogging.Set {
		values[k] = v
	}
	return &helmv1.HelmChart{
		ObjectMeta: v1.ObjectMeta{
			Name: "rancher-logging",
		},
		Spec: helmv1.HelmChartSpec{
			Chart:   "rancher-logging",
			Repo:    "https://raw.githubusercontent.com/rancher/opni-charts/main",
			Version: "3.10.0",
			Set:     values,
		},
	}
}

func BuildTraefikHelmChart(spec *demov1alpha1.OpniDemo) *helmv1.HelmChart {
	values := map[string]intstr.IntOrString{
		"ports.websecure.nodePort": intstr.FromInt(32222),
	}
	for k, v := range spec.Spec.Components.Opni.Traefik.Set {
		values[k] = v
	}
	return &helmv1.HelmChart{
		ObjectMeta: v1.ObjectMeta{
			Name: "traefik",
		},
		Spec: helmv1.HelmChartSpec{
			Chart:   "traefik",
			Repo:    "https://helm.traefik.io/traefik",
			Version: spec.Spec.TraefikVersion,
			Set:     values,
		},
	}
}
