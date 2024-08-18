package config_test

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"git.cryptic.systems/volker.raschek/tarr/pkg/config"
	"github.com/stretchr/testify/require"
)

//go:embed test/assets/xml/config.xml
var expectedXMLConfig string

//go:embed test/assets/yaml/config.yaml
var expectedYAMLConfig string

func TestReadWriteConfig_XML(t *testing.T) {
	require := require.New(t)

	tmpDir, err := os.MkdirTemp("", "*")
	require.NoError(err)
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	expectedXMLConfigName := filepath.Join(tmpDir, "expected_config.xml")
	f, err := os.Create(expectedXMLConfigName)
	require.NoError(err)

	_, err = f.WriteString(expectedXMLConfig)
	require.NoError(err)

	actualConfig, err := config.ReadConfig(expectedXMLConfigName)
	require.NoError(err)
	require.NotNil(actualConfig)

	actualXMLConfigName := filepath.Join(tmpDir, "actual_config.xml")
	err = config.WriteConfig(actualXMLConfigName, actualConfig)
	require.NoError(err)

	b, err := os.ReadFile(actualXMLConfigName)
	require.NoError(err)
	require.Equal(expectedXMLConfig, string(b))
}

func TestReadWriteConfig_YAML(t *testing.T) {
	require := require.New(t)

	tmpDir, err := os.MkdirTemp("", "*")
	require.NoError(err)
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	expectedYAMLConfigName := filepath.Join(tmpDir, "expected_config.yaml")
	f, err := os.Create(expectedYAMLConfigName)
	require.NoError(err)

	_, err = f.WriteString(expectedYAMLConfig)
	require.NoError(err)

	actualConfig, err := config.ReadConfig(expectedYAMLConfigName)
	require.NoError(err)
	require.NotNil(actualConfig)

	actualYAMLConfigName := filepath.Join(tmpDir, "actual_config.yaml")
	err = config.WriteConfig(actualYAMLConfigName, actualConfig)
	require.NoError(err)

	b, err := os.ReadFile(actualYAMLConfigName)
	require.NoError(err)
	require.Equal(expectedYAMLConfig, string(b))
}
