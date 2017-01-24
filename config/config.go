package config

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

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

var log = logging.MustGetLogger("fabric_sdk_go")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} [%{module}] %{level:.4s} : %{message}`,
)

// initConfig reads in config file
func InitConfig(configFile string) error {

	if configFile != "" {
		viper.SetConfigFile(configFile)
		// If a config file is found, read it in.
		err := viper.ReadInConfig()

		if err == nil {
			log.Infof("Using config file: %s", viper.ConfigFileUsed())
		} else {
			return fmt.Errorf("Fatal error config file: %v", err)
		}
		//		viper.WatchConfig()
		//		viper.OnConfigChange(func(e fsnotify.Event) {
		//			log.Infof("Config file changed: %s", e.Name)
		//		})

	}

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	loggingLevelString := viper.GetString("client.logging.level")
	logLevel := logging.INFO
	if loggingLevelString != "" {
		log.Infof("fabric_sdk_go Logging level: %v", loggingLevelString)
		var err error
		logLevel, err = logging.LogLevel(loggingLevelString)
		if err != nil {
			panic(err)
		}
	}
	logging.SetBackend(backendFormatter).SetLevel(logging.Level(logLevel), "fabric_sdk_go")

	return nil
}

func GetPeersConfig() []PeerConfig {
	peersConfig := []PeerConfig{}
	peers := viper.GetStringMap("client.peers")
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
	return viper.GetBool("client.tls.enabled")
}

func GetTlsCACertPool() *x509.CertPool {
	certPool := x509.NewCertPool()
	if viper.GetString("client.tls.certificate") != "" {
		rawData, err := ioutil.ReadFile(viper.GetString("tls.certificate"))
		if err != nil {
			panic(err)
		}
		certPool.AddCert(loadCAKey(rawData))
	}
	return certPool
}

func GetTlsServerHostOverride() string {
	return viper.GetString("client.tls.serverhostoverride")
}

func IsSecurityEnabled() bool {
	return viper.GetBool("client.security.enabled")
}
func TcertBatchSize() int {
	return viper.GetInt("client.tcert.batch.size")
}
func GetSecurityAlgorithm() string {
	return viper.GetString("client.security.hashAlgorithm")
}
func GetSecurityLevel() int {
	return viper.GetInt("client.security.level")

}
func GetOrdererHost() string {
	return viper.GetString("client.orderer.host")
}

func GetMspUrl() string {
	return viper.GetString("client.msp.url")
}

func GetMspId() string {
	return viper.GetString("client.msp.id")
}

func GetMspClientPath() string {
	return viper.GetString("client.msp.clientPath")
}

func GetKeyStorePath() string {
	return viper.GetString("client.keystore.path")
}

func GetOrdererPort() string {
	return strconv.Itoa(viper.GetInt("client.orderer.port"))
}

func loadCAKey(rawData []byte) *x509.Certificate {
	block, _ := pem.Decode(rawData)

	pub, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}
	return pub
}
