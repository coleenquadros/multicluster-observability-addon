package otelcol

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_ConfigureExportersSecrets(t *testing.T) {
	b, err := os.ReadFile("./test_data/simplest.yaml")
	require.NoError(t, err)
	otelColConfig := string(b)
	cfg, err := ConfigFromString(otelColConfig)
	require.NoError(t, err)

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tracing-otlphttp-auth",
			Namespace: "cluster-1",
		},
		Data: map[string][]byte{
			"tls.crt": []byte("data"),
			"ca.crt":  []byte("data"),
			"tls.key": []byte("data"),
		},
	}

	err = ConfigureExportersSecrets(cfg, "otlp", secret)
	require.NoError(t, err)

	exportersField := cfg["exporters"]
	exporters := exportersField.(map[string]interface{})
	require.Nil(t, exporters["debug"])

	b, err = os.ReadFile("./test_data/basic_otelhttp.yaml")
	require.NoError(t, err)
	otelColConfig = string(b)
	cfg, err = ConfigFromString(otelColConfig)
	require.NoError(t, err)

	err = ConfigureExportersSecrets(cfg, "otlphttp", secret)
	require.NoError(t, err)

	exportersField = cfg["exporters"]
	exporters = exportersField.(map[string]interface{})
	otlphttpField := exporters["otlphttp"]
	otlphttp := otlphttpField.(map[string]interface{})
	require.NotNil(t, otlphttp["tls"])
}

func Test_getExporters(t *testing.T) {
	b, err := os.ReadFile("./test_data/simplest.yaml")
	require.NoError(t, err)
	otelColConfig := string(b)
	cfg, err := ConfigFromString(otelColConfig)
	require.NoError(t, err)

	exporters, err := GetExporters(cfg)
	require.NoError(t, err)
	require.Len(t, exporters, 1)

	cfg = make(map[string]interface{})
	_, err = GetExporters(cfg)
	require.Error(t, err)
}

func Test_configureExporterSecrets(t *testing.T) {
	exporter := make(map[string]interface{})
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tracing-otlphttp-auth",
			Namespace: "cluster-1",
		},
		Data: map[string][]byte{
			"tls.crt": []byte("data"),
			"ca.crt":  []byte("data"),
			"tls.key": []byte("data"),
		},
	}
	configureExporterSecrets(exporter, secret)
	require.NotNil(t, exporter["tls"])
}
