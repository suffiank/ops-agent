package port_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/ops-agent/internal/port"
	"github.com/google/go-cmp/cmp"
)

const testDataDir = "testdata/config"

var baseConfigPath = fmt.Sprintf("%s/base.yaml", testDataDir)

func TestConfigReadExisting(t *testing.T) {
	_, err := port.ReadConfig(baseConfigPath)
	if err != nil {
		t.Fatalf("expected to successfully read config file, got error: %v", err)
	}
}

func TestConfigReadCreatesIfNotExists(t *testing.T) {
	nonexistentPath := fmt.Sprintf("%s/nonexistent.yaml", testDataDir)

	config, err := port.ReadConfig(nonexistentPath)
	if err != nil {
		t.Fatalf("expected to successfully create new config file, got error: %v", err)
	}
	if len(config.ReservedPorts) > 0 {
		t.Fatalf("expected no ports to be reserved, reserved ports list is: %v", config.ReservedPorts)
	}
	if _, err := os.Stat(nonexistentPath); err != nil {
		t.Fatalf("expected %s to exist, got error: %v", nonexistentPath, err)
	}

	os.Remove(nonexistentPath)
}

func TestConfigWrite(t *testing.T) {
	writeTestPath := fmt.Sprintf("%s/written.yaml", testDataDir)

	t.Run("config write test", func(t *testing.T) {
		config, err := port.ReadConfig(baseConfigPath)
		if err != nil {
			t.Fatalf("expected to successfully read config file, got error: %v", err)
		}
		err = port.WriteConfig(writeTestPath, config)
		if err != nil {
			t.Fatalf("expected to successfully write config file, got error: %v", err)
		}

		baseConfigContent, err := os.ReadFile(baseConfigPath)
		if err != nil {
			t.Fatalf("expected to successfully read config file contents, got error: %v", err)
		}
		writtenConfigContent, err := os.ReadFile(writeTestPath)
		if err != nil {
			t.Fatalf("expected to successfully read written config file contents, got error: %v", err)
		}

		if diff := cmp.Diff(string(baseConfigContent), string(writtenConfigContent)); diff != "" {
			t.Fatalf("expected written config to match base config contents, got diff: %s", diff)
		}
	})

	os.Remove(writeTestPath)
}

func TestChangeConfigOverwrites(t *testing.T) {
	t.Parallel()
	overwriteTestPath := fmt.Sprintf("%s/overwritten.yaml", testDataDir)

	t.Run("config overwrite test", func(t *testing.T) {
		config, err := port.ReadConfig(baseConfigPath)
		if err != nil {
			t.Fatalf("expected to successfully read config file, got error: %v", err)
		}
		err = port.WriteConfig(overwriteTestPath, config)
		if err != nil {
			t.Fatalf("expected to successfully write config file, got error: %v", err)
		}

		config.ReservedPorts["key1"] = 3000
		config.ReservedPorts["new_key"] = 8080

		err = port.WriteConfig(overwriteTestPath, config)
		if err != nil {
			t.Fatalf("expected to successfully write config file, got error: %v", err)
		}

		writtenConfig, err := port.ReadConfig(overwriteTestPath)
		if err != nil {
			t.Fatalf("expected to successfully read config file, got error: %v", err)
		}

		if diff := cmp.Diff(config.ReservedPorts, writtenConfig.ReservedPorts); diff != "" {
			t.Fatalf("expected read final config to match the modified config used to write, got diff: %s", diff)
		}
	})

	os.Remove(overwriteTestPath)
}
