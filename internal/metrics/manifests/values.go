package manifests

import (
	"encoding/json"

	"github.com/rhobs/multicluster-observability-addon/internal/metrics/config"
	"github.com/rhobs/multicluster-observability-addon/internal/metrics/handlers"
	corev1 "k8s.io/api/core/v1"
)

type MetricsValues struct {
	PlatformEnabled           bool          `json:"platformEnabled"`
	UserWorkloadsEnabled      bool          `json:"userWorkloadsEnabled"`
	Secrets                   []ConfigValue `json:"secrets"`
	Images                    ImagesValues  `json:"images"`
	PrometheusControllerID    string        `json:"prometheusControllerID"`
	PrometheusCAConfigMapName string        `json:"prometheusCAConfigMapName"`
	Platform                  Collector     `json:"platform"`
	UserWorkload              Collector     `json:"userWorkload"`
}

type Collector struct {
	AppName             string        `json:"appName"`
	ConfigMaps          []ConfigValue `json:"configMaps"`
	PrometheusAgentSpec string        `json:"prometheusAgentSpec"`
	ScrapeConfigs       []ConfigValue `json:"scrapeConfigs"`
	Rules               []ConfigValue `json:"rules"`
}

type ImagesValues struct {
	PrometheusOperator       string `json:"prometheusOperator"`
	PrometheusConfigReloader string `json:"prometheusConfigReloader"`
}

type ConfigValue struct {
	Name   string            `json:"name"`
	Data   string            `json:"data"`
	Labels map[string]string `json:"labels"`
}

func BuildValues(opts handlers.Options) (MetricsValues, error) {
	ret := MetricsValues{
		PrometheusControllerID:    config.PrometheusControllerID,
		PrometheusCAConfigMapName: config.PrometheusCAConfigMapName,
		Platform: Collector{
			AppName: config.PlatformMetricsCollectorApp,
		},
		UserWorkload: Collector{
			AppName: config.UserWorkloadMetricsCollectorApp,
		},
	}

	// Build Prometheus Agent Spec for Platform
	if opts.Platform.PrometheusAgent != nil {
		ret.PlatformEnabled = true

		agentJson, err := json.Marshal(opts.Platform.PrometheusAgent.Spec)
		if err != nil {
			return ret, err
		}

		ret.Platform.PrometheusAgentSpec = string(agentJson)
	}

	// Build Prometheus Agent Spec for User Workloads
	if opts.UserWorkloads.PrometheusAgent != nil {
		ret.UserWorkloadsEnabled = true

		agentJson, err := json.Marshal(opts.UserWorkloads.PrometheusAgent.Spec)
		if err != nil {
			return ret, err
		}

		ret.UserWorkload.PrometheusAgentSpec = string(agentJson)
	}

	// Build scrape configs
	for _, scrapeConfig := range opts.Platform.ScrapeConfigs {
		scrapeConfigJson, err := json.Marshal(scrapeConfig.Spec)
		if err != nil {
			return ret, err
		}

		ret.Platform.ScrapeConfigs = append(ret.Platform.ScrapeConfigs, ConfigValue{
			Name:   scrapeConfig.Name,
			Data:   string(scrapeConfigJson),
			Labels: scrapeConfig.Labels,
		})
	}

	for _, scrapeConfig := range opts.UserWorkloads.ScrapeConfigs {
		scrapeConfigJson, err := json.Marshal(scrapeConfig.Spec)
		if err != nil {
			return ret, err
		}

		ret.UserWorkload.ScrapeConfigs = append(ret.UserWorkload.ScrapeConfigs, ConfigValue{
			Name:   scrapeConfig.Name,
			Data:   string(scrapeConfigJson),
			Labels: scrapeConfig.Labels,
		})
	}

	// Build rules
	for _, rule := range opts.Platform.Rules {
		ruleJson, err := json.Marshal(rule.Spec)
		if err != nil {
			return ret, err
		}

		ret.Platform.Rules = append(ret.Platform.Rules, ConfigValue{
			Name:   rule.Name,
			Data:   string(ruleJson),
			Labels: rule.Labels,
		})
	}

	for _, rule := range opts.UserWorkloads.Rules {
		ruleJson, err := json.Marshal(rule.Spec)
		if err != nil {
			return ret, err
		}

		ret.UserWorkload.Rules = append(ret.UserWorkload.Rules, ConfigValue{
			Name:   rule.Name,
			Data:   string(ruleJson),
			Labels: rule.Labels,
		})
	}

	// Build secrets and config maps
	var err error
	ret.Secrets, err = buildSecrets(opts.Secrets)
	if err != nil {
		return ret, err
	}

	// Set config maps
	ret.Platform.ConfigMaps, err = buildConfigMaps(opts.Platform.ConfigMaps)
	if err != nil {
		return ret, err
	}

	ret.UserWorkload.ConfigMaps, err = buildConfigMaps(opts.UserWorkloads.ConfigMaps)
	if err != nil {
		return ret, err
	}

	// Set images
	ret.Images = ImagesValues{
		PrometheusOperator:       opts.Images.PrometheusOperator,
		PrometheusConfigReloader: opts.Images.PrometheusConfigReloader,
	}

	return ret, nil
}

func buildSecrets(secrets []*corev1.Secret) ([]ConfigValue, error) {
	secretsValue := []ConfigValue{}
	for _, secret := range secrets {
		dataJSON, err := json.Marshal(secret.Data)
		if err != nil {
			return secretsValue, err
		}
		secretValue := ConfigValue{
			Name:   secret.Name,
			Data:   string(dataJSON),
			Labels: secret.Labels,
		}
		secretsValue = append(secretsValue, secretValue)
	}
	return secretsValue, nil
}

func buildConfigMaps(configMaps []*corev1.ConfigMap) ([]ConfigValue, error) {
	configMapsValue := []ConfigValue{}
	for _, configMap := range configMaps {
		dataJSON, err := json.Marshal(configMap.Data)
		if err != nil {
			return configMapsValue, err
		}
		configMapValue := ConfigValue{
			Name:   configMap.Name,
			Data:   string(dataJSON),
			Labels: configMap.Labels,
		}
		configMapsValue = append(configMapsValue, configMapValue)
	}
	return configMapsValue, nil
}
