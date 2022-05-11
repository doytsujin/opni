package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	loggingv1beta1 "github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/output"
	"github.com/cenkalti/backoff"
	"github.com/rancher/opni/apis/v1beta2"
	"github.com/rancher/opni/pkg/bootstrap"
	"github.com/rancher/opni/pkg/capabilities/wellknown"
	gatewayclients "github.com/rancher/opni/pkg/clients"
	"github.com/rancher/opni/pkg/opni/common"
	"github.com/rancher/opni/pkg/tokens"
	"github.com/rancher/opni/pkg/trust"
	loggingplugin "github.com/rancher/opni/plugins/logging/pkg/logging"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	secretName        = "opni-opensearch-auth"
	secretKey         = "password"
	dataPrepperName   = "opni-shipper"
	clusterOutputName = "opni-output"
	clusterFlowName   = "opni-flow"
)

var (
	skipTLSVerify   bool
	rancherLogging  bool
	inCluster       bool
	gatewayEndpoint string
	bootstrapToken  string
	provider        string
	namespace       string
	pins            []string

	k8sClient client.Client
)

type simpleIdentProvider struct {
	Client client.Client
}

func BuildBootstrapLoggingCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "logging cluster-name",
		Short: "Bootstrap the logging capability for a cluster",
		Args:  cobra.ExactArgs(1),
		RunE:  doBootstrap,
	}

	command.Flags().BoolVar(&skipTLSVerify, "opensearch-insecure", false, "skip Opensearch tls verification")
	command.Flags().BoolVar(&rancherLogging, "use-rancher-logging", false, "manually configure log shipping with rancher-logging")
	command.Flags().BoolVar(&inCluster, "in-cluster", false, "set bootstrap to run in cluster")
	command.Flags().StringVar(&gatewayEndpoint, "gateway-url", "https://localhost:8443", "upstream Opni gateway")
	command.Flags().StringVar(&provider, "provider", "rke", "the Kubernetes distribution")
	command.Flags().StringVar(&bootstrapToken, "token", "", "bootstrap token")
	command.Flags().StringVar(&namespace, "namespace", common.DefaultOpniNamespace, "namespace to use")
	trust.BindFlags(command.Flags())
	command.MarkFlagRequired("token")

	return command
}

func doBootstrap(cmd *cobra.Command, args []string) error {
	trustConfig, err := trust.BuildConfigFromFlags(cmd.Flags())
	if err != nil {
		return err
	}
	trustStrategy, err := trustConfig.Build()
	if err != nil {
		return err
	}

	k8sClient = common.GetClientOrDie(inCluster)
	identifier := &simpleIdentProvider{
		Client: k8sClient,
	}

	clusterID, err := identifier.UniqueIdentifier(cmd.Context())
	if err != nil {
		return err
	}

	bootstrapConfig, err := buildBoostrapClient(trustStrategy)
	if err != nil {
		return err
	}

	keyring, err := bootstrapConfig.Bootstrap(cmd.Context(), identifier)
	if err != nil {
		return err
	}

	gatewayClient, err := gatewayclients.NewGatewayHTTPClient(gatewayEndpoint, identifier, keyring, trustStrategy)
	if err != nil {
		return err
	}
	rb := gatewayClient.Get(cmd.Context(), "/logging/v1/cluster")

	code, body, err := rb.Send()
	if err != nil {
		return err
	}
	if code != 200 {
		return errors.New("unsuccessful request to fetch credentials")
	}

	var detailsResp loggingplugin.OpensearchDetailsResponse

	if err := json.Unmarshal(body, &detailsResp); err != nil {
		return err
	}

	if err := createAuthSecret(cmd.Context(), detailsResp.Password); err != nil {
		return err
	}

	if err := createDataPrepper(
		cmd.Context(),
		detailsResp.Username,
		clusterID,
		detailsResp.ExternalURL,
	); err != nil {
		return err
	}

	if !rancherLogging {
		if err := createOpniClusterOutput(cmd.Context()); err != nil {
			return err
		}
		if err := createOpniClusterFlow(cmd.Context(), clusterID); err != nil {
			return err
		}
		if err := createLogAdapter(cmd.Context()); err != nil {
			return err
		}
	}

	return nil
}

func createAuthSecret(ctx context.Context, password string) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		StringData: map[string]string{
			secretKey: password,
		},
	}

	return k8sClient.Create(ctx, secret)
}

