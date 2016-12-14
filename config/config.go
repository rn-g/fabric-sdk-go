package config

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

const (
	PEERS = "peers"
)

type PeerConfig struct {
	Host      string
	Port      string
	EventHost string
	EventPort string
}

var myLogger = logging.MustGetLogger("config")

// initConfig reads in config file
func InitConfig(configFile string) error {

	if configFile != "" {
		viper.SetConfigFile(configFile)
		// If a config file is found, read it in.
		err := viper.MergeInConfig()

		if err == nil {
			myLogger.Infof("Using config file: %s", viper.ConfigFileUsed())
		} else {
			return fmt.Errorf("Fatal error config file: %v", err)
		}
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			myLogger.Infof("Config file changed: %s", e.Name)
		})

	}

	return nil
}

func GetPeersConfig() []PeerConfig {
	peersConfig := []PeerConfig{}
	peers := viper.GetStringMap("peers")
	for key, value := range peers {
		mm := value.(map[interface{}]interface{})
		host, _ := mm["host"].(string)
		port, _ := mm["port"].(int)
		eventHost, _ := mm["event_host"].(string)
		eventPort, _ := mm["event_port"].(int)

		p := PeerConfig{Host: host, Port: strconv.Itoa(port), EventHost: eventHost, EventPort: strconv.Itoa(eventPort)}
		if p.Host == "" {
			panic(fmt.Sprintf("host key not exist or empty for %s", key))
		}
		if p.Port == "" {
			panic(fmt.Sprintf("port key not exist or empty for %s", key))
		}
		if p.EventHost == "" {
			panic(fmt.Sprintf("event_host not exist or empty for %s", key))
		}
		if p.EventPort == "" {
			panic(fmt.Sprintf("event_port not exist or empty for %s", key))
		}
		peersConfig = append(peersConfig, p)
	}
	return peersConfig

}

func IsTlsEnabled() bool {
	return viper.GetBool("tls.enabled")
}

func GetTlsCACertPool() *x509.CertPool {
	certPool := x509.NewCertPool()
	if viper.GetString("tls.certificate") != "" {
		rawData, err := ioutil.ReadFile(viper.GetString("tls.certificate"))
		if err != nil {
			panic(err)
		}
		certPool.AddCert(loadCAKey(rawData))
	}
	return certPool
}

func GetTlsServerHostOverride() string {
	return viper.GetString("tls.serverhostoverride")
}

func IsSecurityEnabled() bool {
	return viper.GetBool("security.enabled")
}
func TcertBatchSize() int {
	return viper.GetInt("tcert.batch.size")
}
func GetSecurityAlgorithm() string {
	return viper.GetString("security.hashAlgorithm")
}
func GetSecurityLevel() int {
	return viper.GetInt("security.level")

}

func loadCAKey(rawData []byte) *x509.Certificate {
	block, _ := pem.Decode(rawData)

	pub, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}
	return pub
}
