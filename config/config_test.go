package config

import (
	"fmt"
	"os"
	"testing"
)

func TestGetPeersConfig(t *testing.T) {
	pc := GetPeersConfig()

	for _, value := range pc {
		fmt.Printf("Host: %s, Port:%s, EventHost:%s, EventPort:%s \n",
			value.Host, value.Port, value.EventHost, value.EventPort)

	}

}

func TestGetTlsEnabled(t *testing.T) {
	fmt.Printf("IsTlsEnabled: %v\n", IsTlsEnabled())

}

func TestMain(m *testing.M) {
	err := InitConfig("../test_resources/config_test.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(m.Run())
}