func createDataPrepper(
	ctx context.Context,
	username string,
	clusterID string,
	opensearchEndpoint string,
) error {
	dataPrepper := v1beta2.DataPrepper{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dataPrepperName,
			Namespace: namespace,
		},
		Spec: v1beta2.DataPrepperSpec{
			Username: username,
			PasswordFrom: &corev1.SecretKeySelector{
				Key: secretKey,
				LocalObjectReference: corev1.LocalObjectReference{
					Name: secretName,
				},
			},
			Opensearch: &v1beta2.OpensearchSpec{
				Endpoint:                 opensearchEndpoint,
				InsecureDisableSSLVerify: skipTLSVerify,
			},
			ClusterID: clusterID,
		},
	}

	return k8sClient.Create(ctx, &dataPrepper)
}

func createOpniClusterOutput(ctx context.Context) error {
	clusterOutput := &loggingv1beta1.ClusterOutput{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterOutputName,
			Namespace: namespace,
		},
		Spec: loggingv1beta1.ClusterOutputSpec{
			OutputSpec: loggingv1beta1.OutputSpec{
				HTTPOutput: &output.HTTPOutputConfig{
					Endpoint:    fmt.Sprintf("http://%s.%s:2021/log/ingest", dataPrepperName, namespace),
					ContentType: "application/json",
					JsonArray:   true,
					Buffer: &output.Buffer{
						Tags:           "[]",
						FlushInterval:  "2s",
						ChunkLimitSize: "1mb",
					},
				},
			},
		},
	}
	return k8sClient.Create(ctx, clusterOutput)
}

func createOpniClusterFlow(ctx context.Context, clusterID string) error {
	clusterFlow := &loggingv1beta1.ClusterFlow{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterFlowName,
			Namespace: namespace,
		},
		Spec: loggingv1beta1.ClusterFlowSpec{
			Filters: []loggingv1beta1.Filter{
				{
					Dedot: &filter.DedotFilterConfig{
						Separator: "-",
						Nested:    true,
					},
				},
				{
					Grep: &filter.GrepConfig{
						Exclude: []filter.ExcludeSection{
							{
								Key:     "log",
								Pattern: `^\n$`,
							},
						},
					},
				},
				{
					DetectExceptions: &filter.DetectExceptions{
						Languages: []string{
							"java",
							"python",
							"go",
							"ruby",
							"js",
							"csharp",
							"php",
						},
						MultilineFlushInterval: "0.1",
					},
				},
				{
					RecordTransformer: &filter.RecordTransformer{
						Records: []filter.Record{
							{
								"cluster_id": clusterID,
							},
						},
					},
				},
			},
			Match: []loggingv1beta1.ClusterMatch{
				{
					ClusterExclude: &loggingv1beta1.ClusterExclude{
						Namespaces: []string{
							"opni-system",
						},
					},
				},
				{
					ClusterSelect: &loggingv1beta1.ClusterSelect{},
				},
			},
			GlobalOutputRefs: []string{
				clusterOutputName,
			},
		},
	}

	return k8sClient.Create(ctx, clusterFlow)
}

func createLogAdapter(ctx context.Context) error {
	i := 1
	lga := &v1beta2.LogAdapter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "opni-logging",
			Namespace: common.NamespaceFlagValue,
		},
		Spec: v1beta2.LogAdapterSpec{
			Provider: v1beta2.LogProvider(provider),
		},
	}
	for {
		err := k8sClient.Create(ctx, lga)
		if err == nil {
			return nil
		}
		// TODO: fix k8s apiserver and check for webserver error instead
		if !k8serrors.IsInternalError(err) {
			return err
		}
		retryBackoff := backoff.NewExponentialBackOff()
		if i == 5 {
			return err
		}
		time.Sleep(retryBackoff.NextBackOff())
		i += 1
	}

}

func buildBoostrapClient(trustStrategy trust.Strategy) (*bootstrap.ClientConfig, error) {
	token, err := tokens.ParseHex(bootstrapToken)
	if err != nil {
		return nil, err
	}

	return &bootstrap.ClientConfig{
		Capability:    wellknown.CapabilityLogs,
		Token:         token,
		Endpoint:      gatewayEndpoint,
		TrustStrategy: trustStrategy,
		K8sNamespace:  common.NamespaceFlagValue,
		K8sConfig:     common.RestConfig,
	}, nil
}

func (p *simpleIdentProvider) UniqueIdentifier(ctx context.Context) (string, error) {
	systemNamespace := &corev1.Namespace{}
	if err := p.Client.Get(ctx, types.NamespacedName{
		Name: "kube-system",
	}, systemNamespace); err != nil {
		return "", err
	}

	return string(systemNamespace.GetUID()), nil
}