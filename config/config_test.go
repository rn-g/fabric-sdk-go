package config

import (
	"fmt"
	"os"
	"testing"
)

func TestGetPeersConfig(t *testing.T) {
	pc := GetPeersConfig()

	for _, value := range pc {
		if value.Host == "" {
			t.Fatalf("Host value is empty")
		}
		if value.Port == "" {
			t.Fatalf("Port value is empty")
		}
		if value.Port == "" {
			t.Fatalf("EventHost value is empty")
		}
		if value.Port == "" {
			t.Fatalf("EventPort value is empty")
		}

	}

}

func TestMain(m *testing.M) {
	err := InitConfig("../integration_test/test_resources/config1/config_test.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(m.Run())
}
