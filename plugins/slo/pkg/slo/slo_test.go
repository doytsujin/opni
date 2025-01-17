package slo_test

import (
	"context"
	"fmt"
	"os"
	"strings"

	oslov1 "github.com/alexandreLamarre/oslo/pkg/manifest/v1"
	"github.com/alexandreLamarre/sloth/core/app/generate"
	"github.com/alexandreLamarre/sloth/core/prometheus"
	"github.com/hashicorp/go-hclog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rancher/opni/pkg/test"
	apis "github.com/rancher/opni/plugins/slo/pkg/apis/slo"
	"github.com/rancher/opni/plugins/slo/pkg/slo"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Converting ServiceLevelObjective Messages to Prometheus Rules", Ordered, Label(test.Unit, test.Slow), func() {
	sloTestDataDir := "../../../../pkg/test/testdata/slo/"
	slo1 := &apis.ServiceLevelObjective{ // no alerts
		Id:          "foo-id",
		Name:        "foo-name",
		Datasource:  "monitoring",
		Description: "Some SLO",
		Services: []*apis.Service{
			{JobId: "foo-service", ClusterId: "foo-cluster", MetricName: "uptime", MetricId: "up"},
		},
		MonitorWindow:     "30d",
		MetricDescription: "Some metric",
		BudgetingInterval: "5m",
		Labels:            []*apis.Label{},
		Targets: []*apis.Target{
			{ValueX100: 9999},
		},
		Alerts: []*apis.Alert{},
	}

	multipleObjectives := proto.Clone(slo1).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
	multiClusterMultiService := proto.Clone(slo1).ProtoReflect().Interface().(*apis.ServiceLevelObjective)

	alertSLO := proto.Clone(slo1).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
	alertSLO.Alerts = []*apis.Alert{ // multiple alerts
		{
			Name:                    "alert-foo",
			NotificationTarget:      "email",
			NotificationDescription: "Send to email",
			Description:             "Alert when we breach the objective",
			ConditionType:           "budget",
			ThresholdType:           ">",
			OnNoData:                true,
			OnCreate:                true,
			OnBreach:                true,
			OnResolved:              true,
		},
	}

	multiAlerts := proto.Clone(slo1).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
	multiAlerts.Alerts = alertSLO.Alerts
	multiAlerts.Alerts = append(multiAlerts.Alerts, &apis.Alert{
		Name:                    "alert-bar",
		NotificationTarget:      "slack",
		NotificationDescription: "Send to slack",
		Description:             "Alert on burn rate",
		ConditionType:           "burnrate",
		ThresholdType:           ">",
		OnNoData:                true,
		OnCreate:                true,
		OnBreach:                true,
		OnResolved:              true,
	})

	var simpleSpec []oslov1.SLO
	var objectiveSpecs []oslov1.SLO
	var multiClusterSpecs []oslov1.SLO
	var alertSpecs []oslov1.SLO
	var multiAlertSpecs []oslov1.SLO

	var simplePrometheusIR []*prometheus.SLOGroup
	var objectivePrometheusIR []*prometheus.SLOGroup
	var multiClusterPrometheusIR []*prometheus.SLOGroup
	// var alertPrometheusIR []*prometheus.SLOGroup
	// var multiAlertPrometheusIR []*prometheus.SLOGroup

	var simplePrometheusResponse []*generate.Response
	var objectivePrometheusResponse []*generate.Response
	var multiClusterPrometheusResponse []*generate.Response
	// var alertPrometheusIR []*generate.Response
	// var multiAlertPrometheusIR []*generate.Response

	When("A ServiceLevelObjective message is given", func() {
		It("should validate proper input", func() {
			Expect(slo.ValidateInput(slo1)).To(Succeed())

			sloLogging := proto.Clone(slo1).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
			sloLogging.Datasource = "logging"
			Expect(slo.ValidateInput(sloLogging)).To(Succeed())

			Expect(slo.ValidateInput(alertSLO)).To(Succeed())
			Expect(slo.ValidateInput(multiAlerts)).To(Succeed())

			for _, atype := range []string{slo.NotifHook, slo.NotifPager, slo.NotifMail, slo.NotifSlack} {
				sloNewAlert := proto.Clone(alertSLO).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
				alertSLO.Alerts[0].NotificationTarget = atype
				Expect(slo.ValidateInput(sloNewAlert)).To(Succeed())
			}

			for _, ctype := range []string{slo.AlertingBurnRate, slo.AlertingBudget, slo.AlertingTarget} {
				sloNewAlert := proto.Clone(alertSLO).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
				alertSLO.Alerts[0].ConditionType = ctype
				Expect(slo.ValidateInput(sloNewAlert)).To(Succeed())
			}

			for _, ttype := range []string{slo.GTThresholdType, slo.LTThresholdType} {
				sloNewAlert := proto.Clone(alertSLO).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
				alertSLO.Alerts[0].ThresholdType = ttype
				Expect(slo.ValidateInput(sloNewAlert)).To(Succeed())
			}

		})
		It("should reject improper input", func() {
			invalidDesc := proto.Clone(slo1).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
			invalidDesc.Description = strings.Repeat("a", 1056)
			Expect(slo.ValidateInput(invalidDesc)).To(MatchError(slo.ErrInvalidDescription))

			invalidSource := proto.Clone(slo1).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
			invalidSource.Datasource = strings.Repeat("a", 256)
			Expect(slo.ValidateInput(invalidSource)).To(MatchError(slo.ErrInvalidDatasource))

			missingId := proto.Clone(slo1).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
			missingId.Id = ""
			Expect(slo.ValidateInput(missingId)).To(MatchError(slo.ErrInvalidId))

			sloInvalidAlertTarget := proto.Clone(alertSLO).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
			sloInvalidAlertTarget.Alerts[0].NotificationTarget = "invalid-234987ukjas"
			Expect(slo.ValidateInput(sloInvalidAlertTarget)).To(MatchError(slo.ErrInvalidAlertTarget))

			sloInvalidAlertCondition := proto.Clone(alertSLO).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
			sloInvalidAlertCondition.Alerts[0].ConditionType = "invalid-234987ukjas"
			Expect(slo.ValidateInput(sloInvalidAlertCondition)).To(MatchError(slo.ErrInvalidAlertCondition))

			sloInvalidAlertThreshold := proto.Clone(alertSLO).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
			sloInvalidAlertThreshold.Alerts[0].ThresholdType = "invalid-234987ukjas"
			Expect(slo.ValidateInput(sloInvalidAlertThreshold)).To(MatchError(slo.ErrInvalidAlertThreshold))

		})
	})
	When("We convert a ServiceLevelObjective to OpenSLO", func() {
		It("Should create a valid OpenSLO spec", func() {
			//Monitoring SLOs
			var err error
			Expect(slo1.Datasource).To(Equal("monitoring")) // make sure we didn't mutate original message
			simpleSpec, err = slo.ParseToOpenSLO(slo1, context.Background(), hclog.New(&hclog.LoggerOptions{}))
			Expect(err).To(Succeed())
			Expect(simpleSpec).To(HaveLen(1))

			expectedData, err := os.ReadFile(fmt.Sprintf("%s/slo1.yaml", sloTestDataDir))
			Expect(err).To(Succeed())

			createdData, err := yaml.Marshal(&simpleSpec[0])
			Expect(createdData).To(MatchYAML(expectedData))
			multipleObjectives.Targets = []*apis.Target{
				{ValueX100: 9999},
				{ValueX100: 9995},
			}

			objectiveSpecs, err = slo.ParseToOpenSLO(multipleObjectives, context.Background(), hclog.New(&hclog.LoggerOptions{}))
			Expect(err).To(Succeed())

			Expect(objectiveSpecs).To(HaveLen(1))
			expectedObjectives, err := os.ReadFile(fmt.Sprintf("%s/objectives.yaml", sloTestDataDir))
			Expect(yaml.Marshal(&objectiveSpecs[0])).To(MatchYAML(expectedObjectives))

			multiClusterMultiService.Services = []*apis.Service{
				{JobId: "foo-service", ClusterId: "foo-cluster", MetricName: "uptime", MetricId: "up"},
				{JobId: "foo-service2", ClusterId: "foo-cluster", MetricName: "uptime", MetricId: "up"},
				{JobId: "foo-service", ClusterId: "bar-cluster", MetricName: "uptime", MetricId: "up"},
				{JobId: "foo-service2", ClusterId: "bar-cluster", MetricName: "uptime", MetricId: "up"},
			}

			multiClusterSpecs, err = slo.ParseToOpenSLO(multiClusterMultiService, context.Background(), hclog.New(&hclog.LoggerOptions{}))
			Expect(err).To(Succeed())

			Expect(multiClusterSpecs).To(HaveLen(4))
			expectedMulti, err := os.ReadFile(fmt.Sprintf("%s/multiSpecs.yaml", sloTestDataDir))
			Expect(err).To(Succeed())
			Expect(yaml.Marshal(multiClusterSpecs)).To(MatchYAML(expectedMulti))

			// Alerting SLOS
			alertSpecs, err = slo.ParseToOpenSLO(alertSLO, context.Background(), hclog.New(&hclog.LoggerOptions{}))
			Expect(err).To(Succeed())
			Expect(alertSpecs).To(HaveLen(1))
			Expect(alertSpecs[0].Spec.AlertPolicies).To(HaveLen(1))
			//FIXME: this is a check to verify we haven't implemented this functionality yet
			Expect(alertSpecs[0].Spec.AlertPolicies[0].Spec.Conditions).To(HaveLen(0))
			expectedAlert, err := os.ReadFile(fmt.Sprintf("%s/alertSLO.yaml", sloTestDataDir))
			Expect(err).To(Succeed())
			Expect(yaml.Marshal(alertSpecs[0])).To(MatchYAML(expectedAlert))

			multiAlertSpecs, err = slo.ParseToOpenSLO(multiAlerts, context.Background(), hclog.New(&hclog.LoggerOptions{}))
			Expect(err).To(Succeed())
			expectedMultiAlert, err := os.ReadFile(fmt.Sprintf("%s/multiAlertSLO.yaml", sloTestDataDir))
			Expect(err).To(Succeed())
			Expect(yaml.Marshal(multiAlertSpecs)).To(MatchYAML(expectedMultiAlert))

			// Logging SLOs
			sloLogging := proto.Clone(slo1).ProtoReflect().Interface().(*apis.ServiceLevelObjective)
			sloLogging.Datasource = "logging"
			_, err = slo.ParseToOpenSLO(sloLogging, context.Background(), hclog.New(&hclog.LoggerOptions{}))
			Expect(err).To(MatchError("Not implemented"))
		})
	})

	When("We convert an OpenSLO to Sloth IR for OpniMonitoring", func() {
		It("should create a valid Sloth IR", func() {
			var err error
			simplePrometheusIR, err = slo.ParseToPrometheusModel(simpleSpec)
			Expect(err).To(Succeed())
			Expect(simplePrometheusIR).To(HaveLen(1))

			objectivePrometheusIR, err = slo.ParseToPrometheusModel(objectiveSpecs)
			Expect(err).To(Succeed())
			Expect(objectivePrometheusIR).To(HaveLen(1))
			// For each objective defined by user we expect to have a corresponding SLOGroup
			Expect(objectivePrometheusIR[0].SLOs).To(HaveLen(2))

			multiClusterPrometheusIR, err = slo.ParseToPrometheusModel(multiClusterSpecs)
			Expect(err).To(Succeed())
			Expect(multiClusterPrometheusIR).To(HaveLen(4))
			for _, ir := range multiClusterPrometheusIR {
				Expect(ir.SLOs).To(HaveLen(1))
			}
		})
	})

	When("We convert Sloth IR to prometheus rules", func() {
		It("Should create valid prometheus rules", func() {
			var err error
			simplePrometheusResponse, err = slo.GeneratePrometheusRule(simplePrometheusIR, context.Background())
			Expect(err).To(MatchError("Prometheus generator failed to start"))
			Expect(simplePrometheusResponse).To(BeNil())

			objectivePrometheusResponse, err = slo.GeneratePrometheusRule(objectivePrometheusIR, context.Background())
			Expect(err).To(MatchError("Prometheus generator failed to start"))
			Expect(objectivePrometheusResponse).To(BeNil())

			multiClusterPrometheusResponse, err = slo.GeneratePrometheusRule(multiClusterPrometheusIR, context.Background())
			Expect(err).To(MatchError("Prometheus generator failed to start"))
			Expect(multiClusterPrometheusResponse).To(BeNil())
		})
	})

})
